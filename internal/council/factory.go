package council

import (
	"fmt"

	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/memory"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type CouncilNodeFactory struct {
	AgentRepo     agent.Repository
	Registry      *llm.Registry
	MemoryManager memory.MemoryManager
	baseFactory   *nodes.GenericNodeFactory
}

func NewCouncilNodeFactory(repo agent.Repository, registry *llm.Registry, memManager memory.MemoryManager) *CouncilNodeFactory {
	return &CouncilNodeFactory{
		AgentRepo:     repo,
		Registry:      registry,
		MemoryManager: memManager,
		baseFactory:   nodes.NewGenericNodeFactory(registry, repo, memManager),
	}
}

func (f *CouncilNodeFactory) CreateNode(node *workflow.Node, deps workflow.FactoryDeps) (workflow.NodeProcessor, error) {
	switch node.Type {
	case workflow.NodeTypeStart:
		// Map Council StartOutputKeys to StartProcessor configuration
		return &StartProcessor{
			OutputKeys: StartOutputKeys,
		}, nil

	case workflow.NodeTypeEnd:
		// Map EndInputKeys to PromptSections
		sections := []workflow.PromptSection{}
		for _, key := range EndInputKeys {
			sections = append(sections, workflow.PromptSection{Key: key, Label: key})
		}

		// We need parameters from node properties
		model, _ := node.Properties["model"].(string)
		prompt, _ := node.Properties["system_prompt"].(string)

		// Resolve default model if not specified
		if model == "" {
			model = f.Registry.GetDefaultModel()
		}

		// Get LLM Provider
		var llmProvider llm.LLMProvider
		var err error
		if model != "" {
			// Try to get provider by model name
			llmProvider, err = f.Registry.GetProviderByModel(model)
		}

		// Fallback to provider property or default deepseek
		if llmProvider == nil || err != nil {
			providerName, _ := node.Properties["provider"].(string)
			if providerName == "" {
				providerName = "default"
			}
			llmProvider, err = f.Registry.GetLLMProvider(providerName)
		}

		if err != nil {
			return nil, fmt.Errorf("failed to resolve llm provider for end node: %w", err)
		}

		return &nodes.EndProcessor{
			NodeID:         node.ID,
			LLM:            llmProvider,
			Model:          model,
			Prompt:         prompt,
			PromptSections: sections,
			OutputKey:      "final_report", // Council-specific key
		}, nil

	case workflow.NodeTypeLoop:
		maxRounds, _ := node.Properties["max_rounds"].(int)
		if maxRounds == 0 {
			if f, ok := node.Properties["max_rounds"].(float64); ok {
				maxRounds = int(f)
			}
		}

		exitOnScore, _ := node.Properties["exit_on_score"].(int)
		if exitOnScore == 0 {
			if f, ok := node.Properties["exit_on_score"].(float64); ok {
				exitOnScore = int(f)
			}
		}

		return &nodes.LoopProcessor{
			MaxRounds:       maxRounds,
			ExitOnScore:     exitOnScore,
			Session:         deps.Session,
			PassthroughKeys: LoopPassthroughKeys,
		}, nil

	case workflow.NodeTypeAgent:
		// Create AgentProcessor with Council config
		agentID, _ := node.Properties["agent_uuid"].(string)
		if agentID == "" {
			agentID, _ = node.Properties["agent_id"].(string)
		}
		if agentID == "" {
			return nil, fmt.Errorf("agent node %s missing agent_uuid or agent_id", node.ID)
		}

		// Construct generic PromptSections for Agent from Council context keys
		// We use a fixed set of sections for all Council agents for now
		sections := []workflow.PromptSection{
			{Key: "document_content", Label: "document_content"},
			{Key: "proposal", Label: "proposal"},
			{Key: "combined_context", Label: "combined_context"},
			{Key: "aggregated_outputs", Label: "previous_analyses"},
			{Key: "optimization_objective", Label: "optimization_objective"},
		}

		return &nodes.AgentProcessor{
			NodeID:          node.ID,
			AgentID:         agentID,
			AgentRepo:       f.AgentRepo,
			Registry:        f.Registry,
			Session:         deps.Session,
			PassthroughKeys: AgentPassthroughKeys,
			PromptSections:  sections,
			OutputKey:       "agent_output", // Council-specific key
			// Tools: f.resolveTools(node) // TODO: Implement tool resolution
		}, nil

	case workflow.NodeTypeParallel:
		// Parallel logic is handled by Engine (structural), but we return nil/error
		// so engine knows to handle it or we can return a Dummy processor if Engine requires one.
		// Our generic engine handles Parallel explicitly.
		return nil, nil // Or generic error "Handled by Engine"

	default:
		// Delegate generic nodes to core factory
		return f.baseFactory.CreateNode(node, deps)
	}
}
