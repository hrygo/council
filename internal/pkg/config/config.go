package config

import (
	"os"
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

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/council?sslmode=disable"
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	// 1536 Dimension Embedding Models Recommendation
	// OpenAI:      text-embedding-3-small
	// DashScope:   text-embedding-v1
	// SiliconFlow: Qwen/Qwen3-Embedding-8B
	// Ollama:      gte-qwen2-1.5b-instruct-embed-f16

	embeddingProvider := os.Getenv("EMBEDDING_PROVIDER")
	if embeddingProvider == "" {
		embeddingProvider = "siliconflow"
	}

	embeddingModel := os.Getenv("EMBEDDING_MODEL")
	if embeddingModel == "" {
		switch embeddingProvider {
		case "openai":
			embeddingModel = "text-embedding-3-small"
		case "dashscope":
			embeddingModel = "text-embedding-v1"
		case "siliconflow":
			embeddingModel = "Qwen/Qwen3-Embedding-8B"
		case "ollama":
			// Note: User must pull this model manually in Ollama
			embeddingModel = "gte-qwen2-1.5b-instruct-embed-f16"
		case "gemini":
			// Warning: Gemini defaults to 768.
			// Users should use gemini-embedding-001 with explicit config if supported,
			// or accept dimensionality mismatch (store as 1536 pad/truncate?).
			// For now, default to text-embedding-004 as closest standard.
			embeddingModel = "text-embedding-004"
		}
	}

	embeddingKey := os.Getenv("EMBEDDING_API_KEY")
	if embeddingKey == "" {
		switch embeddingProvider {
		case "openai":
			embeddingKey = os.Getenv("OPENAI_API_KEY")
		case "dashscope":
			embeddingKey = os.Getenv("DASHSCOPE_API_KEY")
		case "siliconflow":
			embeddingKey = os.Getenv("SILICONFLOW_API_KEY")
		case "gemini":
			embeddingKey = os.Getenv("GEMINI_API_KEY")
		}
	}

	return &Config{
		Port:         port,
		DatabaseURL:  dbURL,
		RedisURL:     redisURL,
		TavilyAPIKey: os.Getenv("TAVILY_API_KEY"),

		LLM: LLMConfig{
			Provider: os.Getenv("LLM_PROVIDER"),
			APIKey:   os.Getenv("LLM_API_KEY"),
			BaseURL:  os.Getenv("LLM_BASE_URL"),
			Model:    os.Getenv("LLM_MODEL"),
		},
		Embedding: EmbeddingConfig{
			Provider: embeddingProvider,
			APIKey:   embeddingKey,
			BaseURL:  os.Getenv("EMBEDDING_BASE_URL"),
			Model:    embeddingModel,
		},

		OpenAIKey:      os.Getenv("OPENAI_API_KEY"),
		DeepSeekKey:    os.Getenv("DEEPSEEK_API_KEY"),
		DashScopeKey:   os.Getenv("DASHSCOPE_API_KEY"),
		GeminiKey:      os.Getenv("GEMINI_API_KEY"),
		SiliconFlowKey: os.Getenv("SILICONFLOW_API_KEY"),
	}
}
