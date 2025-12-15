package llm

import (
	"context"
	"fmt"
)

// DeepSeekClient wraps OpenAIClient to explicitly handle unsupported methods.
type DeepSeekClient struct {
	*OpenAIClient
}

// NewDeepSeekClient returns a Client that uses DeepSeek's OpenAI-compatible API.
func NewDeepSeekClient(apiKey string) *DeepSeekClient {
	// DeepSeek API endpoint
	// https://api-docs.deepseek.com/
	baseURL := "https://api.deepseek.com"
	client := NewOpenAICompatibleClient(apiKey, baseURL)
	return &DeepSeekClient{
		OpenAIClient: client,
	}
}

// Ensure DeepSeekClient implements LLMProvider and Embedder interfaces
var _ LLMProvider = (*DeepSeekClient)(nil)
var _ Embedder = (*DeepSeekClient)(nil)

// Embed overrides the default OpenAI Embed method to return an error,
// as DeepSeek does not support embeddings.
func (c *DeepSeekClient) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	return nil, fmt.Errorf("deepseek provider does not support embeddings")
}
