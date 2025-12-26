package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
)

// MemoryRetrievalProcessor retrieves historical context from Memory system.
type MemoryRetrievalProcessor struct {
	MemoryManager memory.MemoryManager
}

// NewMemoryRetrievalProcessor creates a new MemoryRetrievalProcessor.
func NewMemoryRetrievalProcessor(mm memory.MemoryManager) *MemoryRetrievalProcessor {
	return &MemoryRetrievalProcessor{MemoryManager: mm}
}

// Process retrieves relevant historical context and injects it into the output.
func (p *MemoryRetrievalProcessor) Process(
	ctx context.Context,
	input map[string]interface{},
	stream chan<- workflow.StreamEvent,
) (map[string]interface{}, error) {
	// Send start event
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// 1. Extract topic/query from input
	topic, _ := input["topic"].(string)
	if topic == "" {
		// Fallback: use content or document as query
		if content, ok := input["content"].(string); ok {
			topic = content
		} else if doc, ok := input["document"].(string); ok {
			topic = doc
		}
	}

	// 2. Extract groupID from input (required for Memory service)
	groupID, _ := input["group_id"].(string)

	// 3. Retrieve from Memory system if available
	var historyContext string
	if p.MemoryManager != nil && topic != "" {
		sessionID, _ := input["session_id"].(string)
		items, err := p.MemoryManager.Retrieve(ctx, topic, groupID, sessionID)
		if err != nil {
			// Log but don't fail - memory retrieval is optional
			stream <- workflow.StreamEvent{
				Type:      "memory_retrieval_error",
				Timestamp: time.Now(),
				Data:      map[string]interface{}{"error": err.Error()},
			}
		} else {
			historyContext = formatHistorySummary(items)
		}
	}

	// 4. Build output with injected history context
	output := make(map[string]interface{})
	for k, v := range input {
		output[k] = v
	}
	output["history_context"] = historyContext

	// 5. Emit completion event
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"status":         "completed",
			"context_length": len(historyContext),
		},
	}

	return output, nil
}

// formatHistorySummary formats memory items into a markdown summary.
func formatHistorySummary(items []memory.ContextItem) string {
	if len(items) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("## 历史上下文\n\n")

	for i, item := range items {
		sb.WriteString(fmt.Sprintf("### 记录 %d (相关度: %.2f)\n", i+1, item.Score))
		sb.WriteString(item.Content)
		sb.WriteString("\n\n")
	}

	return sb.String()
}
