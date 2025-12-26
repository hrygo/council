package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
)

func TestVoteProcessor_Process(t *testing.T) {
	p := &VoteProcessor{
		Threshold: 0.6,
	}

	stream := make(chan workflow.StreamEvent, 10)

	// Test case 1: 2 YES, 1 NO (ratio 0.66 >= 0.6)
	input1 := map[string]interface{}{
		"a1": "YES",
		"a2": "YES",
		"a3": "NO",
	}
	output1, _ := p.Process(context.Background(), input1, stream)
	if !output1["approved"].(bool) {
		t.Error("expected approved to be true for 2/3 ratio")
	}

	// Test case 2: 1 YES, 2 NO (ratio 0.33 < 0.6)
	input2 := map[string]interface{}{
		"a1": "YES",
		"a2": "NO",
		"a3": "NO",
	}
	output2, _ := p.Process(context.Background(), input2, stream)
	if output2["approved"].(bool) {
		t.Error("expected approved to be false for 1/3 ratio")
	}

	// Test case 3: empty input (fallthrough mock)
	input3 := map[string]interface{}{}
	output3, _ := p.Process(context.Background(), input3, stream)
	if !output3["approved"].(bool) {
		t.Error("expected approved to be true for empty input mock")
	}
}

func TestVoteProcessor_ExactThreshold(t *testing.T) {
	p := &VoteProcessor{Threshold: 0.5}
	stream := make(chan workflow.StreamEvent, 10)

	// Exactly at threshold: 1 YES, 1 NO (ratio 0.5 >= 0.5)
	input := map[string]interface{}{"a1": "YES", "a2": "NO"}
	output, _ := p.Process(context.Background(), input, stream)
	if !output["approved"].(bool) {
		t.Error("expected approved when ratio equals threshold")
	}
}

func TestVoteProcessor_AllApproved(t *testing.T) {
	p := &VoteProcessor{Threshold: 0.9}
	stream := make(chan workflow.StreamEvent, 10)

	// All approved
	input := map[string]interface{}{"a1": "APPROVED", "a2": "APPROVED"}
	output, _ := p.Process(context.Background(), input, stream)
	if !output["approved"].(bool) {
		t.Error("expected approved for all APPROVED votes")
	}
	if output["yes_votes"].(int) != 2 {
		t.Errorf("expected 2 yes votes, got %v", output["yes_votes"])
	}
}
