package workflow

import (
	"testing"
)

func TestEstimateWorkflowCost(t *testing.T) {
	// Define a simple graph
	graph := &GraphDefinition{
		Nodes: map[string]*Node{
			"start": {
				ID:   "start",
				Type: NodeTypeStart,
			},
			"agent_1": {
				ID:   "agent_1",
				Type: NodeTypeAgent,
				Properties: map[string]interface{}{
					"agent_name": "Writer",
					"model":      "gpt-4",
				},
			},
			"agent_2": {
				ID:   "agent_2",
				Type: NodeTypeAgent,
				Properties: map[string]interface{}{
					"agent_name": "Reviewer",
					"model":      "gemini-1.5-flash", // Should be cheap
				},
			},
			"vote": {
				ID:   "vote",
				Type: NodeTypeVote,
			},
			"end": {
				ID:   "end",
				Type: NodeTypeEnd,
			},
		},
	}

	estimate := EstimateWorkflowCost(graph)

	if estimate.TotalTokens != 3000 { // 2 Agents * (1000 input + 500 output) = 3000
		t.Errorf("Expected 3000 tokens, got %d", estimate.TotalTokens)
	}

	// Calculate expected cost
	// gpt-4: (1*0.03) + (0.5*0.06) = 0.03 + 0.03 = 0.06
	// gemini-1.5-flash: (1*0.00035) + (0.5*0.00105) = 0.00035 + 0.000525 = 0.000875
	// Total: 0.060875
	expectedCost := 0.060875
	if estimate.TotalCostUSD != expectedCost {
		t.Errorf("Expected cost %f, got %f", expectedCost, estimate.TotalCostUSD)
	}

	// Check breakdown
	if val, ok := estimate.AgentBreakdown["Writer"]; !ok || val != 0.06 {
		t.Errorf("Writer cost incorrect, got %v", val)
	}
}

func TestEstimateWorkflowCost_LogicNodesFree(t *testing.T) {
	graph := &GraphDefinition{
		Nodes: map[string]*Node{
			"loop": {ID: "loop", Type: NodeTypeLoop},
			"vote": {ID: "vote", Type: NodeTypeVote},
		},
	}

	estimate := EstimateWorkflowCost(graph)
	if estimate.TotalCostUSD != 0 {
		t.Errorf("Logic nodes should be free, got %f", estimate.TotalCostUSD)
	}
}
