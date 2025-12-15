package llm

import (
	"fmt"
)

// Router creates Providers based on configuration.
type Router struct{}

// NewRouter creates a new Router.
func NewRouter() *Router {
	return &Router{}
}

// GetLLMProvider creates a new Chat Model provider based on the configuration.
func (r *Router) GetLLMProvider(config LLMConfig) (LLMProvider, error) {
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

// NewEmbedder creates a new Embedder based on the configuration.
func (r *Router) NewEmbedder(config EmbeddingConfig) (Embedder, error) {
	switch config.Type {
	case "openai":
		return NewOpenAIClient(config.APIKey), nil
	case "google", "gemini":
		return NewGeminiClient(config.APIKey), nil
	case "ollama":
		return NewOllamaClient(config.BaseURL), nil
	// DeepSeek and DashScope also use OpenAI compatible client for embeddings usually
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
