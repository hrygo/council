package config

import (
	"os"
	"strings"
)

type Config struct {
	Port           string
	DatabaseURL    string
	TavilyAPIKey   string
	OpenAIKey      string
	DeepSeekKey    string
	DashScopeKey   string
	GeminiKey      string
	SiliconFlowKey string
	RedisURL       string

	LLM       LLMConfig
	Embedding EmbeddingConfig
}

// LLMConfig defines the System Default Chat Model configuration.
// This is used for:
// 1. System-level tasks (e.g., Wizard Mode, background analysis).
// 2. Fallback when an Agent does not have a specific model configured.
// NOTE: Individual Agents (in DB) can and should have their own independent model configurations.
type LLMConfig struct {
	Provider string // e.g., "siliconflow", "openai"
	APIKey   string // Optional override. Usually defaults to the global key for the chosen provider.
	BaseURL  string // Used primarily for "ollama" or custom OpenAI-compatible proxies (OneAPI).
	Model    string // e.g., "deepseek-v3", "gpt-4o"
}

// EmbeddingConfig defines the System Wide Embedding configuration.
// Unlike Chat Models, Embeddings should generally be consistent across the entire system
// to ensure vector compatibility in the "memories" table (VECTOR column).
type EmbeddingConfig struct {
	Provider string // e.g., "siliconflow", "openai"
	APIKey   string // Optional override.
	BaseURL  string // Used for "ollama" or custom endpoints.
	Model    string // Critical: Must match the database VECTOR dimension (e.g., 1536).
}

const (
	DefaultPort        = "8080"
	DefaultDBURL       = "postgres://user:password@localhost:5432/council?sslmode=disable"
	DefaultRedisURL    = "localhost:6379"
	DefaultEmbedding   = "siliconflow"
	DefaultVectorModel = "Qwen/Qwen3-Embedding-8B"
)

func Load() *Config {
	cfg := &Config{
		Port:        getEnv("PORT", DefaultPort),
		DatabaseURL: getEnv("DATABASE_URL", DefaultDBURL),
		RedisURL:    getEnv("REDIS_URL", DefaultRedisURL),
	}

	// LLM Config
	cfg.LLM = LLMConfig{
		Provider: getEnv("LLM_PROVIDER", "siliconflow"),
		APIKey:   os.Getenv("LLM_API_KEY"),
		BaseURL:  os.Getenv("LLM_BASE_URL"),
		Model:    getEnv("LLM_MODEL", "deepseek-v3"),
	}

	// Embedding Config
	provider := getEnv("EMBEDDING_PROVIDER", DefaultEmbedding)
	cfg.Embedding = EmbeddingConfig{
		Provider: provider,
		APIKey:   getEnv("EMBEDDING_API_KEY", getProviderKey(provider)),
		BaseURL:  os.Getenv("EMBEDDING_BASE_URL"),
		Model:    getEnv("EMBEDDING_MODEL", getDefaultForProvider(provider)),
	}

	// Legacy keys mapping
	cfg.TavilyAPIKey = os.Getenv("TAVILY_API_KEY")
	cfg.OpenAIKey = os.Getenv("OPENAI_API_KEY")
	cfg.DeepSeekKey = os.Getenv("DEEPSEEK_API_KEY")
	cfg.DashScopeKey = os.Getenv("DASHSCOPE_API_KEY")
	cfg.GeminiKey = os.Getenv("GEMINI_API_KEY")
	cfg.SiliconFlowKey = os.Getenv("SILICONFLOW_API_KEY")

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func getProviderKey(provider string) string {
	switch provider {
	case "openai":
		return os.Getenv("OPENAI_API_KEY")
	case "dashscope":
		return os.Getenv("DASHSCOPE_API_KEY")
	case "siliconflow":
		return os.Getenv("SILICONFLOW_API_KEY")
	case "gemini":
		return os.Getenv("GEMINI_API_KEY")
	default:
		return ""
	}
}

func getDefaultForProvider(provider string) string {
	switch provider {
	case "openai":
		return "text-embedding-3-small"
	case "dashscope":
		return "text-embedding-v1"
	case "siliconflow":
		return "Qwen/Qwen3-Embedding-8B"
	case "ollama":
		return "gte-qwen2-1.5b-instruct-embed-f16"
	case "gemini":
		return "text-embedding-004"
	default:
		return DefaultVectorModel
	}
}
