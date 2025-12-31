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
	NodeID         string // Dynamic Node ID
	LLM            llm.LLMProvider
	Model          string
	Prompt         string
	PromptSections []workflow.PromptSection // Configuration
	OutputKey      string                   // Configuration: Key for summary (e.g. "final_report")
}

func (e *EndProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- workflow.StreamEvent) (map[string]interface{}, error) {
	// 1. Logic start
	_ = stream // for now use stream for tokens below

	// 2. Aggregate Content - Build structured context for report generation
	var contentBuilder strings.Builder

	for _, section := range e.PromptSections {
		if val, ok := input[section.Key]; ok {
			strVal := fmt.Sprintf("%v", val)
			if strVal != "" {
				contentBuilder.WriteString("## " + section.Label + "\n")
				contentBuilder.WriteString(strVal)
				contentBuilder.WriteString("\n\n")
			}
		}
	}

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
						Data:      map[string]interface{}{"node_id": e.NodeID, "chunk": chunk.Content},
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
	outputKey := e.OutputKey
	if outputKey == "" {
		outputKey = "final_report"
	}
	output := map[string]interface{}{
		outputKey:  finalSummary.String(),
		"ended_at": time.Now(),
	}

	_ = 0 // end logic complete

	return output, nil
}
