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
	// Better: Engine handles NodeTypeHumanReview specially, or we return an error
	// Emit event to notify UI
	stream <- workflow.StreamEvent{
		Type:      "human_interaction_required",
		Timestamp: time.Now(),
		// Assuming 'p' is the receiver and 'Node' is a field, and 'timeout' is defined.
		// Since these are not in the original code, I'll use 'h' and 'h.TimeoutMinutes'
		// to maintain consistency with the provided context.
		// If 'p.Node.ID' and 'timeout' are intended, the struct and method signature
		// would need to be updated, which is beyond the scope of this specific instruction.
		Data: map[string]interface{}{
			"reason":  "Human review required",
			"timeout": h.TimeoutMinutes,
		},
	}

	// Return ErrSuspended to pause execution at this node
	return nil, workflow.ErrSuspended
}
