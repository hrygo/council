package config

import (
	"os"
	"testing"

	"github.com/hrygo/dialecta/internal/llm"
)

func TestNew(t *testing.T) {
	cfg := New()

	if cfg == nil {
		t.Fatal("New() returned nil")
	}

	// Check ProRole defaults
	if cfg.ProRole.Provider != llm.ProviderDeepSeek {
		t.Errorf("ProRole.Provider = %v, want %v", cfg.ProRole.Provider, llm.ProviderDeepSeek)
	}
	if cfg.ProRole.Model != "deepseek-chat" {
		t.Errorf("ProRole.Model = %v, want %v", cfg.ProRole.Model, "deepseek-chat")
	}
	if cfg.ProRole.Temperature != 0.8 {
		t.Errorf("ProRole.Temperature = %v, want %v", cfg.ProRole.Temperature, 0.8)
	}
	if cfg.ProRole.MaxTokens != 4096 {
		t.Errorf("ProRole.MaxTokens = %v, want %v", cfg.ProRole.MaxTokens, 4096)
	}

	// Check ConRole defaults
	if cfg.ConRole.Provider != llm.ProviderDashScope {
		t.Errorf("ConRole.Provider = %v, want %v", cfg.ConRole.Provider, llm.ProviderDashScope)
	}
	if cfg.ConRole.Model != "qwen-plus" {
		t.Errorf("ConRole.Model = %v, want %v", cfg.ConRole.Model, "qwen-plus")
	}

	// Check JudgeRole defaults
	if cfg.JudgeRole.Provider != llm.ProviderGemini {
		t.Errorf("JudgeRole.Provider = %v, want %v", cfg.JudgeRole.Provider, llm.ProviderGemini)
	}
	if cfg.JudgeRole.Model != "gemini-3-pro-preview" {
		t.Errorf("JudgeRole.Model = %v, want %v", cfg.JudgeRole.Model, "gemini-3-pro-preview")
	}
}

func TestRoleConfig_ToLLMConfig(t *testing.T) {
	role := RoleConfig{
		Provider:    llm.ProviderDeepSeek,
		Model:       "test-model",
		Temperature: 0.5,
		MaxTokens:   2048,
	}

	cfg := role.ToLLMConfig()

	if cfg.Provider != role.Provider {
		t.Errorf("Provider = %v, want %v", cfg.Provider, role.Provider)
	}
	if cfg.Model != role.Model {
		t.Errorf("Model = %v, want %v", cfg.Model, role.Model)
	}
	if cfg.Temperature != role.Temperature {
		t.Errorf("Temperature = %v, want %v", cfg.Temperature, role.Temperature)
	}
	if cfg.MaxTokens != role.MaxTokens {
		t.Errorf("MaxTokens = %v, want %v", cfg.MaxTokens, role.MaxTokens)
	}
}

func TestConfigError_Error(t *testing.T) {
	err := ConfigError("test error message")
	if err.Error() != "test error message" {
		t.Errorf("Error() = %v, want %v", err.Error(), "test error message")
	}
}

func TestConfig_Validate(t *testing.T) {
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
		name        string
		setup       func()
		cfg         *Config
		wantErr     bool
		errContains string
	}{
		{
			name: "all keys set - valid",
			setup: func() {
				os.Setenv("DEEPSEEK_API_KEY", "test-key")
				os.Setenv("GEMINI_API_KEY", "test-key")
				os.Setenv("DASHSCOPE_API_KEY", "test-key")
			},
			cfg:     New(),
			wantErr: false,
		},
		{
			name: "GOOGLE_API_KEY works for Gemini",
			setup: func() {
				os.Setenv("DEEPSEEK_API_KEY", "test-key")
				os.Unsetenv("GEMINI_API_KEY")
				os.Setenv("GOOGLE_API_KEY", "test-key")
				os.Setenv("DASHSCOPE_API_KEY", "test-key")
			},
			cfg:     New(),
			wantErr: false,
		},
		{
			name: "missing DeepSeek key",
			setup: func() {
				os.Unsetenv("DEEPSEEK_API_KEY")
				os.Setenv("GEMINI_API_KEY", "test-key")
				os.Setenv("DASHSCOPE_API_KEY", "test-key")
			},
			cfg:         New(),
			wantErr:     true,
			errContains: "DEEPSEEK_API_KEY",
		},
		{
			name: "missing Gemini key",
			setup: func() {
				os.Setenv("DEEPSEEK_API_KEY", "test-key")
				os.Unsetenv("GEMINI_API_KEY")
				os.Unsetenv("GOOGLE_API_KEY")
				os.Setenv("DASHSCOPE_API_KEY", "test-key")
			},
			cfg:         New(),
			wantErr:     true,
			errContains: "GEMINI_API_KEY",
		},
		{
			name: "missing DashScope key",
			setup: func() {
				os.Setenv("DEEPSEEK_API_KEY", "test-key")
				os.Setenv("GEMINI_API_KEY", "test-key")
				os.Unsetenv("DASHSCOPE_API_KEY")
			},
			cfg:         New(),
			wantErr:     true,
			errContains: "DASHSCOPE_API_KEY",
		},
		{
			name: "only DeepSeek provider - only needs DeepSeek key",
			setup: func() {
				os.Setenv("DEEPSEEK_API_KEY", "test-key")
				os.Unsetenv("GEMINI_API_KEY")
				os.Unsetenv("GOOGLE_API_KEY")
				os.Unsetenv("DASHSCOPE_API_KEY")
			},
			cfg: &Config{
				ProRole:   RoleConfig{Provider: llm.ProviderDeepSeek},
				ConRole:   RoleConfig{Provider: llm.ProviderDeepSeek},
				JudgeRole: RoleConfig{Provider: llm.ProviderDeepSeek},
			},
			wantErr: false,
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

			err := tt.cfg.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() error = nil, want error containing %q", tt.errContains)
				} else if tt.errContains != "" && err.Error() != tt.errContains && !contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %q", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() error = %v, want nil", err)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
