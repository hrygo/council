package llm

import (
	"context"
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

func (c *GeminiClient) Stream(ctx context.Context, req *CompletionRequest) (<-chan string, <-chan error) {
	if c.client == nil {
		return nil, nil // Or return error generic
	}

	outputChan := make(chan string)
	errChan := make(chan error, 1)

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

		iter := c.client.Models.GenerateContentStream(ctx, req.Model, contents, config)

		for resp, err := range iter {
			if err != nil {
				errChan <- fmt.Errorf("gemini stream error: %w", err)
				return // Or continue if stream, but typically error ends stream
			}

			if resp != nil && len(resp.Candidates) > 0 {
				for _, part := range resp.Candidates[0].Content.Parts {
					if part.Text != "" {
						outputChan <- part.Text
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
