package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
)

// AgentMockRepository implements Repository for testing.
type AgentMockRepository struct {
	Agents map[uuid.UUID]*agent.Agent
	Err    error
}

func NewAgentMockRepository() *AgentMockRepository {
	return &AgentMockRepository{
		Agents: make(map[uuid.UUID]*agent.Agent),
	}
}

func (m *AgentMockRepository) Create(ctx context.Context, a *agent.Agent) error {
	if m.Err != nil {
		return m.Err
	}
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	m.Agents[a.ID] = a
	return nil
}

func (m *AgentMockRepository) GetByID(ctx context.Context, id uuid.UUID) (*agent.Agent, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	a, ok := m.Agents[id]
	if !ok {
		return nil, context.DeadlineExceeded
	}
	return a, nil
}

func (m *AgentMockRepository) List(ctx context.Context) ([]*agent.Agent, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	var list []*agent.Agent
	for _, a := range m.Agents {
		list = append(list, a)
	}
	return list, nil
}

func (m *AgentMockRepository) Update(ctx context.Context, a *agent.Agent) error {
	if m.Err != nil {
		return m.Err
	}
	if _, ok := m.Agents[a.ID]; !ok {
		return context.DeadlineExceeded
	}
	m.Agents[a.ID] = a
	return nil
}

func (m *AgentMockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.Err != nil {
		return m.Err
	}
	delete(m.Agents, id)
	return nil
}
