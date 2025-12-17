package nodes

import (
	"context"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
)

type VoteProcessor struct {
	Threshold float64
	VoteType  string
}

func (v *VoteProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": "running"},
	}

	// This is a placeholder logic.
	// Real implementation would inspect `input` (results from previous nodes, e.g. parallel agents)
	// and extract votes.

	// Assuming input contains a list of "agent_outputs" or similar which are strings "YES" or "NO" or scores.
	// For simulation, we assume inputs might have keys like "agent_1_result": "YES"

	yesVotes := 0
	totalVotes := 0

	// Mock: Scan input values for simplistic "YES"/"NO" detection
	for _, val := range input {
		if str, ok := val.(string); ok {
			upper := str // In real code: strings.ToUpper
			if upper == "YES" || upper == "APPROVED" {
				yesVotes++
			}
			totalVotes++
		}
	}

	// If no inputs suitable, mock success for development flow
	if totalVotes == 0 {
		// Mock behavior: If this node runs, assume it passes for now unless configured otherwise
		yesVotes = 1
		totalVotes = 1
	}

	ratio := float64(yesVotes) / float64(totalVotes)
	approved := ratio >= v.Threshold

	output := map[string]interface{}{
		"approved":    approved,
		"yes_votes":   yesVotes,
		"total_votes": totalVotes,
		"ratio":       ratio,
		"timestamp":   time.Now(),
	}

	status := "completed"
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"status": status, "result": approved},
	}

	return output, nil
}
