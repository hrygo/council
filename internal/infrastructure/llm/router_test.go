package llm

import (
	"testing"

	"github.com/hrygo/council/internal/pkg/config"
)

func TestRegistry_NewEmbedder(t *testing.T) {
	cfg := &config.Config{}
	registry := NewRegistry(cfg)

	tests := []struct {
		name    string
		conf    EmbeddingConfig
		wantErr bool
	}{
		{"OpenAI", EmbeddingConfig{Type: "openai", APIKey: "test", Model: "text-embedding-3-small"}, false},
		{"Unknown", EmbeddingConfig{Type: "unknown"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := registry.NewEmbedder(tt.conf)
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

func TestRegistry_GetLLMProvider(t *testing.T) {
	cfg := &config.Config{
		LLM: config.LLMConfig{
			Provider: "openai",
			APIKey:   "sk-test",
			Model:    "gpt-4",
		},
		OpenAIKey: "sk-test-global",
	}
	registry := NewRegistry(cfg)

	tests := []struct {
		name         string
		providerName string
		wantErr      bool
	}{
		{"Default", "default", false},
		{"OpenAI Explicit", "openai", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := registry.GetLLMProvider(tt.providerName)
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
