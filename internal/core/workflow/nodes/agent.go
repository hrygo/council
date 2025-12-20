package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrygo/council/internal/core/agent"
	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type AgentProcessor struct {
	AgentID   string
	AgentRepo agent.Repository
	Registry  *llm.Registry
}

func (a *AgentProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Notify Start
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": a.AgentID, "status": "running"},
	}

	// 2. Fetch Agent Persona
	ag, err := a.AgentRepo.GetByID(ctx, parseUUID(a.AgentID)) // Assuming helper or proper UUID handling
	if err != nil {
		return nil, fmt.Errorf("failed to fetch agent %s: %w", a.AgentID, err)
	}

	// 3. Construct Context from Input
	// Simplified: Join all string values
	var contextBuilder strings.Builder
	for k, v := range input {
		if str, ok := v.(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("%s: %s\n", k, str))
		}
	}
	userContent := contextBuilder.String()
	if userContent == "" {
		userContent = "Begin task."
	}

	// 4. Resolve LLM Provider
	providerName := ag.ModelConfig.Provider
	provider, err := a.Registry.GetLLMProvider(providerName)
	if err != nil {
		return nil, fmt.Errorf("failed to get LLM provider '%s': %w", providerName, err)
	}

	// 5. Call LLM
	req := &llm.CompletionRequest{
		Model: ag.ModelConfig.Model,
		Messages: []llm.Message{
			{Role: "system", Content: ag.PersonaPrompt},
			{Role: "user", Content: userContent},
		},
		Temperature: float32(ag.ModelConfig.Temperature),
		MaxTokens:   ag.ModelConfig.MaxTokens,
		TopP:        float32(ag.ModelConfig.TopP),
		Stream:      true,
	}
	// Fallback model if config missing
	if req.Model == "" {
		req.Model = "gpt-4"
	}

	tokenStream, errChan := provider.Stream(ctx, req)
	var responseBuilder strings.Builder

	// 5. Stream Tokens
Loop:
	for {
		select {
		case token, ok := <-tokenStream:
			if !ok {
				tokenStream = nil
			} else {
				responseBuilder.WriteString(token)
				stream <- workflow.StreamEvent{
					Type:      "token_stream",
					Timestamp: time.Now(),
					Data:      map[string]interface{}{"node_id": a.AgentID, "chunk": token},
				}
			}
		case err, ok := <-errChan:
			if ok && err != nil {
				return nil, err
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		if tokenStream == nil {
			break Loop
		}
	}

	// 6. Output
	output := map[string]interface{}{
		"agent_output": responseBuilder.String(),
		"agent_id":     a.AgentID,
		"timestamp":    time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": a.AgentID, "status": "completed"},
	}

	return output, nil
}

func parseUUID(id string) uuid.UUID {
	u, _ := uuid.Parse(id)
	return u
}
