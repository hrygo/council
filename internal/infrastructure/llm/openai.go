package llm

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	return NewOpenAICompatibleClient(apiKey, "https://api.openai.com/v1")
}

func NewOpenAICompatibleClient(apiKey, baseURL string) *OpenAIClient {
	config := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		config.BaseURL = baseURL
	}
	client := openai.NewClientWithConfig(config)
	return &OpenAIClient{
		client: client,
	}
}

// Ensure OpenAIClient implements LLMProvider and Embedder interfaces
var _ LLMProvider = (*OpenAIClient)(nil)
var _ Embedder = (*OpenAIClient)(nil)

func (c *OpenAIClient) Generate(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	// Model is required. If empty, the upstream API will return an error which we handle below.

	messages := make([]openai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		messages[i] = openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       req.Model,
			Messages:    messages,
			Temperature: req.Temperature,
			TopP:        req.TopP,
			MaxTokens:   req.MaxTokens,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("openai generate error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	return &CompletionResponse{
		Content: resp.Choices[0].Message.Content,
		Usage: Usage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}, nil
}

func (c *OpenAIClient) Stream(ctx context.Context, req *CompletionRequest) (<-chan CompletionChunk, <-chan error) {
	outputChan := make(chan CompletionChunk)
	errChan := make(chan error, 1)

	go func() {
		defer close(outputChan)
		defer close(errChan)

		// Model is required for OpenAI-style APIs.
		// If empty, the provider usually returns an error (handled below).
		// We avoid hardcoding GPT-4o here to support compatible providers (DeepSeek, GLM, etc.).

		messages := make([]openai.ChatCompletionMessage, len(req.Messages))
		for i, msg := range req.Messages {
			messages[i] = openai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			}
			// Map ToolCalls if they form part of history
			if len(msg.ToolCalls) > 0 {
				tcList := make([]openai.ToolCall, len(msg.ToolCalls))
				for j, tc := range msg.ToolCalls {
					tcList[j] = openai.ToolCall{
						ID:   tc.ID,
						Type: openai.ToolType(tc.Type),
						Function: openai.FunctionCall{
							Name:      tc.Function.Name,
							Arguments: tc.Function.Arguments,
						},
					}
				}
				messages[i].ToolCalls = tcList
			}
			if msg.ToolCallID != "" {
				messages[i].ToolCallID = msg.ToolCallID
			}
		}

		var tools []openai.Tool
		if len(req.Tools) > 0 {
			tools = make([]openai.Tool, len(req.Tools))
			for i, t := range req.Tools {
				tools[i] = openai.Tool{
					Type: openai.ToolType(t.Type),
					Function: &openai.FunctionDefinition{
						Name:        t.Function.Name,
						Description: t.Function.Description,
						Parameters:  t.Function.Parameters,
					},
				}
			}
		}

		stream, err := c.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:       req.Model,
				Messages:    messages,
				Temperature: req.Temperature,
				TopP:        req.TopP,
				MaxTokens:   req.MaxTokens,
				Stream:      true,
				Tools:       tools,
				StreamOptions: &openai.StreamOptions{
					IncludeUsage: true,
				},
			},
		)
		if err != nil {
			errChan <- fmt.Errorf("stream creation failed: %w", err)
			return
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				errChan <- fmt.Errorf("stream error: %w", err)
				return
			}

			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta
				chunk := CompletionChunk{
					Content: delta.Content,
				}

				if len(delta.ToolCalls) > 0 {
					for _, tc := range delta.ToolCalls {
						idx := 0
						if tc.Index != nil {
							idx = *tc.Index
						}
						chunk.ToolCalls = append(chunk.ToolCalls, ToolCall{
							Index: idx,
							ID:    tc.ID, // Often only present in first chunk
							Type:  string(tc.Type),
							Function: FunctionCall{
								Name:      tc.Function.Name,
								Arguments: tc.Function.Arguments,
							},
						})
					}
				}
				outputChan <- chunk
			}

			// Handle usage if present (usually in last chunk with no choices)
			if response.Usage != nil {
				outputChan <- CompletionChunk{
					Usage: &Usage{
						PromptTokens:     response.Usage.PromptTokens,
						CompletionTokens: response.Usage.CompletionTokens,
						TotalTokens:      response.Usage.TotalTokens,
					},
				}
			}
		}
	}()

	return outputChan, errChan
}

func (c *OpenAIClient) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	if model == "" {
		model = string(openai.SmallEmbedding3)
	}

	resp, err := c.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.EmbeddingModel(model),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("openai embed error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
