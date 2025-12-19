package mocks

import (
	"context"

	"github.com/hrygo/council/internal/core/workflow"
)

type WorkflowMockRepository struct {
	CreateFunc func(ctx context.Context, graph *workflow.GraphDefinition) error
	GetFunc    func(ctx context.Context, id string) (*workflow.GraphDefinition, error)
	UpdateFunc func(ctx context.Context, graph *workflow.GraphDefinition) error
	ListFunc   func(ctx context.Context) ([]*workflow.WorkflowEntity, error)
}

func (m *WorkflowMockRepository) Create(ctx context.Context, graph *workflow.GraphDefinition) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, graph)
	}
	return nil
}

func (m *WorkflowMockRepository) Get(ctx context.Context, id string) (*workflow.GraphDefinition, error) {
	if m.GetFunc != nil {
		return m.GetFunc(ctx, id)
	}
	return nil, nil
}

func (m *WorkflowMockRepository) Update(ctx context.Context, graph *workflow.GraphDefinition) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, graph)
	}
	return nil
}

func (m *WorkflowMockRepository) List(ctx context.Context) ([]*workflow.WorkflowEntity, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx)
	}
	return nil, nil
}
