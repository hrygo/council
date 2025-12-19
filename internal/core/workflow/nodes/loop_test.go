package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
)

func TestLoopProcessor_Process(t *testing.T) {
	p := &LoopProcessor{
		MaxRounds:     3,
		ExitCondition: "consensus",
	}

	stream := make(chan workflow.StreamEvent, 10)

	// Test case 1: Iteration 1, not consensus
	input1 := map[string]interface{}{"iteration": 1, "approved": false}
	output1, _ := p.Process(context.Background(), input1, stream)
	if output1["should_exit"].(bool) {
		t.Error("expected should_exit to be false for round 1")
	}

	// Test case 2: Iteration 3 (Max)
	input2 := map[string]interface{}{"iteration": 3}
	output2, _ := p.Process(context.Background(), input2, stream)
	if !output2["should_exit"].(bool) {
		t.Error("expected should_exit to be true for round 3")
	}

	// Test case 3: Consensus reached early
	input3 := map[string]interface{}{"iteration": 2, "approved": true}
	output3, _ := p.Process(context.Background(), input3, stream)
	if !output3["should_exit"].(bool) {
		t.Error("expected should_exit to be true for consensus")
	}
}
