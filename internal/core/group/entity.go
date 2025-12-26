package group

import (
	"time"

	"github.com/google/uuid"
)

// Group represents a collaboration group (Project/Context).
type Group struct {
	ID                uuid.UUID   `json:"group_uuid" db:"group_uuid"`
	Name              string      `json:"name" db:"name"`
	Icon              *string     `json:"icon" db:"icon"`
	SystemPrompt      *string     `json:"system_prompt" db:"system_prompt"`
	DefaultAgentUUIDs []uuid.UUID `json:"default_agent_uuids" db:"default_agent_uuids"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
}
