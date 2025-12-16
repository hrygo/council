package workflow

import (
	"context"
	"testing"
)

func TestSession_Transitions(t *testing.T) {
	graph := &GraphDefinition{ID: "test-graph"}
	session := NewSession(graph, nil)

	// Initial State
	if session.Status != SessionPending {
		t.Errorf("expected Pending, got %s", session.Status)
	}

	// Start
	ctx := context.Background()
	session.Start(ctx)
	if session.Status != SessionRunning {
		t.Errorf("expected Running, got %s", session.Status)
	}
	if session.StartTime.IsZero() {
		t.Error("expected StartTime to be set")
	}

	// Complete (Simulate)
	session.Complete()
	if session.Status != SessionCompleted {
		t.Errorf("expected Completed, got %s", session.Status)
	}
	if session.EndTime.IsZero() {
		t.Error("expected EndTime to be set")
	}
}

func TestSession_Stop(t *testing.T) {
	session := NewSession(&GraphDefinition{}, nil)
	session.Start(context.Background())

	session.Stop()
	if session.Status != SessionFailed { // Or Stopped/Cancelled depending on design
		t.Errorf("expected Failed (or Stopped), got %s", session.Status)
	}
}
