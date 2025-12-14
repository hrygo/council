package llm

import (
	"os"
	"testing"
)

func TestParseProvider(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Provider
		wantErr bool
	}{
		{
			name:    "deepseek",
			input:   "deepseek",
			want:    ProviderDeepSeek,
			wantErr: false,
		},
		{
			name:    "gemini",
			input:   "gemini",
			want:    ProviderGemini,
			wantErr: false,
		},
		{
			name:    "google alias for gemini",
			input:   "google",
			want:    ProviderGemini,
			wantErr: false,
		},
		{
			name:    "dashscope",
			input:   "dashscope",
			want:    ProviderDashScope,
			wantErr: false,
		},
		{
			name:    "qwen alias for dashscope",
			input:   "qwen",
			want:    ProviderDashScope,
			wantErr: false,
		},
		{
			name:    "alibaba alias for dashscope",
			input:   "alibaba",
			want:    ProviderDashScope,
			wantErr: false,
		},
		{
			name:    "unknown provider",
			input:   "unknown",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseProvider(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseProvider(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseProvider(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	// Save original env vars
	origDeepSeek := os.Getenv("DEEPSEEK_API_KEY")
	origGemini := os.Getenv("GEMINI_API_KEY")
	origGoogle := os.Getenv("GOOGLE_API_KEY")
	origDashScope := os.Getenv("DASHSCOPE_API_KEY")
	defer func() {
		os.Setenv("DEEPSEEK_API_KEY", origDeepSeek)
		os.Setenv("GEMINI_API_KEY", origGemini)
		os.Setenv("GOOGLE_API_KEY", origGoogle)
		os.Setenv("DASHSCOPE_API_KEY", origDashScope)
	}()

	tests := []struct {
		name    string
		setup   func()
		cfg     Config
		wantErr bool
	}{
		{
			name: "DeepSeek with key",
			setup: func() {
				os.Setenv("DEEPSEEK_API_KEY", "test-key")
			},
			cfg: Config{
				Provider: ProviderDeepSeek,
				Model:    "deepseek-chat",
			},
			wantErr: false,
		},
		{
			name: "DeepSeek without key",
			setup: func() {
				os.Unsetenv("DEEPSEEK_API_KEY")
			},
			cfg: Config{
				Provider: ProviderDeepSeek,
			},
			wantErr: true,
		},
		{
			name: "Gemini with GEMINI_API_KEY",
			setup: func() {
				os.Setenv("GEMINI_API_KEY", "test-key")
				os.Unsetenv("GOOGLE_API_KEY")
			},
			cfg: Config{
				Provider: ProviderGemini,
				Model:    "gemini-pro",
			},
			wantErr: false,
		},
		{
			name: "Gemini with GOOGLE_API_KEY fallback",
			setup: func() {
				os.Unsetenv("GEMINI_API_KEY")
				os.Setenv("GOOGLE_API_KEY", "test-key")
			},
			cfg: Config{
				Provider: ProviderGemini,
				Model:    "gemini-pro",
			},
			wantErr: false,
		},
		{
			name: "Gemini without any key",
			setup: func() {
				os.Unsetenv("GEMINI_API_KEY")
				os.Unsetenv("GOOGLE_API_KEY")
			},
			cfg: Config{
				Provider: ProviderGemini,
			},
			wantErr: true,
		},
		{
			name: "DashScope with key",
			setup: func() {
				os.Setenv("DASHSCOPE_API_KEY", "test-key")
			},
			cfg: Config{
				Provider: ProviderDashScope,
				Model:    "qwen-plus",
			},
			wantErr: false,
		},
		{
			name: "DashScope without key",
			setup: func() {
				os.Unsetenv("DASHSCOPE_API_KEY")
			},
			cfg: Config{
				Provider: ProviderDashScope,
			},
			wantErr: true,
		},
		{
			name:  "Unknown provider",
			setup: func() {},
			cfg: Config{
				Provider: "unknown",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all env vars first
			os.Unsetenv("DEEPSEEK_API_KEY")
			os.Unsetenv("GEMINI_API_KEY")
			os.Unsetenv("GOOGLE_API_KEY")
			os.Unsetenv("DASHSCOPE_API_KEY")

			tt.setup()

			got, err := NewClient(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewClient() returned nil client without error")
			}
		})
	}
}

func TestProviderConstants(t *testing.T) {
	// Ensure provider constants have expected values
	if ProviderDeepSeek != "deepseek" {
		t.Errorf("ProviderDeepSeek = %v, want %v", ProviderDeepSeek, "deepseek")
	}
	if ProviderGemini != "gemini" {
		t.Errorf("ProviderGemini = %v, want %v", ProviderGemini, "gemini")
	}
	if ProviderDashScope != "dashscope" {
		t.Errorf("ProviderDashScope = %v, want %v", ProviderDashScope, "dashscope")
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		Role:    "user",
		Content: "test content",
	}

	if msg.Role != "user" {
		t.Errorf("Message.Role = %v, want %v", msg.Role, "user")
	}
	if msg.Content != "test content" {
		t.Errorf("Message.Content = %v, want %v", msg.Content, "test content")
	}
}

func TestConfig(t *testing.T) {
	cfg := Config{
		Provider:    ProviderDeepSeek,
		Model:       "test-model",
		Temperature: 0.7,
		MaxTokens:   1024,
	}

	if cfg.Provider != ProviderDeepSeek {
		t.Errorf("Config.Provider = %v, want %v", cfg.Provider, ProviderDeepSeek)
	}
	if cfg.Model != "test-model" {
		t.Errorf("Config.Model = %v, want %v", cfg.Model, "test-model")
	}
	if cfg.Temperature != 0.7 {
		t.Errorf("Config.Temperature = %v, want %v", cfg.Temperature, 0.7)
	}
	if cfg.MaxTokens != 1024 {
		t.Errorf("Config.MaxTokens = %v, want %v", cfg.MaxTokens, 1024)
	}
}
