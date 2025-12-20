package agent

import (
	"time"

	"github.com/google/uuid"
)

// Agent represents an AI persona.
type Agent struct {
	ID            uuid.UUID    `json:"id" db:"id"`
	Name          string       `json:"name" db:"name"`
	Avatar        *string      `json:"avatar" db:"avatar"`
	Description   *string      `json:"description" db:"description"`
	PersonaPrompt string       `json:"persona_prompt" db:"persona_prompt"`
	ModelConfig   ModelConfig  `json:"model_config" db:"model_config"`
	Capabilities  Capabilities `json:"capabilities" db:"capabilities"`
	CreatedAt     time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at" db:"updated_at"`
}

// ModelConfig defines the specific LLM configuration for this agent.
type ModelConfig struct {
	Provider    string  `json:"provider"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p"`
	MaxTokens   int     `json:"max_tokens"`
}

// Capabilities defines what this agent can do.
type Capabilities struct {
	WebSearch      bool   `json:"web_search"`
	SearchProvider string `json:"search_provider"` // e.g., "tavily"
	CodeExecution  bool   `json:"code_execution"`
}
