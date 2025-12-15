package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// SiliconFlowClient wraps OpenAIClient to handle SiliconFlow specific logic.
type SiliconFlowClient struct {
	*OpenAIClient
}

// NewSiliconFlowClient returns a Client that uses SiliconFlow's OpenAI-compatible API.
func NewSiliconFlowClient(apiKey string) *SiliconFlowClient {
	// SiliconFlow API endpoint
	// https://siliconflow.cn/
	baseURL := "https://api.siliconflow.cn/v1"
	client := NewOpenAICompatibleClient(apiKey, baseURL)
	return &SiliconFlowClient{
		OpenAIClient: client,
	}
}

// Ensure SiliconFlowClient implements LLMProvider and Embedder interfaces
var _ LLMProvider = (*SiliconFlowClient)(nil)
var _ Embedder = (*SiliconFlowClient)(nil)

func (c *SiliconFlowClient) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	if model == "" {
		// SiliconFlow has popular models like BAAI/bge-m3
		model = "BAAI/bge-m3"
	}

	resp, err := c.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.EmbeddingModel(model),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("siliconflow embed error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
