package resources

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/hrygo/council/internal/pkg/config"
)

// SystemNamespace is the UUID namespace for system resources
var SystemNamespace = uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")

// Seeder handles database seeding for default data.
type Seeder struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

// NewSeeder creates a new Seeder instance.
func NewSeeder(db *pgxpool.Pool, cfg *config.Config) *Seeder {
	return &Seeder{db: db, cfg: cfg}
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

		// Use prompt config if available, otherwise fallback to system config
		provider := prompt.Config.Provider
		if provider == "" {
			provider = s.cfg.LLM.Provider
		}
		model := prompt.Config.Model
		if model == "" {
			model = s.cfg.LLM.Model
		}

		modelConfig, err := json.Marshal(map[string]interface{}{
			"provider":    provider,
			"model":       model,
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
			INSERT INTO agents (agent_uuid, name, persona_prompt, model_config, capabilities, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			ON CONFLICT (agent_uuid) DO UPDATE SET
				name = EXCLUDED.name,
				persona_prompt = EXCLUDED.persona_prompt,
				model_config = EXCLUDED.model_config,
				updated_at = NOW()
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
		"system_surgeon",
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
		INSERT INTO groups (group_uuid, name, system_prompt, default_agent_uuids, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		ON CONFLICT (group_uuid) DO NOTHING
	`, groupUUID, "The Council", councilSystemPrompt, agentIDsJSON)

	if err != nil {
		return fmt.Errorf("failed to seed group: %w", err)
	}

	return nil
}

// debateWorkflowGraph is the JSON graph definition for 'Council Debate' workflow.
const debateWorkflowGraph = `{
	"workflow_uuid": "c0deb47e-0000-0000-0000-000000000001", "start_node_id": "start",
	"nodes": {
		"start": {
			"node_id": "start",
			"type": "start",
			"name": "Input Document",
			"next_ids": ["parallel_analysis"]
		},
		"parallel_analysis": {
			"node_id": "parallel_analysis",
			"type": "parallel",
			"name": "Parallel Analysis",
			"next_ids": ["agent_affirmative", "agent_negative"]
		},
		"agent_affirmative": {
			"node_id": "agent_affirmative",
			"type": "agent",
			"name": "Affirmative",
			"properties": {"agent_uuid": "system_affirmative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_negative": {
			"node_id": "agent_negative",
			"type": "agent",
			"name": "Negative",
			"properties": {"agent_uuid": "system_negative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_adjudicator": {
			"node_id": "agent_adjudicator",
			"type": "agent",
			"name": "Adjudicator",
			"properties": {"agent_uuid": "system_adjudicator"},
			"next_ids": ["end"]
		},
		"end": {
			"node_id": "end",
			"type": "end",
			"name": "Generate Report"
		}
	}
}`

// optimizeWorkflowGraph is the JSON graph definition for 'Council Optimize' workflow.
const optimizeWorkflowGraph = `{
	"workflow_uuid": "c00p71m3-0000-0000-0000-000000000001", "start_node_id": "start",
	"nodes": {
		"start": {
			"node_id": "start",
			"type": "start",
			"name": "Input Document",
			"next_ids": ["memory_retrieval"]
		},
		"memory_retrieval": {
			"node_id": "memory_retrieval",
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
			"node_id": "parallel_debate",
			"type": "parallel",
			"name": "Debate Round",
			"next_ids": ["agent_affirmative", "agent_negative"]
		},
		"agent_affirmative": {
			"node_id": "agent_affirmative",
			"type": "agent",
			"name": "Affirmative",
			"properties": {"agent_uuid": "system_affirmative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_negative": {
			"node_id": "agent_negative",
			"type": "agent",
			"name": "Negative",
			"properties": {"agent_uuid": "system_negative"},
			"next_ids": ["agent_adjudicator"]
		},
		"agent_adjudicator": {
			"node_id": "agent_adjudicator",
			"type": "agent",
			"name": "Adjudicator",
			"properties": {"agent_uuid": "system_adjudicator", "output_format": "structured_verdict"},
			"next_ids": ["agent_surgeon"]
		},
		"agent_surgeon": {
			"node_id": "agent_surgeon",
			"type": "agent",
			"name": "Surgeon",
			"properties": {"agent_uuid": "system_surgeon", "tools": ["write_file", "read_file"]},
			"next_ids": ["human_review"]
		},
		"human_review": {
			"node_id": "human_review",
			"type": "human_review",
			"name": "Review & Apply",
			"properties": {"show_score": true, "actions": ["continue", "apply", "exit", "rollback"]},
			"next_ids": ["loop_decision"]
		},
		"loop_decision": {
			"node_id": "loop_decision",
			"type": "loop",
			"name": "Continue?",
			"properties": {"max_rounds": 5, "exit_on_score": 90},
			"next_ids": ["memory_retrieval", "end"]
		},
		"end": {
			"node_id": "end",
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
			ID:          "c0deb47e-0000-0000-0000-000000000001",
			Name:        "Council Debate",
			Description: "三方辩论，生成综合裁决报告",
			Graph:       debateWorkflowGraph,
		},
		{
			ID:          "c00p71m3-0000-0000-0000-000000000001",
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
		graph = strings.ReplaceAll(graph, "system_surgeon", uuid.NewSHA1(SystemNamespace, []byte("system_surgeon")).String())

		// Compact the JSON to ensure valid format
		var compactGraph bytes.Buffer
		if err := json.Compact(&compactGraph, []byte(graph)); err != nil {
			return fmt.Errorf("invalid workflow graph JSON for %s: %w", wf.ID, err)
		}

		_, err := s.db.Exec(ctx, `
			INSERT INTO workflow_templates (template_uuid, name, description, graph_definition, is_system, created_at, updated_at)
			VALUES ($1, $2, $3, $4::jsonb, true, NOW(), NOW())
			ON CONFLICT (template_uuid) DO UPDATE SET
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				graph_definition = EXCLUDED.graph_definition,
				is_system = EXCLUDED.is_system,
				updated_at = NOW()
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
	if err := s.SeedLLMOptions(ctx); err != nil {
		return fmt.Errorf("seed llm options: %w", err)
	}
	return nil
}

// SeedLLMOptions seeds the LLM providers and models from 2025 verified list.
func (s *Seeder) SeedLLMOptions(ctx context.Context) error {
	// 1. Create Tables if not exist (Migration replacement for robustness)
	_, err := s.db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS llm_providers (
			provider_id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			icon TEXT NOT NULL,
			sort_order INTEGER NOT NULL DEFAULT 0,
			is_enabled BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS llm_models (
			model_id TEXT PRIMARY KEY,
			provider_id TEXT NOT NULL REFERENCES llm_providers(provider_id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			is_mainstream BOOLEAN DEFAULT FALSE,
			sort_order INTEGER NOT NULL DEFAULT 0,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure LLM tables: %w", err)
	}

	// 2. Define Data
	providers := []struct {
		ID        string
		Name      string
		Icon      string
		SortOrder int
		Models    []string
	}{
		{
			ID: "openai", Name: "OpenAI", Icon: "🟢", SortOrder: 1,
			Models: []string{"gpt-4o", "o1", "gpt-5-mini", "gpt-4.5-preview"},
		},
		{
			ID: "google", Name: "Google", Icon: "🔵", SortOrder: 2,
			Models: []string{"gemini-3-pro", "gemini-3-flash", "gemini-2.0-flash", "gemini-2.0-pro"},
		},
		{
			ID: "deepseek", Name: "DeepSeek", Icon: "🟣", SortOrder: 3,
			Models: []string{"deepseek-chat", "deepseek-reasoner"}, // verified: deepseek-chat (V3), deepseek-reasoner (R1)
		},
		{
			ID: "dashscope", Name: "DashScope", Icon: "🟡", SortOrder: 4,
			Models: []string{"qwen-max", "qwen-plus", "qwen-turbo"},
		},
		{
			ID: "siliconflow", Name: "SiliconFlow", Icon: "🟠", SortOrder: 5,
			Models: []string{
				"zai-org/GLM-4.6", // Verified: SiliconFlow uses repo format
				"Qwen/Qwen2.5-72B-Instruct",
				"Qwen/Qwen2.5-Coder-32B-Instruct",
				"deepseek-ai/DeepSeek-V3",
				"deepseek-ai/DeepSeek-R1",
			},
		},
	}

	// 3. Upsert Logic
	for _, p := range providers {
		// Upsert Provider
		_, err := s.db.Exec(ctx, `
			INSERT INTO llm_providers (provider_id, name, icon, sort_order, updated_at)
			VALUES ($1, $2, $3, $4, NOW())
			ON CONFLICT (provider_id) DO UPDATE SET
				name = EXCLUDED.name,
				icon = EXCLUDED.icon,
				sort_order = EXCLUDED.sort_order,
				updated_at = NOW()
		`, p.ID, p.Name, p.Icon, p.SortOrder)
		if err != nil {
			return fmt.Errorf("failed to upsert provider %s: %w", p.ID, err)
		}

		// Upsert Models
		for i, mID := range p.Models {
			_, err := s.db.Exec(ctx, `
				INSERT INTO llm_models (model_id, provider_id, name, is_mainstream, sort_order, updated_at)
				VALUES ($1, $2, $3, true, $4, NOW())
				ON CONFLICT (model_id) DO UPDATE SET
					name = EXCLUDED.name,
					is_mainstream = true,
					sort_order = EXCLUDED.sort_order,
					updated_at = NOW()
			`, mID, p.ID, mID, i+1)
			if err != nil {
				return fmt.Errorf("failed to upsert model %s for %s: %w", mID, p.ID, err)
			}
		}
	}

	return nil
}
