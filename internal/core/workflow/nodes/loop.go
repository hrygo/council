package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type LoopProcessor struct {
	MaxRounds   int
	ExitOnScore int // Score threshold for automatic exit (e.g., 90)
	Session     *workflow.Session
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
	// Check Exit on Score (SPEC-609 Defect-3 fix)
	currentScore := 0.0
	if val, ok := input["score"].(float64); ok {
		currentScore = val
	} else if val, ok := input["score"].(int); ok {
		currentScore = float64(val)
	} else {
		// Try to find score in history if not in direct input
		// Assuming input contains all node outputs, look for "agent_adjudicator" or "structured_score"
		// This part depends on how Engine passes inputs. Assuming flat map of all outputs.
		if scoreMap, ok := input["structured_score"].(map[string]interface{}); ok {
			if s, ok := scoreMap["weighted_score"].(float64); ok {
				currentScore = s
			}
		}
	}

	// Persist History
	if l.Session != nil {
		history, _ := l.Session.GetContext("score_history").([]float64)
		history = append(history, currentScore)
		l.Session.SetContext("score_history", history)

		// Calculate Delta
		if len(history) > 1 {
			prevScore := history[len(history)-2]
			delta := currentScore - prevScore

			// Detect Regression (Rollback) - Threshold -10
			if delta < -10 {
				// TODO: Trigger VFS Rollback?
				// For now just log usage
				exitReason = "regression_detected"
				// shouldExit = true? Or allow loop to try fix?
				// Usually loop restarts or tries again.
			}
		}
	}

	if l.ExitOnScore > 0 && currentScore >= float64(l.ExitOnScore) {
		shouldExit = true
		exitReason = "score_threshold_reached"
	}

	output := map[string]interface{}{
		"should_exit":   shouldExit,
		"exit_reason":   exitReason,
		"current_round": currentRound,
		"timestamp":     time.Now(),
	}

	// Passthrough context fields for loop continuation (Fix-C6)
	loopPassthroughKeys := []string{
		"document_content",
		"proposal",
		"optimization_objective",
		"combined_context",
		"session_id",
	}
	for _, key := range loopPassthroughKeys {
		if val, ok := input[key]; ok {
			output[key] = val
		}
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "completed", "should_exit": shouldExit, "exit_reason": exitReason},
	}

	return output, nil
}
