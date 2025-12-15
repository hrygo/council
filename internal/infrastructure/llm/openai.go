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
	if req.Model == "" {
		req.Model = openai.GPT4o
	}

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

func (c *OpenAIClient) Stream(ctx context.Context, req *CompletionRequest) (<-chan string, <-chan error) {
	// TODO: Implement streaming
	outputChan := make(chan string)
	errChan := make(chan error, 1)

	go func() {
		defer close(outputChan)
		defer close(errChan)

		if req.Model == "" {
			req.Model = openai.GPT4o
		}

		messages := make([]openai.ChatCompletionMessage, len(req.Messages))
		for i, msg := range req.Messages {
			messages[i] = openai.ChatCompletionMessage{
				Role:    msg.Role,
				Content: msg.Content,
			}
		}

		stream, err := c.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model:       req.Model,
				Messages:    messages,
				Temperature: req.Temperature,
				MaxTokens:   req.MaxTokens,
				Stream:      true,
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
				outputChan <- response.Choices[0].Delta.Content
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
