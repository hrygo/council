package mocks

import (
	"context"

	"github.com/hrygo/council/internal/core/memory"
)

type MemoryMockManager struct {
	CapturedQuarantine []string
	CapturedWM         []string
	RetrieveResult     []memory.ContextItem
	Err                error
}

func (m *MemoryMockManager) LogQuarantine(ctx context.Context, sessionID string, nodeID string, content string, metadata map[string]interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	m.CapturedQuarantine = append(m.CapturedQuarantine, content)
	return nil
}

func (m *MemoryMockManager) UpdateWorkingMemory(ctx context.Context, groupID string, content string, metadata map[string]interface{}) error {
	if m.Err != nil {
		return m.Err
	}
	m.CapturedWM = append(m.CapturedWM, content)
	return nil
}

func (m *MemoryMockManager) Promote(ctx context.Context, groupID string, digest string) error {
	return m.Err
}

func (m *MemoryMockManager) Retrieve(ctx context.Context, query string, groupID string, sessionID string) ([]memory.ContextItem, error) {
	return m.RetrieveResult, m.Err
}
