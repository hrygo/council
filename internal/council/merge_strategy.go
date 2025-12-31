// Package council 包含 Council Debate/Optimize 工作流的应用层实现。
package council

import (
	"fmt"
	"strings"

	"github.com/hrygo/council/internal/core/workflow"
)

// Compile-time check: CouncilMergeStrategy implements MergeStrategy
var _ workflow.MergeStrategy = (*CouncilMergeStrategy)(nil)

// CouncilMergeStrategy 是 Council 工作流专用的聚合策略。
// 它将多个 Agent 的 agent_output 聚合为 aggregated_outputs，
// 并透传其他上下文字段。
type CouncilMergeStrategy struct{}

// Merge 聚合多个上游节点的输出。
func (s *CouncilMergeStrategy) Merge(inputs []map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	var agentOutputs []string

	for i, input := range inputs {
		// 聚合 agent_output
		if output, ok := input["agent_output"].(string); ok && output != "" {
			agentOutputs = append(agentOutputs, output)
		}

		// 透传上下文字段 (first-win 策略)
		for _, key := range CouncilContextKeys {
			if _, exists := merged[key]; !exists {
				if val, ok := input[key]; ok {
					merged[key] = val
				}
			}
		}

		// 仍然保留 branch_N 供调试，但不作为核心依赖
		merged[fmt.Sprintf("branch_%d", i)] = input
	}

	// 将多个 agent_output 聚合为 aggregated_outputs
	if len(agentOutputs) > 0 {
		merged["aggregated_outputs"] = strings.Join(agentOutputs, "\n\n---\n\n")
	}

	return merged
}
