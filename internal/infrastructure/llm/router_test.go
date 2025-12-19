package llm

import (
	"testing"
)

func TestRouter_GetLLMProvider(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		name    string
		config  LLMConfig
		wantErr bool
	}{
		{"OpenAI", LLMConfig{Type: "openai", APIKey: "test"}, false},
		{"Gemini", LLMConfig{Type: "gemini", APIKey: "test"}, false},
		{"Google", LLMConfig{Type: "google", APIKey: "test"}, false},
		{"DeepSeek", LLMConfig{Type: "deepseek", APIKey: "test"}, false},
		{"Ollama", LLMConfig{Type: "ollama", BaseURL: "http://localhost"}, false},
		{"DashScope", LLMConfig{Type: "dashscope", APIKey: "test"}, false},
		{"SiliconFlow", LLMConfig{Type: "siliconflow", APIKey: "test"}, false},
		{"Unknown", LLMConfig{Type: "unknown"}, true},
		{"Anthropic Not Impl", LLMConfig{Type: "anthropic"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := router.GetLLMProvider(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLLMProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && p == nil {
				t.Error("GetLLMProvider() returned nil provider without error")
			}
		})
	}
}

func TestRouter_NewEmbedder(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		name    string
		config  EmbeddingConfig
		wantErr bool
	}{
		{"OpenAI", EmbeddingConfig{Type: "openai", APIKey: "test"}, false},
		{"Gemini", EmbeddingConfig{Type: "gemini", APIKey: "test"}, false},
		{"Ollama", EmbeddingConfig{Type: "ollama", BaseURL: "http://localhost"}, false},
		{"DeepSeek", EmbeddingConfig{Type: "deepseek", APIKey: "test"}, false},
		{"Unknown", EmbeddingConfig{Type: "unknown"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := router.NewEmbedder(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEmbedder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && p == nil {
				t.Error("NewEmbedder() returned nil embedder without error")
			}
		})
	}
}
