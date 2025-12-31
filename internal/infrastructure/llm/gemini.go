package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"google.golang.org/genai"
)

type GeminiClient struct {
	client *genai.Client
}

func NewGeminiClient(apiKey string) *GeminiClient {
	ctx := context.Background()
	// Removed Backend field as it caused undefined error, assuming default or auto-detection.
	// If needed, check library version specifics.
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		log.Printf("Failed to create Gemini client: %v", err)
		return nil
	}
	return &GeminiClient{
		client: client,
	}
}

// Ensure GeminiClient implements LLMProvider and Embedder interface
var _ LLMProvider = (*GeminiClient)(nil)
var _ Embedder = (*GeminiClient)(nil)

func (c *GeminiClient) Generate(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if c.client == nil {
		return nil, fmt.Errorf("gemini client is not initialized")
	}

	if req.Model == "" {
		req.Model = "gemini-1.5-flash"
	}

	// Map Messages to Gemini Parts
	// Note: The officially genai library `GenerateContent` usually takes parts directly for a single turn
	// or uses a chat session for multi-turn.
	// For compat with "Provider" interface stateless style, we construct a prompt.
	// However, `genai.GenerateContent` accepts `*Content` list if using `Models.GenerateContent`.
	// But `client.Models.GenerateContent` signature in v0.1.0 (checking docs via assumption/standard):
	// func (s *ModelsService) GenerateContent(ctx context.Context, model string, contents []*Content, config *GenerateContentConfig) (*GenerateContentResponse, error)

	// Let's assume standard usage for the library.

	var contents []*genai.Content
	for _, msg := range req.Messages {
		role := "user"
		if msg.Role == "assistant" || msg.Role == "model" {
			role = "model"
		}
		// Map system to user or model based on requirements, or omit/prepend.
		if msg.Role == "system" {
			role = "user"
		}

		contents = append(contents, &genai.Content{
			Role: role,
			Parts: []*genai.Part{
				{Text: msg.Content},
			},
		})
	}

	// Config
	config := &genai.GenerateContentConfig{
		Temperature: &req.Temperature,
		TopP:        &req.TopP,
	}
	if req.MaxTokens > 0 {
		val := int32(req.MaxTokens)
		config.MaxOutputTokens = val
	}

	// s.Models.GenerateContent signature check fix
	resp, err := c.client.Models.GenerateContent(ctx, req.Model, contents, config)
	if err != nil {
		return nil, fmt.Errorf("gemini generate error: %w", err)
	}

	// Extract Content
	// Response structure: resp.Candidates[0].Content.Parts[0].Text
	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no candidates returned")
	}

	text := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			text += part.Text
		}
	}

	// Usage Metadata - fixing access assuming they are values based on error "int32 and nil"
	usage := Usage{}
	if resp.UsageMetadata != nil {
		usage.PromptTokens = int(resp.UsageMetadata.PromptTokenCount)
		usage.CompletionTokens = int(resp.UsageMetadata.CandidatesTokenCount)
		usage.TotalTokens = int(resp.UsageMetadata.TotalTokenCount)
	}

	return &CompletionResponse{
		Content: text,
		Usage:   usage,
	}, nil
}

