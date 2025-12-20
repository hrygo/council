package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// SystemNamespace is the UUID namespace for system resources
var SystemNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

// Seeder handles database seeding for default data.
type Seeder struct {
	db *pgxpool.Pool
}

// NewSeeder creates a new Seeder instance.
func NewSeeder(db *pgxpool.Pool) *Seeder {
	return &Seeder{db: db}
}

// SeedAgents seeds default agents into the database.
// Uses ON CONFLICT DO NOTHING for idempotency.
func (s *Seeder) SeedAgents(ctx context.Context) error {
	prompts, err := LoadAllPrompts()
	if err != nil {
		return fmt.Errorf("failed to load prompts: %w", err)
	}

	for agentID, prompt := range prompts {
		// Generate deterministic UUID from agent ID string
		agentUUID := uuid.NewSHA1(SystemNamespace, []byte(agentID))

		modelConfig, err := json.Marshal(map[string]interface{}{
			"provider":    prompt.Config.Provider,
			"model":       prompt.Config.Model,
			"temperature": prompt.Config.Temperature,
			"max_tokens":  prompt.Config.MaxTokens,
			"top_p":       prompt.Config.TopP,
		})
		if err != nil {
			return fmt.Errorf("failed to marshal model config for %s: %w", agentID, err)
		}

		capabilities, _ := json.Marshal(map[string]bool{
			"web_search":     false,
			"code_execution": false,
		})

		_, err = s.db.Exec(ctx, `
			INSERT INTO agents (id, name, persona_prompt, model_config, capabilities, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, agentUUID, prompt.Config.Name, prompt.Content, modelConfig, capabilities)

		if err != nil {
			return fmt.Errorf("failed to seed agent %s: %w", agentID, err)
		}
	}

	return nil
}

// councilSystemPrompt is the system prompt for The Council group.
const councilSystemPrompt = `# The Council - 多智能体协作治理体

你是 "理事会" (The Council) 的一部分，这是一个由多个 AI 专家组成的治理机构。

## 核心原则

1. **对抗性协作**: 通过正反辩论锻造卓越决策。
2. **工业级标准**: 所有输出必须具备可执行性。
3. **全局统筹**: 始终以用户的"初始目标"为最高准则。
`

// SeedGroups seeds default groups into the database.
func (s *Seeder) SeedGroups(ctx context.Context) error {
	defaultAgentIDs := []string{
		"system_affirmative",
		"system_negative",
		"system_adjudicator",
	}

	var agentUUIDs []string
	for _, id := range defaultAgentIDs {
		agentUUIDs = append(agentUUIDs, uuid.NewSHA1(SystemNamespace, []byte(id)).String())
	}

	// Generate deterministic UUID for the group
	groupID := "system_council"
	groupUUID := uuid.NewSHA1(SystemNamespace, []byte(groupID))

	// Convert agent IDs to JSON array for storage
	agentIDsJSON, _ := json.Marshal(agentUUIDs)

	_, err := s.db.Exec(ctx, `
		INSERT INTO groups (id, name, system_prompt, default_agent_ids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (id) DO NOTHING
	`, groupUUID, "The Council", councilSystemPrompt, agentIDsJSON)

	if err != nil {
		return fmt.Errorf("failed to seed group: %w", err)
	}

	return nil
}

// debateWorkflowGraph is the JSON graph definition for council_debate workflow.
const debateWorkflowGraph = `{
	"start_node_id": "start",
	"nodes": {
		"start": {
			"id": "start",
			"type": "start",
			"name": "Input Document",
			"next_ids": ["parallel_analysis"]
		},
		"parallel_analysis": {
			"id": "parallel_analysis",
			"type": "parallel",
			"name": "Parallel Analysis",
			"next_ids": ["agent_affirmative", "agent_negative"]
		},
		"agent_affirmative": {
			"id": "agent_affirmative",
			"type": "agent",
			"name": "Affirmative",
			"properties": {"agent_id": "system_affirmative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_negative": {
			"id": "agent_negative",
			"type": "agent",
			"name": "Negative",
			"properties": {"agent_id": "system_negative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_adjudicator": {
			"id": "agent_adjudicator",
			"type": "agent",
			"name": "Adjudicator",
			"properties": {"agent_id": "system_adjudicator"},
			"next_ids": ["end"]
		},
		"end": {
			"id": "end",
			"type": "end",
			"name": "Generate Report"
		}
	}
}`

// optimizeWorkflowGraph is the JSON graph definition for council_optimize workflow.
const optimizeWorkflowGraph = `{
	"start_node_id": "start",
	"nodes": {
		"start": {
			"id": "start",
			"type": "start",
			"name": "Input Document",
			"next_ids": ["memory_retrieval"]
		},
		"memory_retrieval": {
			"id": "memory_retrieval",
			"type": "memory_retrieval",
			"name": "Load History",
			"properties": {
				"max_results": 5,
				"time_range_days": 7,
				"include_verdicts": true
			},
			"next_ids": ["parallel_debate"]
		},
		"parallel_debate": {
			"id": "parallel_debate",
			"type": "parallel",
			"name": "Debate Round",
			"next_ids": ["agent_affirmative", "agent_negative"]
		},
		"agent_affirmative": {
			"id": "agent_affirmative",
			"type": "agent",
			"name": "Affirmative",
			"properties": {"agent_id": "system_affirmative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_negative": {
			"id": "agent_negative",
			"type": "agent",
			"name": "Negative",
			"properties": {"agent_id": "system_negative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_adjudicator": {
			"id": "agent_adjudicator",
			"type": "agent",
			"name": "Adjudicator",
			"properties": {"agent_id": "system_adjudicator", "output_format": "structured_verdict"},
			"next_ids": ["human_review"]
		},
		"human_review": {
			"id": "human_review",
			"type": "human_review",
			"name": "Review & Apply",
			"properties": {"show_score": true, "actions": ["continue", "apply", "exit", "rollback"]},
			"next_ids": ["loop_decision"]
		},
		"loop_decision": {
			"id": "loop_decision",
			"type": "loop",
			"name": "Continue?",
			"properties": {"max_rounds": 5, "exit_on_score": 90},
			"next_ids": ["memory_retrieval", "end"]
		},
		"end": {
			"id": "end",
			"type": "end",
			"name": "Final Report"
		}
	}
}`

// SeedWorkflows seeds default workflow templates into the database.
func (s *Seeder) SeedWorkflows(ctx context.Context) error {
	workflows := []struct {
		ID          string
		Name        string
		Description string
		Graph       string
	}{
		{
			ID:          "council_debate",
			Name:        "Council Debate",
			Description: "三方辩论，生成综合裁决报告",
			Graph:       debateWorkflowGraph,
		},
		{
			ID:          "council_optimize",
			Name:        "Council Optimize",
			Description: "迭代优化循环，含历史上下文检索",
			Graph:       optimizeWorkflowGraph,
		},
	}

	for _, wf := range workflows {
		// Generate deterministic UUID for workflow
		wfUUID := uuid.NewSHA1(SystemNamespace, []byte(wf.ID))

		// Replace agent IDs in graph with UUIDs
		graph := wf.Graph
		graph = strings.ReplaceAll(graph, "system_affirmative", uuid.NewSHA1(SystemNamespace, []byte("system_affirmative")).String())
		graph = strings.ReplaceAll(graph, "system_negative", uuid.NewSHA1(SystemNamespace, []byte("system_negative")).String())
		graph = strings.ReplaceAll(graph, "system_adjudicator", uuid.NewSHA1(SystemNamespace, []byte("system_adjudicator")).String())

		// Compact the JSON to ensure valid format
		var compactGraph bytes.Buffer
		if err := json.Compact(&compactGraph, []byte(graph)); err != nil {
			return fmt.Errorf("invalid workflow graph JSON for %s: %w", wf.ID, err)
		}

		_, err := s.db.Exec(ctx, `
			INSERT INTO workflow_templates (id, name, description, graph_definition, created_at, updated_at)
			VALUES ($1, $2, $3, $4::jsonb, NOW(), NOW())
			ON CONFLICT (id) DO NOTHING
		`, wfUUID, wf.Name, wf.Description, compactGraph.String())

		if err != nil {
			return fmt.Errorf("failed to seed workflow %s: %w", wf.ID, err)
		}
	}

	return nil
}

// SeedAll seeds all default data: agents, groups, and workflows.
// Executes in dependency order: Agents -> Groups -> Workflows.
func (s *Seeder) SeedAll(ctx context.Context) error {
	if err := s.SeedAgents(ctx); err != nil {
		return fmt.Errorf("seed agents: %w", err)
	}
	if err := s.SeedGroups(ctx); err != nil {
		return fmt.Errorf("seed groups: %w", err)
	}
	if err := s.SeedWorkflows(ctx); err != nil {
		return fmt.Errorf("seed workflows: %w", err)
	}
	return nil
}
