package mocks

import (
	"context"

	"github.com/hrygo/council/internal/core/workflow"
)

type TemplateMockRepository struct {
	Templates []workflow.Template
	Err       error
}

func (m *TemplateMockRepository) List(ctx context.Context) ([]workflow.Template, error) {
	return m.Templates, m.Err
}

func (m *TemplateMockRepository) Create(ctx context.Context, t *workflow.Template) error {
	if m.Err != nil {
		return m.Err
	}
	m.Templates = append(m.Templates, *t)
	return nil
}

func (m *TemplateMockRepository) Get(ctx context.Context, id string) (*workflow.Template, error) {
	for _, t := range m.Templates {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, m.Err
}

func (m *TemplateMockRepository) Delete(ctx context.Context, id string) error {
	if m.Err != nil {
		return m.Err
	}
	for i, t := range m.Templates {
		if t.ID == id {
			m.Templates = append(m.Templates[:i], m.Templates[i+1:]...)
			return nil
		}
	}
	return nil
}
