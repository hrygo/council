package nodes

import (
	"fmt"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

// Dependencies for creating nodes
type NodeDependencies struct {
	Registry      *llm.Registry
	AgentRepo     agent.Repository
	MemoryManager memory.MemoryManager
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
				model = deps.Registry.GetDefaultModel()
			}

			// EndProcessor currently uses LLMProvider directly.
			// We need to resolve it. Since EndProcessor logic might be simple,
			// we can just use Default provider or resolve if it supports it.
			// Let's assume EndProcessor needs a specific provider.
			// For now, let's pass the Registry and let it resolve or resolve System Default here.
			provider, err := deps.Registry.GetLLMProvider("default")
			if err != nil {
				return nil, fmt.Errorf("failed to get default LLM provider: %w", err)
			}

			return &EndProcessor{
				LLM:    provider,
				Model:  model,
				Prompt: prompt,
			}, nil

		case workflow.NodeTypeAgent:
			agentID, _ := node.Properties["agent_id"].(string)
			if agentID == "" {
				return nil, fmt.Errorf("agent_id property missing for node %s", node.ID)
			}
			return &AgentProcessor{
				NodeID:    node.ID,
				AgentID:   agentID,
				AgentRepo: deps.AgentRepo,
				Registry:  deps.Registry,
			}, nil

		case workflow.NodeTypeVote:
			threshold, _ := node.Properties["threshold"].(float64)
			voteType, _ := node.Properties["vote_type"].(string)
			return &VoteProcessor{
				Threshold: threshold,
				VoteType:  voteType,
			}, nil

		case workflow.NodeTypeLoop:
			maxRounds, _ := node.Properties["max_rounds"].(float64)      // JSON numbers often float64
			exitOnScore, _ := node.Properties["exit_on_score"].(float64) // SPEC-609 Defect-3 fix
			return &LoopProcessor{
				MaxRounds:   int(maxRounds),
				ExitOnScore: int(exitOnScore),
			}, nil

		case workflow.NodeTypeFactCheck:
			threshold, _ := node.Properties["verify_threshold"].(float64)
			provider, err := deps.Registry.GetLLMProvider("default")
			if err != nil {
				return nil, fmt.Errorf("failed to get default LLM for fact check: %w", err)
			}
			return &FactCheckProcessor{
				LLM:             provider,
				VerifyThreshold: threshold,
			}, nil

		case workflow.NodeTypeHumanReview:
			timeout, _ := node.Properties["timeout_minutes"].(float64)
			allowSkip, _ := node.Properties["allow_skip"].(bool)
			return &HumanReviewProcessor{
				TimeoutMinutes: int(timeout),
				AllowSkip:      allowSkip,
			}, nil

		case workflow.NodeTypeMemoryRetrieval:
			// SPEC-607: Memory Retrieval Node
			return NewMemoryRetrievalProcessor(deps.MemoryManager), nil

		default:
			return nil, fmt.Errorf("unsupported node type: %s", node.Type)
		}
	}
}
