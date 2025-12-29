package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/core/workflow/nodes/tools"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type AgentProcessor struct {
	NodeID    string // Graph node ID (e.g., "agent_affirmative")
	AgentID   string // Agent UUID from database
	AgentRepo agent.Repository
	Registry  *llm.Registry
	Tools     []tools.Tool      // Injected tools
	Session   *workflow.Session // Injected Session
}

func (a *AgentProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Notify Start
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": a.NodeID, "status": "running"},
	}

	// 2. Fetch Agent Persona
	ag, err := a.AgentRepo.GetByID(ctx, parseUUID(a.AgentID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent %s: %w", a.AgentID, err)
	}

	// 3. Construct Context from Input
	history := constructHistory(ag.PersonaPrompt, input)

	// Prepare Tools
	var llmTools []llm.Tool
	for _, t := range a.Tools {
		llmTools = append(llmTools, llm.Tool{
			Type: "function",
			Function: llm.ToolFunction{
				Name:        t.Name(),
				Description: t.Description(),
				Parameters:  t.Parameters(),
			},
		})
	}

	// 4. Resolve LLM Provider
	providerName := ag.ModelConfig.Provider
	provider, err := a.Registry.GetLLMProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM provider '%s': %w", providerName, err)
	}

	// 5. Re-Act Loop
	var finalResponse string
	maxIterations := 5

	for i := 0; i < maxIterations; i++ {
		req := &llm.CompletionRequest{
			Model:       ag.ModelConfig.Model,
			Messages:    history,
			Temperature: float32(ag.ModelConfig.Temperature),
			MaxTokens:   ag.ModelConfig.MaxTokens,
			TopP:        float32(ag.ModelConfig.TopP),
			Stream:      false,
			Tools:       llmTools,
		}
		if req.Model == "" {
			req.Model = a.Registry.GetDefaultModel()
		}

		resp, err := provider.Generate(ctx, req)
		if err != nil {
			return nil, err
		}

		// Append Assistant Message
		msg := llm.Message{
			Role:      "assistant",
			Content:   resp.Content,
			ToolCalls: resp.ToolCalls,
		}
		history = append(history, msg)

		// Check if we need to execute tools
		if len(resp.ToolCalls) > 0 {
			finalResponse = resp.Content // Could be empty or partial

			// Execute Tools
			for _, tc := range resp.ToolCalls {
				toolName := tc.Function.Name
				toolArgs := tc.Function.Arguments

				// Find tool
				var selectedTool tools.Tool
				for _, t := range a.Tools {
					if t.Name() == toolName {
						selectedTool = t
						break
					}
				}

				var result string
				if selectedTool == nil {
					result = fmt.Sprintf("Error: Tool %s not found", toolName)
				} else {
					// Parse Args
					var argsMap map[string]interface{}
					if err := json.Unmarshal([]byte(toolArgs), &argsMap); err != nil {
						result = fmt.Sprintf("Error: Invalid JSON arguments: %v", err)
					} else {
						// Execute
						if sat, ok := selectedTool.(tools.SessionAwareTool); ok {
							result, err = sat.ExecuteWithSession(ctx, a.Session, argsMap)
						} else {
							result, err = selectedTool.Execute(ctx, argsMap)
						}
						if err != nil {
							result = fmt.Sprintf("Error: %v", err)
						}
					}
				}

				// Append Tool Result
				history = append(history, llm.Message{
					Role:       "tool",
					Content:    result,
					ToolCallID: tc.ID,
				})

				// Notify Stream
				stream <- workflow.StreamEvent{
					Type:      "tool_execution",
					Timestamp: time.Now(),
					Data: map[string]interface{}{
						"node_id": a.NodeID,
						"tool":    toolName,
						"input":   toolArgs,
						"output":  result,
					},
				}
			}
			// Continue loop to let LLM process tool result
		} else {
			// No tool calls, we are done
			finalResponse = resp.Content

			// Notify Content Stream (Simulated for non-streaming)
			stream <- workflow.StreamEvent{
				Type:      "token_stream",
				Timestamp: time.Now(),
				Data:      map[string]interface{}{"node_id": a.NodeID, "chunk": finalResponse},
			}
			break
		}
	}

	// 6. Output
	output := map[string]interface{}{
		"agent_output": finalResponse,
		"agent_id":     a.AgentID,
		"timestamp":    time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": a.NodeID, "status": "completed"},
	}

	return output, nil
}

func constructHistory(systemPrompt string, input map[string]interface{}) []llm.Message {
	keys := make([]string, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var contextBuilder strings.Builder
	for _, k := range keys {
		if str, ok := input[k].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("%s: %s\n", k, str))
		}
	}
	userContent := contextBuilder.String()
	if userContent == "" {
		userContent = "Begin task."
	}

	return []llm.Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}
}

func parseUUID(id string) uuid.UUID {
	u, _ := uuid.Parse(id)
	return u
}