func (c *GeminiClient) Stream(ctx context.Context, req *CompletionRequest) (<-chan CompletionChunk, <-chan error) {
	outputChan := make(chan CompletionChunk, 10)
	errChan := make(chan error, 1)

	if c.client == nil {
		errChan <- fmt.Errorf("gemini client is not initialized")
		close(outputChan)
		close(errChan)
		return outputChan, errChan
	}

	go func() {
		defer close(outputChan)
		defer close(errChan)

		if req.Model == "" {
			req.Model = "gemini-1.5-flash"
		}

		// Construct prompt parts (same as Generate)
		var contents []*genai.Content
		for _, msg := range req.Messages {
			role := "user"
			if msg.Role == "assistant" || msg.Role == "model" {
				role = "model"
			}
			if msg.Role == "system" {
				role = "user"
			}
			contents = append(contents, &genai.Content{
				Role: role,
				Parts: []*genai.Part{
					{Text: msg.Content},
				},
			})
		}

		config := &genai.GenerateContentConfig{
			Temperature: &req.Temperature,
			TopP:        &req.TopP,
		}
		if req.MaxTokens > 0 {
			val := int32(req.MaxTokens)
			config.MaxOutputTokens = val
		}
		// Enable Tools if present
		// Enable Tools if present
		if len(req.Tools) > 0 {
			var functionDecls []*genai.FunctionDeclaration
			for _, t := range req.Tools {
				paramsMap, ok := t.Function.Parameters.(map[string]interface{})
				var schema *genai.Schema
				if ok && paramsMap != nil {
					schema = schemaFromMap(paramsMap)
				}

				functionDecls = append(functionDecls, &genai.FunctionDeclaration{
					Name:        t.Function.Name,
					Description: t.Function.Description,
					Parameters:  schema,
				})
			}
			config.Tools = []*genai.Tool{
				{
					FunctionDeclarations: functionDecls,
				},
			}
		}

		iter := c.client.Models.GenerateContentStream(ctx, req.Model, contents, config)

		for resp, err := range iter {
			if err != nil {
				errChan <- fmt.Errorf("gemini stream error: %w", err)
				return
			}

			if resp != nil {
				if len(resp.Candidates) > 0 {
					for _, part := range resp.Candidates[0].Content.Parts {
						chunk := CompletionChunk{}

						// 1. Text Content
						if part.Text != "" {
							chunk.Content = part.Text
						}

						// 2. Function Calls
						if part.FunctionCall != nil {
							// Map to ToolCalls
							// Gemini streams one function call per part usually?
							// We need a dummy ID as Gemini might not send one in the stream same way OpenAI does?
							// Using "call_" + timestamp or random?
							// Actually, for streaming *aggregation*, we need a stable ID if it's split.
							// But Gemini usually sends the whole FunctionCall in one chunk (not streamed char-by-char for args).
							// So we can generate an ID here.

							argsBytes, _ := json.Marshal(part.FunctionCall.Args) // Args is map[string]interface{}
							chunk.ToolCalls = []ToolCall{
								{
									ID:   "call_" + part.FunctionCall.Name, // distinct enough
									Type: "function",
									Function: FunctionCall{
										Name:      part.FunctionCall.Name,
										Arguments: string(argsBytes),
									},
								},
							}
						}

						// Only send if we have something
						if chunk.Content != "" || len(chunk.ToolCalls) > 0 {
							outputChan <- chunk
						}
					}
				}

				if resp.UsageMetadata != nil {
					outputChan <- CompletionChunk{
						Usage: &Usage{
							PromptTokens:     int(resp.UsageMetadata.PromptTokenCount),
							CompletionTokens: int(resp.UsageMetadata.CandidatesTokenCount),
							TotalTokens:      int(resp.UsageMetadata.TotalTokenCount),
						},
					}
				}
			}
		}
	}()

	return outputChan, errChan
}

func (c *GeminiClient) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	if c.client == nil {
		return nil, fmt.Errorf("gemini client is not initialized")
	}
	// Default embedding model
	if model == "" {
		model = "text-embedding-004"
	}

	// Prepare request
	part := &genai.Part{Text: text}
	contents := []*genai.Content{{
		Parts: []*genai.Part{part},
	}}

	resp, err := c.client.Models.EmbedContent(ctx, model, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("gemini embed error: %w", err)
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return resp.Embeddings[0].Values, nil
}

// schemaFromMap converts a generic JSON Schema map to genai.Schema.
// This matches the capabilities needed for standard tool definitions (strings, objects, arrays).
func schemaFromMap(m map[string]interface{}) *genai.Schema {
	s := &genai.Schema{}

	if t, ok := m["type"].(string); ok {
		switch t {
		case "object":
			s.Type = genai.TypeObject
		case "string":
			s.Type = genai.TypeString
		case "integer":
			s.Type = genai.TypeInteger
		case "number":
			s.Type = genai.TypeNumber
		case "boolean":
			s.Type = genai.TypeBoolean
		case "array":
			s.Type = genai.TypeArray
		}
	}

	if desc, ok := m["description"].(string); ok {
		s.Description = desc
	}

	if props, ok := m["properties"].(map[string]interface{}); ok {
		s.Properties = make(map[string]*genai.Schema)
		for k, v := range props {
			if vMap, ok := v.(map[string]interface{}); ok {
				s.Properties[k] = schemaFromMap(vMap)
			}
		}
	}

	if req, ok := m["required"].([]interface{}); ok {
		for _, r := range req {
			if rStr, ok := r.(string); ok {
				s.Required = append(s.Required, rStr)
			}
		}
	}
	
	// Handle Array Items
	if items, ok := m["items"].(map[string]interface{}); ok {
		s.Items = schemaFromMap(items)
	}

	return s
}
