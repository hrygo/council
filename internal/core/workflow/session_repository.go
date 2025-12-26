package workflow

import (
	"context"
)

// SessionEntity represents the persistent state of a session.
type SessionEntity struct {
	ID         string                 `json:"session_uuid" db:"session_uuid"`
	GroupID    string                 `json:"group_uuid" db:"group_uuid"`
	WorkflowID string                 `json:"workflow_uuid" db:"workflow_uuid"`
	Status     SessionStatus          `json:"status"`
	Proposal   map[string]interface{} `json:"proposal"`
	StartedAt  *interface{}           `json:"started_at"` // Simplified for now
	EndedAt    *interface{}           `json:"ended_at"`
}

// SessionRepository defines the interface for session persistence.
type SessionRepository interface {
	Create(ctx context.Context, session *Session, groupID string, workflowID string) error
	Get(ctx context.Context, id string) (*SessionEntity, error)
	UpdateStatus(ctx context.Context, id string, status SessionStatus) error
}
