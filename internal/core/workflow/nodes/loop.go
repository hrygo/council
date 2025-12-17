package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type LoopProcessor struct {
	MaxRounds     int
	ExitCondition string
}

func (l *LoopProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// Loop state management usually happens in the Engine (via Session state).
	// The Node Processor decides if we SHOULD loop.

	// We need to know current iteration.
	// Convention: input["_loop_iteration"] or similar provided by Engine?
	// For now, let's assume input has "iteration" or we just return decision.

	currentRound, _ := input["iteration"].(int)
	if currentRound == 0 {
		currentRound = 1 // Default to 1st round
	}

	shouldExit := false
	if currentRound >= l.MaxRounds {
		shouldExit = true
	}

	// Check Exit Condition (e.g. Consensus from previous Vote node)
	if l.ExitCondition == "consensus" {
		if approved, ok := input["approved"].(bool); ok && approved {
			shouldExit = true
		}
	}

	output := map[string]interface{}{
		"should_exit":   shouldExit,
		"current_round": currentRound,
		"timestamp":     time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "completed", "should_exit": shouldExit},
	}

	return output, nil
}
