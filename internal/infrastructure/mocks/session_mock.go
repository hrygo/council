package mocks

import (
	"context"

	"github.com/hrygo/council/internal/core/workflow"
)

type SessionMockRepository struct {
	CapturedSessions []*workflow.Session
	Err              error
}

func NewSessionMockRepository() *SessionMockRepository {
	return &SessionMockRepository{}
}

func (m *SessionMockRepository) Create(ctx context.Context, session *workflow.Session, groupID string, workflowID string) error {
	if m.Err != nil {
		return m.Err
	}
	m.CapturedSessions = append(m.CapturedSessions, session)
	return nil
}

func (m *SessionMockRepository) Get(ctx context.Context, id string) (*workflow.SessionEntity, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	// Return a dummy entity with a default group
	return &workflow.SessionEntity{
		ID:      id,
		GroupID: "test-group",
		Status:  workflow.SessionRunning,
	}, nil
}

func (m *SessionMockRepository) UpdateStatus(ctx context.Context, id string, status workflow.SessionStatus) error {
	return m.Err
}
