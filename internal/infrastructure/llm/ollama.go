package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// OllamaClient wraps OpenAIClient to handle Ollama specific logic.
type OllamaClient struct {
	*OpenAIClient
}

// NewOllamaClient returns a Client that uses Ollama's OpenAI-compatible API.
func NewOllamaClient(baseURL string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434/v1"
	}
	// "dummy" api key for Ollama
	client := NewOpenAICompatibleClient("ollama", baseURL)
	return &OllamaClient{
		OpenAIClient: client,
	}
}

// Ensure OllamaClient implements LLMProvider and Embedder interfaces
var _ LLMProvider = (*OllamaClient)(nil)
var _ Embedder = (*OllamaClient)(nil)

func (c *OllamaClient) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	if model == "" {
		// Ollama requires a model to be pulled. We can't easily guess a default that exists.
		// But 'nomic-embed-text' or 'mxbai-embed-large' are common.
		// For now, let's error if not provided to force user configuration,
		// OR default to 'mxbai-embed-large' which is popular.
		// Let's rely on user config, but avoid OpenAI default.
		return nil, fmt.Errorf("model is required for ollama embeddings")
	}

	resp, err := c.client.CreateEmbeddings(
		ctx,
		openai.EmbeddingRequest{
			Input: []string{text},
			Model: openai.EmbeddingModel(model),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("ollama embed error: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Data[0].Embedding, nil
}
