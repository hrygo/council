package middleware

import (
	"context"
	"fmt"

	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
)

// MemoryMiddleware handles automatic memory persistence (Quarantine & Working Memory)
type MemoryMiddleware struct {
	Manager memory.MemoryManager
}

func NewMemoryMiddleware(manager memory.MemoryManager) *MemoryMiddleware {
	return &MemoryMiddleware{Manager: manager}
}

func (mm *MemoryMiddleware) Name() string {
	return "MemoryProtocol"
}

func (mm *MemoryMiddleware) BeforeNodeExecution(ctx context.Context, session *workflow.Session, node *workflow.Node) error {
	return nil
}

func (mm *MemoryMiddleware) AfterNodeExecution(ctx context.Context, session *workflow.Session, node *workflow.Node, output map[string]interface{}) (map[string]interface{}, error) {
	// Extract Content
	content, ok := output["content"].(string)
	if !ok {
		// If no content, maybe we skip logging or log partial?
		// For MVP, if no content, we assume it's a structural node and skip,
		// OR we log full output as JSON in content?
		// Let's Log "No Content" but keeping interaction record.
		// Actually typical Agent node output has "content".
		return output, nil
	}

	metadata, _ := output["metadata"].(map[string]interface{})

	// Tier 1: Quarantine Log (Always)
	if err := mm.Manager.LogQuarantine(ctx, session.ID, node.ID, content, metadata); err != nil {
		// Log error but don't fail the workflow flow?
		// "Observation" shouldn't block execution unless critical.
		// For now, allow failure.
		return output, fmt.Errorf("memory log failed: %w", err)
	}

	// Tier 2: Working Memory (Attempt)
	// We need GroupID. Session has GroupID?
	// Session struct currently has: ID, Graph, Inputs...
	// Session doesn't strictly have GroupID unless passed in Inputs or Graph properties.
	// We assume Graph.Properties["group_id"] or Inputs["group_id"]?
	// TDD says: Session belongs to Group.
	// For MVP, we extract group ID from where available.

	// Assuming GroupID is in session Input for now.
	var groupID string
	if val, ok := session.Inputs["group_id"].(string); ok {
		groupID = val
	}

	if groupID != "" {
		if err := mm.Manager.UpdateWorkingMemory(ctx, groupID, content, metadata); err != nil {
			// Log and ignore error for now to prevent workflow interruption
			fmt.Printf("warning: working memory update failed: %v\n", err)
		}
	}

	return output, nil
}
