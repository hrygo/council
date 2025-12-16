package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
)

func TestStartProcessor_Process(t *testing.T) {
	// Mock inputs
	input := map[string]interface{}{
		"proposal": "Let's build a spaceship",
		"attachments": []map[string]interface{}{
			{
				"file_name":    "specs.md",
				"content_type": "text/markdown",
				"content":      "# Specs\n\n1. Engine", // Simulating already read content for simplicity or mock file path logic later
			},
		},
	}

	processor := &StartProcessor{}
	stream := make(chan workflow.StreamEvent, 10)

	ctx := context.Background()
	output, err := processor.Process(ctx, input, stream)
	if err != nil {
		t.Fatalf("Process failed: %v", err)
	}

	// Verify Output
	if output["proposal"] != "Let's build a spaceship" {
		t.Errorf("Expected proposal in output")
	}
	if ctxStr, ok := output["combined_context"].(string); !ok || ctxStr == "" {
		t.Errorf("Expected combined_context in output")
	}

	// Verify Stream Events
	close(stream)
	var events []workflow.StreamEvent
	for e := range stream {
		events = append(events, e)
	}

	if len(events) < 2 {
		t.Errorf("Expected at least 2 events (start/complete), got %d", len(events))
	}
	if events[0].Type != "node_state_change" || events[0].Data["status"] != "running" {
		t.Errorf("Expected start running event")
	}
}
