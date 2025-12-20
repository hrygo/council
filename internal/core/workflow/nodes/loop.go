package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type LoopProcessor struct {
	MaxRounds   int
	ExitOnScore int // Score threshold for automatic exit (e.g., 90)
}

func (l *LoopProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// Loop state management usually happens in the Engine (via Session state).
	// The Node Processor decides if we SHOULD loop.

	currentRound, _ := input["iteration"].(int)
	if currentRound == 0 {
		currentRound = 1 // Default to 1st round
	}

	shouldExit := false
	exitReason := ""

	// Check max rounds
	if currentRound >= l.MaxRounds {
		shouldExit = true
		exitReason = "max_rounds_reached"
	}

	// Check Exit on Score (SPEC-609 Defect-3 fix)
	if l.ExitOnScore > 0 {
		if score, ok := input["score"].(float64); ok && int(score) >= l.ExitOnScore {
			shouldExit = true
			exitReason = "score_threshold_reached"
		}
		// Also check int type
		if score, ok := input["score"].(int); ok && score >= l.ExitOnScore {
			shouldExit = true
			exitReason = "score_threshold_reached"
		}
	}

	output := map[string]interface{}{
		"should_exit":   shouldExit,
		"exit_reason":   exitReason,
		"current_round": currentRound,
		"timestamp":     time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "completed", "should_exit": shouldExit, "exit_reason": exitReason},
	}

	return output, nil
}
