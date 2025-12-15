package group

import (
	"context"

	"github.com/google/uuid"
)

// Repository defines the interface for group persistence.
type Repository interface {
	Create(ctx context.Context, group *Group) error
	GetByID(ctx context.Context, id uuid.UUID) (*Group, error)
	List(ctx context.Context) ([]*Group, error)
	Update(ctx context.Context, group *Group) error
	Delete(ctx context.Context, id uuid.UUID) error
}
