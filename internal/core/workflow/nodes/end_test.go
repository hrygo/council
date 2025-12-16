package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/llm"
)

func TestEndProcessor_Process(t *testing.T) {
	// Setup Mock
	mockLLM := llm.NewMockProvider()
	mockLLM.StreamContent = []string{"Summary", " ", "Result"}

	processor := &EndProcessor{
		LLM:   mockLLM,
		Model: "gpt-4",
	}

	input := map[string]interface{}{
		"proposal": "Some long proposal text",
	}
	stream := make(chan workflow.StreamEvent, 100)

	ctx := context.Background()
	output, err := processor.Process(ctx, input, stream)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify Output
	if report, ok := output["final_report"].(string); !ok || report != "Summary Result" {
		t.Errorf("Expected 'Summary Result', got '%v'", report)
	}

	// Verify Stream
	close(stream)
	var tokens []string
	for e := range stream {
		if e.Type == "token_stream" {
			if chunk, ok := e.Data["chunk"].(string); ok {
				tokens = append(tokens, chunk)
			}
		}
	}

	if len(tokens) != 3 {
		t.Errorf("Expected 3 token events, got %d", len(tokens))
	}
}
