package workflow

import "fmt"

// MergeStrategy defines how to merge outputs from multiple upstream nodes
// when a node has in-degree > 1 (e.g., after parallel branches converge).
// This is a framework-level interface - implementations can be application-specific.
type MergeStrategy interface {
	// Merge receives outputs from all upstream branches and returns a merged input
	// for the downstream node.
	Merge(inputs []map[string]interface{}) map[string]interface{}
}

// DefaultMergeStrategy is the basic merge strategy that:
// 1. Preserves each branch's full output keyed by "branch_N"
// 2. Passes through the first occurrence of any field (first-come-first-served)
type DefaultMergeStrategy struct{}

func (s *DefaultMergeStrategy) Merge(inputs []map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	for i, inp := range inputs {
		// Preserve each branch's complete output for debugging/traceability
		merged[fmt.Sprintf("branch_%d", i)] = inp

		// Pass through first occurrence of each field
		for k, v := range inp {
			if _, exists := merged[k]; !exists {
				merged[k] = v
			}
		}
	}

	return merged
}
