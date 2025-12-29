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

type MockSessionFileRepository struct {
	Files    map[string]*workflow.FileEntity
	Versions map[string][]*workflow.FileEntity
}

func NewMockSessionFileRepository() *MockSessionFileRepository {
	return &MockSessionFileRepository{
		Files:    make(map[string]*workflow.FileEntity),
		Versions: make(map[string][]*workflow.FileEntity),
	}
}

func (m *MockSessionFileRepository) AddVersion(ctx context.Context, sessionID, path, content, author, reason string) (int, error) {
	key := sessionID + ":" + path // Simple keying for mock
	currentVer := 0
	if f, ok := m.Files[key]; ok {
		currentVer = f.Version
	}
	newVer := currentVer + 1

	f := &workflow.FileEntity{
		ID:        "uuid-" + path,
		SessionID: sessionID,
		Path:      path,
		Version:   newVer,
		Content:   content,
		Author:    author,
		Reason:    reason,
	}
	m.Files[key] = f // Latest
	m.Versions[key] = append(m.Versions[key], f)
	return newVer, nil
}

func (m *MockSessionFileRepository) GetLatest(ctx context.Context, sessionID, path string) (*workflow.FileEntity, error) {
	key := sessionID + ":" + path
	if f, ok := m.Files[key]; ok {
		return f, nil
	}
	return nil, nil // Or error not found
}

func (m *MockSessionFileRepository) ListFiles(ctx context.Context, sessionID string) (map[string]*workflow.FileEntity, error) {
	// Simplified: return all in mock
	res := make(map[string]*workflow.FileEntity)
	// Filter by sessionID logic skipped for simplicity if we assume unique sessions in tests or handle carefully
	// But let's do it rightish
	for k, v := range m.Files {
		// key is sessionID:path
		// check prefix
		// implementation detail hidden
		res[k] = v
	}
	return res, nil
}

func (m *MockSessionFileRepository) ListVersions(ctx context.Context, sessionID, path string) ([]*workflow.FileEntity, error) {
	key := sessionID + ":" + path
	return m.Versions[key], nil
}
