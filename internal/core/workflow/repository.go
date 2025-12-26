package workflow

import (
	"context"
	"time"
)

// WorkflowEntity represents the persistent storage for a workflow.
type WorkflowEntity struct {
	ID              string          `json:"workflow_uuid" db:"workflow_uuid"`
	GroupID         string          `json:"group_uuid" db:"group_uuid"`
	Name            string          `json:"name"`
	GraphDefinition GraphDefinition `json:"graph_definition"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// Repository defines the interface for workflow persistence.
type Repository interface {
	Create(ctx context.Context, graph *GraphDefinition) error
	Get(ctx context.Context, id string) (*GraphDefinition, error)
	Update(ctx context.Context, graph *GraphDefinition) error
	List(ctx context.Context) ([]*WorkflowEntity, error)
}
