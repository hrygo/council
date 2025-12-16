package middleware

import (
	"context"
	"testing"
)

func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(10)
	if cb.Name() != "LogicCircuitBreaker" {
		t.Errorf("Expected LogicCircuitBreaker, got %s", cb.Name())
	}

	// Test execution flow
	err := cb.BeforeNodeExecution(context.Background(), nil, nil)
	if err != nil {
		t.Errorf("Expected nil error for simple check, got %v", err)
	}
}

func TestFactCheck(t *testing.T) {
	fc := NewFactCheckTrigger()
	output := map[string]interface{}{
		"content": "This is a fact [Specific Metric] check.",
	}

	newOutput, err := fc.AfterNodeExecution(context.Background(), nil, nil, output)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	meta := newOutput["metadata"].(map[string]interface{})
	if !meta["verify_pending"].(bool) {
		t.Error("Expected verify_pending flag")
	}
}
