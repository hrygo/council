package agent

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for agent persistence.
type Repository interface {
	Create(ctx context.Context, agent *Agent) error
	GetByID(ctx context.Context, id uuid.UUID) (*Agent, error)
	List(ctx context.Context) ([]*Agent, error)
	Update(ctx context.Context, agent *Agent) error
	Delete(ctx context.Context, id uuid.UUID) error
}
