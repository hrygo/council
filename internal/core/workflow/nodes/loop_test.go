package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
)

func TestLoopProcessor_Process(t *testing.T) {
	p := &LoopProcessor{
		MaxRounds:   3,
		ExitOnScore: 90, // Exit when score >= 90
	}

	stream := make(chan workflow.StreamEvent, 10)

	// Test case 1: Iteration 1, score below threshold
	input1 := map[string]interface{}{"iteration": 1, "score": 75.0}
	output1, _ := p.Process(context.Background(), input1, stream)
	if output1["should_exit"].(bool) {
		t.Error("expected should_exit to be false for round 1 with score 75")
	}

	// Test case 2: Iteration 3 (Max rounds reached)
	input2 := map[string]interface{}{"iteration": 3, "score": 80.0}
	output2, _ := p.Process(context.Background(), input2, stream)
	if !output2["should_exit"].(bool) {
		t.Error("expected should_exit to be true when max rounds reached")
	}
	if output2["exit_reason"].(string) != "max_rounds_reached" {
		t.Error("expected exit_reason to be max_rounds_reached")
	}

	// Test case 3: Score threshold reached early
	input3 := map[string]interface{}{"iteration": 2, "score": 92.0}
	output3, _ := p.Process(context.Background(), input3, stream)
	if !output3["should_exit"].(bool) {
		t.Error("expected should_exit to be true when score threshold reached")
	}
	if output3["exit_reason"].(string) != "score_threshold_reached" {
		t.Error("expected exit_reason to be score_threshold_reached")
	}
}
