package nodes

import (
	"fmt"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes/tools"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

// Dependencies for creating nodes
// GenericNodeFactory implements workflow.NodeFactory for standard core nodes.
type GenericNodeFactory struct {
	Registry      *llm.Registry
	AgentRepo     agent.Repository
	MemoryManager memory.MemoryManager
}

// NewGenericNodeFactory creates a new factory with dependencies.
func NewGenericNodeFactory(registry *llm.Registry, agentRepo agent.Repository, memManager memory.MemoryManager) *GenericNodeFactory {
	return &GenericNodeFactory{
		Registry:      registry,
		AgentRepo:     agentRepo,
		MemoryManager: memManager,
	}
}

// CreateNode creates standard nodes.
func (f *GenericNodeFactory) CreateNode(node *workflow.Node, deps workflow.FactoryDeps) (workflow.NodeProcessor, error) {
	switch node.Type {
	case workflow.NodeTypeStart:
		return &StartProcessor{}, nil

	case workflow.NodeTypeEnd:
		prompt, _ := node.Properties["summary_prompt"].(string)
		model, _ := node.Properties["model"].(string)
		if model == "" && f.Registry != nil {
			model = f.Registry.GetDefaultModel()
		}

		// Resolve LLM
		var provider llm.LLMProvider
		var err error
		if f.Registry != nil {
			provider, err = f.Registry.GetLLMProvider("default")
			if err != nil {
				return nil, fmt.Errorf("failed to get default LLM provider: %w", err)
			}
		}

		return &EndProcessor{
			NodeID:    node.ID,
			LLM:       provider,
			Model:     model,
			Prompt:    prompt,
			OutputKey: "summary",
		}, nil

	case workflow.NodeTypeAgent:
		agentID, _ := node.Properties["agent_uuid"].(string)
		if agentID == "" {
			return nil, fmt.Errorf("agent_uuid property missing for node %s", node.ID)
		}

		var processorTools []tools.Tool
		if toolNames, ok := node.Properties["tools"].([]interface{}); ok {
			for _, t := range toolNames {
				if name, ok := t.(string); ok {
					switch name {
					case "write_file":
						processorTools = append(processorTools, &tools.WriteFileTool{})
					case "read_file":
						processorTools = append(processorTools, &tools.ReadFileTool{})
					}
				}
			}
		}

		return &AgentProcessor{
			NodeID:    node.ID,
			AgentID:   agentID,
			AgentRepo: f.AgentRepo,
			Registry:  f.Registry,
			Tools:     processorTools,
			Session:   deps.Session,
			OutputKey: "response",
		}, nil

	case workflow.NodeTypeVote:
		threshold, _ := node.Properties["threshold"].(float64)
		voteType, _ := node.Properties["vote_type"].(string)
		return &VoteProcessor{
			Threshold: threshold,
			VoteType:  voteType,
		}, nil

	case workflow.NodeTypeLoop:
		maxRounds, _ := node.Properties["max_rounds"].(float64)
		exitOnScore, _ := node.Properties["exit_on_score"].(float64)
		return &LoopProcessor{
			MaxRounds:   int(maxRounds),
			ExitOnScore: int(exitOnScore),
			Session:     deps.Session,
		}, nil

	case workflow.NodeTypeFactCheck:
		threshold, _ := node.Properties["verify_threshold"].(float64)
		var provider llm.LLMProvider
		var err error
		if f.Registry != nil {
			provider, err = f.Registry.GetLLMProvider("default")
			if err != nil {
				return nil, fmt.Errorf("failed to get default LLM for fact check: %w", err)
			}
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
		return NewMemoryRetrievalProcessor(f.MemoryManager), nil

	case workflow.NodeTypeContextSynth:
		maxRecent, _ := node.Properties["max_recent_rounds"].(float64)
		if maxRecent == 0 {
			maxRecent = 3
		}
		return &ContextSynthesizerProcessor{
			MaxRecentRounds: int(maxRecent),
		}, nil

	default:
		return nil, fmt.Errorf("unsupported node type: %s", node.Type)
	}
}
