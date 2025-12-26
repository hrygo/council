package workflow

import (
	"bufio"
	"embed"
	"strconv"
	"strings"
	"sync"
)

//go:embed model_pricing.csv
var pricingFile embed.FS

var (
	ModelPricing map[string]struct {
		InputCostPer1K  float64
		OutputCostPer1K float64
	}
	pricingOnce sync.Once
)

func initPricing() {
	ModelPricing = make(map[string]struct {
		InputCostPer1K  float64
		OutputCostPer1K float64
	})

	file, err := pricingFile.Open("model_pricing.csv")
	if err != nil {
		// Fallback to default pricing
		ModelPricing["default"] = struct {
			InputCostPer1K  float64
			OutputCostPer1K float64
		}{0.001, 0.002}
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}

		model := strings.TrimSpace(parts[0])
		inputCost, err1 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		outputCost, err2 := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)

		if err1 == nil && err2 == nil {
			ModelPricing[model] = struct {
				InputCostPer1K  float64
				OutputCostPer1K float64
			}{inputCost, outputCost}
		}
	}

	// Ensure default exists
	if _, ok := ModelPricing["default"]; !ok {
		ModelPricing["default"] = struct {
			InputCostPer1K  float64
			OutputCostPer1K float64
		}{0.001, 0.002}
	}
}

func GetModelPricing() map[string]struct {
	InputCostPer1K  float64
	OutputCostPer1K float64
} {
	pricingOnce.Do(initPricing)
	return ModelPricing
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
	pricing := GetModelPricing()
	model := "default"
	if m, ok := node.Properties["model"].(string); ok {
		model = m
	}
	// Note: If agent node, logic usually comes from the Agent definition.
	// For MVP, if property "model" is missing, we fall through to default.

	price, ok := pricing[model]
	if !ok {
		price = pricing["default"]
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
