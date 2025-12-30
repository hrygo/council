package llm

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hrygo/council/internal/pkg/config"
)

// Registry manages LLM providers.
// It supports checking available providers based on configuration and system defaults.
type Registry struct {
	cfg       *config.Config
	providers map[string]LLMProvider
	mu        sync.RWMutex
}

// NewRegistry creates a new LLM provider registry
func NewRegistry(cfg *config.Config) *Registry {
	return &Registry{
		cfg:       cfg,
		providers: make(map[string]LLMProvider),
	}
}

// RegisterProvider explicitly registers a provider instance (useful for testing or custom providers)
func (r *Registry) RegisterProvider(name string, provider LLMProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[name] = provider
}

// GetLLMProvider retrieves a provider by name.
// If providerName is empty, it returns the system default provider.
// Once a provider's API key is configured, it becomes available for use.
func (r *Registry) GetLLMProvider(providerName string) (LLMProvider, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Normalize
	providerName = strings.ToLower(strings.TrimSpace(providerName))

	// 1. Check Cache explicitly for the requested name (e.g. "default" mock)
	if p, ok := r.providers[providerName]; ok {
		return p, nil
	}

	// 2. Handle "default" or empty resolution if not in cache
	if providerName == "" || providerName == "default" {
		providerName = r.cfg.LLM.Provider
		// Check cache again for resolved name
		if p, ok := r.providers[providerName]; ok {
			return p, nil
		}
	}

	// 3. Initialize New Provider
	// We need to construct a temp config to reuse existing factory methods.
	// We look up API keys from the global config based on provider name.
	apiKey := ""
	baseURL := ""

	switch providerName {
	case "openai":
		apiKey = r.cfg.OpenAIKey
	case "deepseek":
		apiKey = r.cfg.DeepSeekKey
	case "dashscope":
		apiKey = r.cfg.DashScopeKey
	case "siliconflow":
		apiKey = r.cfg.SiliconFlowKey
	case "gemini", "google":
		apiKey = r.cfg.GeminiKey
	case "ollama":
		// Ollama usually uses BaseURL, check LLM BaseURL if provider matches, or default
		if r.cfg.LLM.Provider == "ollama" {
			baseURL = r.cfg.LLM.BaseURL
		}
	default:
		// Attempt to fallback to generic LLM config if names match
		if r.cfg.LLM.Provider == providerName {
			apiKey = r.cfg.LLM.APIKey
			baseURL = r.cfg.LLM.BaseURL
		}
	}

	// If no specific key found, and we are requesting the system default provider, use its generic config
	if apiKey == "" && providerName == r.cfg.LLM.Provider {
		apiKey = r.cfg.LLM.APIKey
	}

	// Create Config wrapper
	llmCfg := LLMConfig{
		Type:    providerName,
		APIKey:  apiKey,
		BaseURL: baseURL,
	}

	provider, err := r.createProvider(llmCfg)
	if err != nil {
		return nil, err
	}

	// Cache it
	r.providers[providerName] = provider
	return provider, nil
}

// createProvider is the internal factory (formerly GetLLMProvider)
func (r *Registry) createProvider(config LLMConfig) (LLMProvider, error) {
	switch config.Type {
	case "openai":
		return NewOpenAIClient(config.APIKey), nil
	case "anthropic":
		return nil, fmt.Errorf("anthropic provider not implemented")
	case "google", "gemini":
		return NewGeminiClient(config.APIKey), nil
	case "deepseek":
		return NewDeepSeekClient(config.APIKey), nil
	case "ollama":
		return NewOllamaClient(config.BaseURL), nil
	case "dashscope":
		return NewDashScopeClient(config.APIKey), nil
	case "siliconflow":
		return NewSiliconFlowClient(config.APIKey), nil
	default:
		return nil, fmt.Errorf("unknown llm provider type: %s", config.Type)
	}
}

// GetDefaultModel returns the system default LLM model name from config
func (r *Registry) GetDefaultModel() string {
	if r.cfg.LLM.Model != "" {
		return r.cfg.LLM.Model
	}
	// Fallback based on provider
	switch r.cfg.LLM.Provider {
	case "gemini", "google":
		return "gemini-2.0-flash"
	case "openai":
		if r.cfg.LLM.Model != "" {
			return r.cfg.LLM.Model
		}
		return "gpt-4o"
	case "deepseek":
		return "deepseek-chat"
	case "dashscope":
		return "qwen-max"
	default:
		if r.cfg.LLM.Model != "" {
			return r.cfg.LLM.Model
		}
		return "gpt-4o"
	}
}

// NewEmbedder creates a new Embedder based on embedding config.
// Embedder configuration is usually stricter (must match vector DB), so we keep it separate from dynamic LLM registry for now.
func (r *Registry) NewEmbedder(config EmbeddingConfig) (Embedder, error) {
	// Re-use logic or keep simple.
	// Embedding usually doesn't need "dynamic switching" per agent,
	// but the Client implementations are often the same.
	// For now, keep as is but use the factory methods.

	switch config.Type {
	case "openai":
		return NewOpenAIClient(config.APIKey), nil
	case "google", "gemini":
		return NewGeminiClient(config.APIKey), nil
	case "ollama":
		return NewOllamaClient(config.BaseURL), nil
	case "deepseek":
		return NewDeepSeekClient(config.APIKey), nil
	case "dashscope":
		return NewDashScopeClient(config.APIKey), nil
	case "siliconflow":
		return NewSiliconFlowClient(config.APIKey), nil
	default:
		return nil, fmt.Errorf("unknown embedding provider type: %s", config.Type)
	}
}
