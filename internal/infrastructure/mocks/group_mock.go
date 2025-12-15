package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/group"
)

// GroupMockRepository implements Repository for testing.
type GroupMockRepository struct {
	Groups map[uuid.UUID]*group.Group
	Err    error
}

func NewGroupMockRepository() *GroupMockRepository {
	return &GroupMockRepository{
		Groups: make(map[uuid.UUID]*group.Group),
	}
}

func (m *GroupMockRepository) Create(ctx context.Context, g *group.Group) error {
	if m.Err != nil {
		return m.Err
	}
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	m.Groups[g.ID] = g
	return nil
}

func (m *GroupMockRepository) GetByID(ctx context.Context, id uuid.UUID) (*group.Group, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	g, ok := m.Groups[id]
	if !ok {
		return nil, context.DeadlineExceeded
	}
	return g, nil
}

func (m *GroupMockRepository) List(ctx context.Context) ([]*group.Group, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var list []*group.Group
	for _, g := range m.Groups {
		list = append(list, g)
	}
	return list, nil
}

func (m *GroupMockRepository) Update(ctx context.Context, g *group.Group) error {
	if m.Err != nil {
		return m.Err
	}
	if _, ok := m.Groups[g.ID]; !ok {
		return context.DeadlineExceeded
	}
	m.Groups[g.ID] = g
	return nil
}

func (m *GroupMockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.Err != nil {
		return m.Err
	}
	delete(m.Groups, id)
	return nil
}
