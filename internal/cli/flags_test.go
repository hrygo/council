package cli

import (
	"testing"

	"github.com/hrygo/dialecta/internal/config"
	"github.com/hrygo/dialecta/internal/llm"
)

func TestOptions_ApplyToConfig(t *testing.T) {
	tests := []struct {
		name           string
		opts           *Options
		wantProProv    llm.Provider
		wantConProv    llm.Provider
		wantJudgeProv  llm.Provider
		wantProModel   string
		wantConModel   string
		wantJudgeModel string
	}{
		{
			name: "apply all options",
			opts: &Options{
				ProProvider:   "gemini",
				ProModel:      "custom-pro",
				ConProvider:   "deepseek",
				ConModel:      "custom-con",
				JudgeProvider: "dashscope",
				JudgeModel:    "custom-judge",
			},
			wantProProv:    llm.ProviderGemini,
			wantConProv:    llm.ProviderDeepSeek,
			wantJudgeProv:  llm.ProviderDashScope,
			wantProModel:   "custom-pro",
			wantConModel:   "custom-con",
			wantJudgeModel: "custom-judge",
		},
		{
			name: "empty models - keep defaults",
			opts: &Options{
				ProProvider:   "deepseek",
				ConProvider:   "dashscope",
				JudgeProvider: "gemini",
			},
			wantProProv:   llm.ProviderDeepSeek,
			wantConProv:   llm.ProviderDashScope,
			wantJudgeProv: llm.ProviderGemini,
			// Models should remain as defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.New()
			tt.opts.ApplyToConfig(cfg)

			if cfg.ProRole.Provider != tt.wantProProv {
				t.Errorf("ProRole.Provider = %v, want %v", cfg.ProRole.Provider, tt.wantProProv)
			}
			if cfg.ConRole.Provider != tt.wantConProv {
				t.Errorf("ConRole.Provider = %v, want %v", cfg.ConRole.Provider, tt.wantConProv)
			}
			if cfg.JudgeRole.Provider != tt.wantJudgeProv {
				t.Errorf("JudgeRole.Provider = %v, want %v", cfg.JudgeRole.Provider, tt.wantJudgeProv)
			}

			if tt.wantProModel != "" && cfg.ProRole.Model != tt.wantProModel {
				t.Errorf("ProRole.Model = %v, want %v", cfg.ProRole.Model, tt.wantProModel)
			}
			if tt.wantConModel != "" && cfg.ConRole.Model != tt.wantConModel {
				t.Errorf("ConRole.Model = %v, want %v", cfg.ConRole.Model, tt.wantConModel)
			}
			if tt.wantJudgeModel != "" && cfg.JudgeRole.Model != tt.wantJudgeModel {
				t.Errorf("JudgeRole.Model = %v, want %v", cfg.JudgeRole.Model, tt.wantJudgeModel)
			}
		})
	}
}

func TestOptions_NeedsHelp(t *testing.T) {
	tests := []struct {
		name string
		opts *Options
		want bool
	}{
		{
			name: "no source, not interactive",
			opts: &Options{
				Source:      "",
				Interactive: false,
			},
			want: true,
		},
		{
			name: "has source",
			opts: &Options{
				Source:      "file.txt",
				Interactive: false,
			},
			want: false,
		},
		{
			name: "stdin source",
			opts: &Options{
				Source:      "-",
				Interactive: false,
			},
			want: false,
		},
		{
			name: "interactive mode",
			opts: &Options{
				Source:      "",
				Interactive: true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.NeedsHelp(); got != tt.want {
				t.Errorf("NeedsHelp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_Defaults(t *testing.T) {
	opts := &Options{
		ProProvider:   "deepseek",
		ConProvider:   "dashscope",
		JudgeProvider: "gemini",
		Stream:        true,
	}

	// Verify default provider strings
	if opts.ProProvider != "deepseek" {
		t.Errorf("Default ProProvider = %v, want deepseek", opts.ProProvider)
	}
	if opts.ConProvider != "dashscope" {
		t.Errorf("Default ConProvider = %v, want dashscope", opts.ConProvider)
	}
	if opts.JudgeProvider != "gemini" {
		t.Errorf("Default JudgeProvider = %v, want gemini", opts.JudgeProvider)
	}
	if !opts.Stream {
		t.Error("Default Stream should be true")
	}
}

func TestOptions_ApplyToConfig_InvalidProvider(t *testing.T) {
	opts := &Options{
		ProProvider:   "invalid",
		ConProvider:   "invalid",
		JudgeProvider: "invalid",
	}

	cfg := config.New()
	originalProProv := cfg.ProRole.Provider
	originalConProv := cfg.ConRole.Provider
	originalJudgeProv := cfg.JudgeRole.Provider

	opts.ApplyToConfig(cfg)

	// Invalid providers should not change the config
	if cfg.ProRole.Provider != originalProProv {
		t.Errorf("Invalid provider should not change ProRole.Provider")
	}
	if cfg.ConRole.Provider != originalConProv {
		t.Errorf("Invalid provider should not change ConRole.Provider")
	}
	if cfg.JudgeRole.Provider != originalJudgeProv {
		t.Errorf("Invalid provider should not change JudgeRole.Provider")
	}
}
