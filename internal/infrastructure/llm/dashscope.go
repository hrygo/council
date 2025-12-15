package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// DashScopeClient wraps OpenAIClient to handle DashScope specific logic.
type DashScopeClient struct {
	*OpenAIClient
}

// NewDashScopeClient returns a Client that uses DashScope's OpenAI-compatible API.
func NewDashScopeClient(apiKey string) *DashScopeClient {
	// DashScope compatible endpoint
	// https://help.aliyun.com/zh/model-studio/developer-reference/use-openai-python-sdk
	baseURL := "https://dashscope.aliyuncs.com/compatible-mode/v1"
	client := NewOpenAICompatibleClient(apiKey, baseURL)
	return &DashScopeClient{
		OpenAIClient: client,
	}
}

// Ensure DashScopeClient implements LLMProvider and Embedder interfaces
var _ LLMProvider = (*DashScopeClient)(nil)
var _ Embedder = (*DashScopeClient)(nil)

func (c *DashScopeClient) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	if model == "" {
		model = "text-embedding-v1" // Common DashScope embedding model
	}

	resp, err := c.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.EmbeddingModel(model),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("dashscope embed error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
