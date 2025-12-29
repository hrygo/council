package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/stretchr/testify/assert"
)

func TestLoopProcessor_HistoryOps(t *testing.T) {
	session := workflow.NewSession(nil, nil)
	session.SetContext("score_history", []float64{80.0})

	proc := &LoopProcessor{
		MaxRounds:   5,
		ExitOnScore: 90,
		Session:     session,
	}

	stream := make(chan workflow.StreamEvent, 100)

	// Round 1: Score 85 (Delta +5)
	input := map[string]interface{}{"iteration": 1, "score": 85.0}
	output, err := proc.Process(context.Background(), input, stream)
	assert.NoError(t, err)
	assert.False(t, output["should_exit"].(bool))

	history := session.GetContext("score_history").([]float64)
	assert.Equal(t, 2, len(history))
	assert.Equal(t, 85.0, history[1])

	// Round 2: Score 70 (Delta -15) -> Regression
	input2 := map[string]interface{}{"iteration": 2, "score": 70.0}
	output2, err := proc.Process(context.Background(), input2, stream)
	assert.NoError(t, err)

	assert.Equal(t, "regression_detected", output2["exit_reason"])
	// In current impl, we didn't force exit on regression, just set reason.
	// Logic says: if l.ExitOnScore...
	// We might want to handle regression explicitly in output.
}
