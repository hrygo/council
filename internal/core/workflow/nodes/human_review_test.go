package nodes

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
)

func TestHumanReviewProcessor_Process(t *testing.T) {
	p := &HumanReviewProcessor{
		TimeoutMinutes: 30,
	}

	stream := make(chan workflow.StreamEvent, 10)
	input := map[string]interface{}{}

	output, err := p.Process(context.Background(), input, stream)

	if err != workflow.ErrSuspended {
		t.Errorf("expected ErrSuspended, got %v", err)
	}
	if output != nil {
		t.Error("expected nil output")
	}

	// Check stream events
	foundRequired := false
	close(stream)
	for event := range stream {
		if event.Type == "human_interaction_required" {
			foundRequired = true
			if event.Data["timeout"] != 30 {
				t.Errorf("expected timeout 30, got %v", event.Data["timeout"])
			}
		}
	}

	if !foundRequired {
		t.Error("human_interaction_required event not found in stream")
	}
}
