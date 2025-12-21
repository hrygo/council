package workflow

// Pricing Model (Simple global map for MVP)
// ideally this comes from configuration
var ModelPricing = map[string]struct {
	InputCostPer1K  float64
	OutputCostPer1K float64
}{
	"gpt-4":             {0.03, 0.06},
	"gpt-4-turbo":       {0.01, 0.03},
	"claude-3-opus":     {0.015, 0.075},
	"claude-3-5-sonnet": {0.003, 0.015},
	"gemini-1.5-pro":    {0.0035, 0.0105},
	"gemini-1.5-flash":  {0.00035, 0.00105},
	"default":           {0.001, 0.002}, // Fallback
}

// EstimateWorkflowCost calculates the approximate cost of a workflow execution.
// This is a heuristic estimation based on node types and configured models.
func EstimateWorkflowCost(graph *GraphDefinition) *CostEstimate {
	estimate := &CostEstimate{
		TotalCostUSD:   0,
		TotalTokens:    0,
		NodeBreakdown:  make(map[string]NodeCostEstimate),
		AgentBreakdown: make(map[string]float64),
	}

	for _, node := range graph.Nodes {
		nodeCost := estimateNodeCost(node)
		estimate.TotalCostUSD += nodeCost.CostUSD
		estimate.TotalTokens += nodeCost.Tokens
		estimate.NodeBreakdown[node.ID] = nodeCost

		if agentName, ok := node.Properties["agent_name"].(string); ok {
			estimate.AgentBreakdown[agentName] += nodeCost.CostUSD
		}
	}

	return estimate
}

func estimateNodeCost(node *Node) NodeCostEstimate {
	// Base assumption: 1000 input tokens, 500 output tokens per "LLM interaction"
	// This is very rough and implies the user knows the prompt size.
	// We could try to count tokens in the prompt property if it exists.

	avgInput := 1000
	avgOutput := 500

	// Logic nodes are free
	if node.Type == NodeTypeStart || node.Type == NodeTypeEnd || node.Type == NodeTypeParallel || node.Type == NodeTypeVote || node.Type == NodeTypeLoop {
		return NodeCostEstimate{CostUSD: 0, Tokens: 0}
	}

	// Try to get model from properties to find price
	model := "default"
	if m, ok := node.Properties["model"].(string); ok {
		model = m
	}
	// Note: If agent node, logic usually comes from the Agent definition. 
	// For MVP, if property "model" is missing, we fall through to default.

	price, ok := ModelPricing[model]
	if !ok {
		price = ModelPricing["default"]
	}

	cost := (float64(avgInput)/1000)*price.InputCostPer1K + (float64(avgOutput)/1000)*price.OutputCostPer1K

	return NodeCostEstimate{
		CostUSD: cost,
		Tokens:  avgInput + avgOutput,
	}
}

type CostEstimate struct {
	TotalCostUSD   float64                     `json:"total_cost_usd"`
	TotalTokens    int                         `json:"total_tokens"`
	NodeBreakdown  map[string]NodeCostEstimate `json:"node_breakdown"`
	AgentBreakdown map[string]float64          `json:"agent_breakdown"`
}

type NodeCostEstimate struct {
	CostUSD float64 `json:"cost_usd"`
	Tokens  int     `json:"tokens"`
}
