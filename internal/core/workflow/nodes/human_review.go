package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type HumanReviewProcessor struct {
	TimeoutMinutes int
	AllowSkip      bool
}

func (h *HumanReviewProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// Human Review logic involves suspending the workflow.
	// The processor typically sets a flag or returns a special error/status that the Engine interprets as "Suspend".

	// However, the Engine's `ExecuteStep` waits for return.
	// We need a mechanism to tell Engine to suspend.
	// For now, we'll return a specific output key that Engine checks, OR we can blocking wait here (bad for resources).
	// Better: Engine handles NodeTypeHumanReview specially, or we return an error "ErrSuspended".

	stream <- workflow.StreamEvent{
		Type:      "human_interaction_required",
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"message": "Waiting for human review",
			"timeout": h.TimeoutMinutes,
		},
	}

	// In a real system, this function might return immediately with a "Suspended" status.
	// The Session wrapper would then pause execution loop.

	// Mock: Auto-approve for dev speed if 'AllowSkip' is true, else pretend we wait (but we can't block indefinitely here in sync model without blocking worker).

	output := map[string]interface{}{
		"status":    "pending_human",
		"timestamp": time.Now(),
	}

	// Note: Engine must handle this output to assume State=Suspended.

	return output, nil
}
