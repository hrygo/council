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
	NodeID          string // Graph node ID (e.g., "agent_affirmative")
	AgentID         string // Agent UUID from database
	AgentRepo       agent.Repository
	Registry        *llm.Registry
	Tools           []tools.Tool             // Injected tools
	Session         *workflow.Session        // Injected Session
	PassthroughKeys []string                 // Configuration: Keys to pass to output
	PromptSections  []workflow.PromptSection // Configuration: Input keys to build prompt
	OutputKey       string                   // Configuration: Key for response content (e.g. "agent_output")
}

func (a *AgentProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Fetch Agent Persona
	_ = stream // usage in tool/tokens below

	// 2. Fetch Agent Persona
	ag, err := a.AgentRepo.GetByID(ctx, parseUUID(a.AgentID))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent %s: %w", a.AgentID, err)
	}

	// 3. Construct Context from Input
	history := constructHistory(ag.PersonaPrompt, input, a.PromptSections)

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
			Stream:      true,
			Tools:       llmTools,
		}
		if req.Model == "" {
			req.Model = a.Registry.GetDefaultModel()
		}

		// Notify "Thinking" (Force frontend to render message bubble)
		stream <- workflow.StreamEvent{
			Type:      "token_stream",
			Timestamp: time.Now(),
			Data:      map[string]interface{}{"node_id": a.NodeID, "agent_id": a.AgentID, "chunk": " "},
		}

		resp, err := a.streamResponse(ctx, provider, req, stream)
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
						"node_id":  a.NodeID,
						"agent_id": a.AgentID,
						"tool":     toolName,
						"input":    toolArgs,
						"output":   result,
					},
				}
			}
			// Continue loop to let LLM process tool result
		} else {
			// No tool calls, we are done
			finalResponse = resp.Content

			// Notify Content Stream (Already done by streamResponse)

			// Notify Token Usage
			stream <- workflow.StreamEvent{
				Type:      "token_usage",
				Timestamp: time.Now(),
				Data: map[string]interface{}{
					"node_id":            a.NodeID,
					"agent_id":           a.AgentID,
					"input_tokens":       resp.Usage.PromptTokens,
					"output_tokens":      resp.Usage.CompletionTokens,
					"estimated_cost_usd": estimateCost(resp.Usage.TotalTokens),
				},
			}
			break
		}
	}

	// 6. Output with context passthrough (Config Driven)
	outputKey := a.OutputKey
	if outputKey == "" {
		outputKey = "agent_output"
	}
	output := map[string]interface{}{
		outputKey:   finalResponse,
		"agent_id":  a.AgentID,
		"timestamp": time.Now(),
	}

	workflow.ApplyPassthrough(input, output, workflow.PassthroughConfig{
		Keys: a.PassthroughKeys,
	})

	finalOutputKeys := a.OutputKey
	_ = finalOutputKeys // end logic complete

	return output, nil
}

func (a *AgentProcessor) streamResponse(ctx context.Context, provider llm.LLMProvider, req *llm.CompletionRequest, stream chan<- workflow.StreamEvent) (*llm.CompletionResponse, error) {
	chunkChan, errChan := provider.Stream(ctx, req)

	fullContent := ""
	toolCallsMap := make(map[int]*llm.ToolCall)
	var usage llm.Usage

	for {
		select {
		case chunk, ok := <-chunkChan:
			if !ok {
				chunkChan = nil
			} else {
				if chunk.Content != "" {
					stream <- workflow.StreamEvent{
						Type:      "token_stream",
						Timestamp: time.Now(),
						Data:      map[string]interface{}{"node_id": a.NodeID, "agent_id": a.AgentID, "chunk": chunk.Content},
					}
					fullContent += chunk.Content
				}

				for _, tc := range chunk.ToolCalls {
					index := tc.Index
					if _, exists := toolCallsMap[index]; !exists {
						toolCallsMap[index] = &llm.ToolCall{
							Index:    index,
							ID:       tc.ID,
							Type:     tc.Type,
							Function: llm.FunctionCall{Name: tc.Function.Name},
						}
					}
					current := toolCallsMap[index]
					if tc.ID != "" {
						current.ID = tc.ID
					}
					if tc.Type != "" {
						current.Type = tc.Type
					}
					if tc.Function.Name != "" {
						current.Function.Name = tc.Function.Name
					}
					current.Function.Arguments += tc.Function.Arguments
				}

				if chunk.Usage != nil {
					usage = *chunk.Usage
				}
			}
		case err, ok := <-errChan:
			if ok {
				return nil, err
			}
			errChan = nil
		}
		if chunkChan == nil && errChan == nil {
			break
		}
	}

	// Flatten tool calls
	var toolCalls []llm.ToolCall
	var indices []int
	for i := range toolCallsMap {
		indices = append(indices, i)
	}
	sort.Ints(indices)
	for _, i := range indices {
		toolCalls = append(toolCalls, *toolCallsMap[i])
	}

	return &llm.CompletionResponse{
		Content:   fullContent,
		ToolCalls: toolCalls,
		Usage:     usage,
	}, nil
}

func constructHistory(systemPrompt string, input map[string]interface{}, sections []workflow.PromptSection) []llm.Message {
	var contextBuilder strings.Builder

	// Build structured context with configured sections
	for _, sec := range sections {
		if val, ok := input[sec.Key].(string); ok && val != "" {
			contextBuilder.WriteString(fmt.Sprintf("<%s>\n%s\n</%s>\n\n", sec.Label, val, sec.Label))
		}
	}

	userContent := contextBuilder.String()
	if userContent == "" {
		// Fallback: dump all string values if no structured content found
		keys := make([]string, 0, len(input))
		for k := range input {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if str, ok := input[k].(string); ok {
				contextBuilder.WriteString(fmt.Sprintf("%s: %s\n", k, str))
			}
		}
		userContent = contextBuilder.String()
	}

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

// estimateCost provides a rough cost estimate based on token count
// Using approximate pricing of $0.002 per 1K tokens (typical for most models)
func estimateCost(tokens int) float64 {
	return float64(tokens) * 0.002 / 1000
}
