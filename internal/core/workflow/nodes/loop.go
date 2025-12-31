package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type LoopProcessor struct {
	MaxRounds       int
	ExitOnScore     int // Score threshold for automatic exit (e.g., 90)
	Session         *workflow.Session
	PassthroughKeys []string // Configuration: Keys to pass for loop continuation
}

// GetNextNodes implements workflow.ConditionalRouter for loop logic.
func (l *LoopProcessor) GetNextNodes(ctx context.Context, output map[string]interface{}, defaultNextIDs []string) ([]string, error) {
	if len(defaultNextIDs) < 2 {
		// If less than 2, standard flow (maybe just loop back or just exit?)
		// But typically Loop must have 2 branches.
		return defaultNextIDs, nil
	}

	shouldExit, _ := output["should_exit"].(bool)
	if shouldExit {
		// Exit path: second next_id
		return []string{defaultNextIDs[1]}, nil
	}
	// Continue path: first next_id
	return []string{defaultNextIDs[0]}, nil
}

func (l *LoopProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Logic start
	_ = stream // usage in events if needed

	// ... (logic omitted for brevity, keeping existing code logic)

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
	currentScore := 0.0
	if val, ok := input["score"].(float64); ok {
		currentScore = val
	} else if val, ok := input["score"].(int); ok {
		currentScore = float64(val)
	} else {
		// Try to find score in history if not in direct input
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

			// Detect Regression
			if delta < -10 {
				exitReason = "regression_detected"
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

	// Passthrough context fields (Config Driven)
	workflow.ApplyPassthrough(input, output, workflow.PassthroughConfig{
		Keys: l.PassthroughKeys,
	})

	_ = 0 // loop logic complete

	return output, nil
}
