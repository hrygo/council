package llm

import (
	"context"
)

// Message represents a chat message.
type Message struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"` // For tool_result messages
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"`
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// CompletionRequest represents a request to the LLM.
type CompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature"`
	TopP        float32   `json:"top_p,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
	Tools       []Tool    `json:"tools,omitempty"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"` // JSON Schema
}

// CompletionResponse represents the response from the LLM.
type CompletionResponse struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
	Usage     Usage      `json:"usage"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// LLMProvider defines the interface for a Chat Model provider.
type LLMProvider interface {
	// Generate returns a complete response.
	Generate(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)
	// Stream returns a channel of chunks.
	Stream(ctx context.Context, req *CompletionRequest) (<-chan string, <-chan error)
}

// Embedder defines the interface for an Embedding provider.
type Embedder interface {
	// Embed generates embeddings for the given text using the specified model.
	Embed(ctx context.Context, model string, text string) ([]float32, error)
}

// LLMConfig holds configuration for creating an LLM provider.
type LLMConfig struct {
	Type    string // openai, anthropic, google, deepseek, ollama, dashscope
	APIKey  string
	BaseURL string
	Model   string // default model
}

// EmbeddingConfig holds configuration for creating an Embedder.
type EmbeddingConfig struct {
	Type    string // openai, google, ollama, etc.
	APIKey  string
	BaseURL string
	Model   string // embedding model name
}
