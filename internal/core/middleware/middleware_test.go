package middleware

import (
	"context"
	"testing"

	"github.com/hrygo/council/internal/core/workflow"
	"github.com/hrygo/council/internal/infrastructure/mocks"
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

func TestMemoryMiddleware(t *testing.T) {
	mockManager := &mocks.MemoryMockManager{}
	mw := NewMemoryMiddleware(mockManager)

	session := &workflow.Session{
		ID:     "s1",
		Inputs: map[string]interface{}{"group_id": "g1"},
	}
	node := &workflow.Node{ID: "n1"}
	output := map[string]interface{}{
		"content": "Secret message",
	}

	_, err := mw.AfterNodeExecution(context.Background(), session, node, output)
	if err != nil {
		t.Fatal(err)
	}

	if len(mockManager.CapturedQuarantine) != 1 || mockManager.CapturedQuarantine[0] != "Secret message" {
		t.Errorf("Expected quarantine log, got %v", mockManager.CapturedQuarantine)
	}
}
