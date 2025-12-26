package memory

import (
	"context"
)

// MemoryManager defines the Three-Tier Memory Protocol interface
type MemoryManager interface {
	// Tier 1: Quarantine
	LogQuarantine(ctx context.Context, sessionID string, nodeID string, content string, metadata map[string]interface{}) error

	// Tier 2: Working Memory
	UpdateWorkingMemory(ctx context.Context, groupID string, content string, metadata map[string]interface{}) error

	// Tier 3: Promotion (Future)
	Promote(ctx context.Context, groupID string, digest string) error

	// Hybrid Retrieval
	Retrieve(ctx context.Context, query string, groupID string, sessionID string) ([]ContextItem, error)
}

type ContextItem struct {
	Content string
	Source  string // "hot", "cold"
	Score   float64
}
