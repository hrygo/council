package llm

import (
	"context"
	"fmt"
)

// MockProvider implements both LLMProvider and Embedder for testing.
type MockProvider struct {
	// Setup expectations
	GenerateResponse      *CompletionResponse
	GenerateResponseQueue []*CompletionResponse // Pop from front
	GenerateError         error

	StreamContent []string
	StreamError   error

	EmbedResponse []float32
	EmbedError    error

	// Method call tracking
	GenerateCalls int
	StreamCalls   int
	EmbedCalls    int
}

// Ensure MockProvider implements interfaces
var _ LLMProvider = (*MockProvider)(nil)
var _ Embedder = (*MockProvider)(nil)

func NewMockProvider() *MockProvider {
	return &MockProvider{
		GenerateResponse: &CompletionResponse{
			Content: "Mocked response content",
			Usage:   Usage{PromptTokens: 10, CompletionTokens: 5, TotalTokens: 15},
		},
		EmbedResponse: make([]float32, 1536), // Default to 1536 dim
	}
}

func (m *MockProvider) Generate(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	m.GenerateCalls++
	if m.GenerateError != nil {
		return nil, m.GenerateError
	}
	if len(m.GenerateResponseQueue) > 0 {
		res := m.GenerateResponseQueue[0]
		m.GenerateResponseQueue = m.GenerateResponseQueue[1:]
		return res, nil
	}
	return m.GenerateResponse, nil
}

func (m *MockProvider) Stream(ctx context.Context, req *CompletionRequest) (<-chan CompletionChunk, <-chan error) {
	m.StreamCalls++
	// Buffer enough for: 1 content + N tool calls + 1 usage
	contentChan := make(chan CompletionChunk, 10)
	errChan := make(chan error, 1)

	go func() {
		defer close(contentChan)
		defer close(errChan)

		if m.StreamError != nil {
			errChan <- m.StreamError
			return
		}

		// Priority 1: Explicit StreamContent (for specific stream testing)
		if len(m.StreamContent) > 0 {
			for _, chunk := range m.StreamContent {
				contentChan <- CompletionChunk{Content: chunk}
			}
			return
		}

		// Priority 2: GenerateResponseQueue or GenerateResponse
		var resp *CompletionResponse
		if len(m.GenerateResponseQueue) > 0 {
			resp = m.GenerateResponseQueue[0]
			m.GenerateResponseQueue = m.GenerateResponseQueue[1:]
		} else {
			resp = m.GenerateResponse
		}

		if resp == nil {
			return
		}

		// Convert Response to Chunks
		if resp.Content != "" {
			contentChan <- CompletionChunk{Content: resp.Content}
		}

		for i, tc := range resp.ToolCalls {
			tcCopy := tc
			tcCopy.Index = i // Ensure Index is set explicitly for streaming aggregation
			contentChan <- CompletionChunk{
				ToolCalls: []ToolCall{tcCopy},
			}
		}

		if resp.Usage.TotalTokens > 0 {
			usageCopy := resp.Usage
			contentChan <- CompletionChunk{Usage: &usageCopy}
		}
	}()

	return contentChan, errChan
}

func (m *MockProvider) Embed(ctx context.Context, model string, text string) ([]float32, error) {
	m.EmbedCalls++
	if m.EmbedError != nil {
		return nil, m.EmbedError
	}
	if len(m.EmbedResponse) == 0 {
		return nil, fmt.Errorf("mock embeddings empty")
	}
	return m.EmbedResponse, nil
}
