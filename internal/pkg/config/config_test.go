package config

import (
	"os"
	"testing"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear relevant env vars
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("EMBEDDING_PROVIDER")
	os.Unsetenv("EMBEDDING_MODEL")

	cfg := Load()

	if cfg.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", cfg.Port)
	}
	if cfg.Embedding.Provider != "siliconflow" {
		t.Errorf("Expected default siliconflow, got %s", cfg.Embedding.Provider)
	}
}

func TestLoad_EmbeddingFallbacks(t *testing.T) {
	tests := []struct {
		provider string
		expected string
	}{
		{"openai", "text-embedding-3-small"},
		{"dashscope", "text-embedding-v1"},
		{"siliconflow", "Qwen/Qwen3-Embedding-8B"},
		{"ollama", "gte-qwen2-1.5b-instruct-embed-f16"},
		{"gemini", "text-embedding-004"},
	}

	for _, tt := range tests {
		os.Setenv("EMBEDDING_PROVIDER", tt.provider)
		os.Unsetenv("EMBEDDING_MODEL")
		cfg := Load()
		if cfg.Embedding.Model != tt.expected {
			t.Errorf("For %s, expected model %s, got %s", tt.provider, tt.expected, cfg.Embedding.Model)
		}
	}
}
