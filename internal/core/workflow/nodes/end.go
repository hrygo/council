package nodes

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

type EndProcessor struct {
	LLM    llm.LLMProvider
	Model  string
	Prompt string
}

func (e *EndProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Notify Start
	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": "end", "status": "running"},
	}

	// 2. Aggregate Content
	// Simplification: Try to find 'combined_context' or 'proposal' or dump everything
	var contentBuilder strings.Builder
	if val, ok := input["combined_context"].(string); ok {
		contentBuilder.WriteString("Context:\n")
		contentBuilder.WriteString(val)
		contentBuilder.WriteString("\n\n")
	}
	if val, ok := input["proposal"].(string); ok {
		contentBuilder.WriteString("Proposal:\n")
		contentBuilder.WriteString(val)
		contentBuilder.WriteString("\n\n")
	}
	// Fallback/Supplement: Check for agent outputs in input map?
	// For now, assuming input carries the necessary context string.

	fullContent := contentBuilder.String()
	if fullContent == "" {
		fullContent = fmt.Sprintf("Raw Input: %v", input)
	}

	// 3. Call LLM
	prompt := e.Prompt
	if prompt == "" {
		prompt = "Please summarize the above content."
	}

	req := &llm.CompletionRequest{
		Model: e.Model,
		Messages: []llm.Message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: fullContent},
		},
		Temperature: 0.7,
		Stream:      true, // We want streaming
	}

	tokenStream, errChan := e.LLM.Stream(ctx, req)

	var finalSummary strings.Builder

	// 4. Stream Tokens
	// We need to read from both channels.
	// Since the interface returns (<-chan string, <-chan error), we must loop.
Loop:
	for {
		select {
		case chunk, ok := <-tokenStream:
			if !ok {
				tokenStream = nil // Channel closed
			} else {
				if chunk.Content != "" {
					finalSummary.WriteString(chunk.Content)
					stream <- workflow.StreamEvent{
						Type:      "token_stream",
						Timestamp: time.Now(),
						Data:      map[string]interface{}{"node_id": "end", "chunk": chunk.Content},
					}
				}
			}
		case err, ok := <-errChan:
			if ok && err != nil {
				return nil, fmt.Errorf("llm stream error: %w", err)
			}
			// If errChan is closed or nil, ignore?
			// Usually check tokenStream nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		if tokenStream == nil {
			break Loop
		}
	}

	// 5. Output
	output := map[string]interface{}{
		"final_report": finalSummary.String(),
		"ended_at":     time.Now(),
	}

	stream <- workflow.StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		Data:      map[string]interface{}{"node_id": "end", "status": "completed"},
	}

	return output, nil
}
