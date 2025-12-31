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
		"ignored":  "should not be in output",
	}

	processor := &StartProcessor{
		OutputKeys: []string{"proposal"},
	}
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
	if _, ok := output["ignored"]; ok {
		t.Errorf("Expected 'ignored' key to be filtered out")
	}

	// Verify Stream Events - Lifecycle events moved to Engine
	close(stream)
	for range stream {
		// drain
	}
}
