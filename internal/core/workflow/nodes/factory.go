package nodes

import (
	"fmt"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

// Dependencies for creating nodes
type NodeDependencies struct {
	LLM       llm.LLMProvider
	AgentRepo agent.Repository
}

// NewNodeFactory returns a factory function compatible with workflow.Engine
func NewNodeFactory(deps NodeDependencies) func(node *workflow.Node) (workflow.NodeProcessor, error) {
	return func(node *workflow.Node) (workflow.NodeProcessor, error) {
		switch node.Type {
		case workflow.NodeTypeStart:
			return &StartProcessor{}, nil

		case workflow.NodeTypeEnd:
			// Extract config from node.Properties usually
			prompt, _ := node.Properties["summary_prompt"].(string)
			model, _ := node.Properties["model"].(string)
			if model == "" {
				model = "gpt-4"
			}
			return &EndProcessor{
				LLM:    deps.LLM,
				Model:  model,
				Prompt: prompt,
			}, nil

		case workflow.NodeTypeAgent:
			agentID, _ := node.Properties["agent_id"].(string)
			if agentID == "" {
				return nil, fmt.Errorf("agent_id property missing for node %s", node.ID)
			}
			return &AgentProcessor{
				AgentID:   agentID,
				AgentRepo: deps.AgentRepo,
				LLM:       deps.LLM,
			}, nil

		case workflow.NodeTypeVote:
			threshold, _ := node.Properties["threshold"].(float64)
			voteType, _ := node.Properties["vote_type"].(string)
			return &VoteProcessor{
				Threshold: threshold,
				VoteType:  voteType,
			}, nil

		case workflow.NodeTypeLoop:
			maxRounds, _ := node.Properties["max_rounds"].(float64) // JSON numbers often float64
			exitCond, _ := node.Properties["exit_condition"].(string)
			return &LoopProcessor{
				MaxRounds:     int(maxRounds),
				ExitCondition: exitCond,
			}, nil

		case workflow.NodeTypeFactCheck:
			threshold, _ := node.Properties["verify_threshold"].(float64)
			return &FactCheckProcessor{
				LLM:             deps.LLM,
				VerifyThreshold: threshold,
			}, nil

		case workflow.NodeTypeHumanReview:
			timeout, _ := node.Properties["timeout_minutes"].(float64)
			allowSkip, _ := node.Properties["allow_skip"].(bool)
			return &HumanReviewProcessor{
				TimeoutMinutes: int(timeout),
				AllowSkip:      allowSkip,
			}, nil

		default:
			return nil, fmt.Errorf("unsupported node type: %s", node.Type)
		}
	}
}
