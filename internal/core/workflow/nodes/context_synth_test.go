package nodes_test

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes"
	"github.com/stretchr/testify/assert"
)

func TestContextSynthesizer_Process(t *testing.T) {
	proc := &nodes.ContextSynthesizerProcessor{
		MaxRecentRounds: 2,
	}

	stream := make(chan workflow.StreamEvent, 100)

	// CASE 1: Append Round 1 (No pruning yet)
	// Input: existing summary (empty), new verdict
	input1 := map[string]interface{}{
		"history_summary": "",
		"new_verdict":     "## Round 1\nVerdict: Good code.",
		"round_summary":   "Round 1: Initial draft", // Heuristic summary
	}

	out1, err := proc.Process(context.Background(), input1, stream)
	assert.NoError(t, err)

	summary1 := out1["history_summary"].(string)
	assert.Contains(t, summary1, "## Chronological Verdicts")
	assert.Contains(t, summary1, "Verdict: Good code")
	assert.NotContains(t, summary1, "## Legacy Context")

	// CASE 2: Append Round 2
	input2 := map[string]interface{}{
		"history_summary": summary1,
		"new_verdict":     "## Round 2\nVerdict: Better code.",
		"round_summary":   "Round 2: Optimization",
	}
	out2, err := proc.Process(context.Background(), input2, stream)
	assert.NoError(t, err)
	summary2 := out2["history_summary"].(string)

	// CASE 3: Append Round 3 (Triggers Pruning of Round 1 because MaxRecentRounds=2)
	input3 := map[string]interface{}{
		"history_summary": summary2,
		"new_verdict":     "## Round 3\nVerdict: Perfect code.",
		"round_summary":   "Round 3: Final polish",
	}
	out3, err := proc.Process(context.Background(), input3, stream)
	assert.NoError(t, err)
	summary3 := out3["history_summary"].(string)

	// Verify Legacy
	assert.Contains(t, summary3, "## Legacy Context")
	assert.Contains(t, summary3, "- Round 1: Initial draft") // Pruned
	// Verify Active
	assert.Contains(t, summary3, "## Round 3")
	assert.Contains(t, summary3, "## Round 2")
	assert.NotContains(t, summary3, "## Round 1\nVerdict: Good code") // Should be gone
}
