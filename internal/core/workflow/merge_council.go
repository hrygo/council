package workflow

import (
	"fmt"
	"strings"
)

// CouncilMergeStrategy is an application-level merge strategy designed for
// Council Debate workflows. It specifically aggregates agent_output fields
// from parallel branches into a single "aggregated_outputs" field.
type CouncilMergeStrategy struct{}

func (s *CouncilMergeStrategy) Merge(inputs []map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	var agentOutputs []string

	for i, inp := range inputs {
		// Collect all agent_output values
		if out, ok := inp["agent_output"].(string); ok && out != "" {
			agentOutputs = append(agentOutputs, out)
		}

		// Pass through first occurrence of context fields (except agent_output)
		for k, v := range inp {
			if k == "agent_output" {
				continue // Already handled specially
			}
			if _, exists := merged[k]; !exists {
				merged[k] = v
			}
		}

		// Preserve branch data for debugging
		merged[fmt.Sprintf("branch_%d", i)] = inp
	}

	// Aggregate all agent outputs into a single field
	if len(agentOutputs) > 0 {
		merged["aggregated_outputs"] = strings.Join(agentOutputs, "\n\n---\n\n")
	}

	return merged
}
