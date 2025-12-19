package workflow

import (
	"context"
	"testing"
	"time"
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
func TestSession_PauseResume(t *testing.T) {
	session := NewSession(&GraphDefinition{}, nil)
	session.Start(context.Background())

	// Pause
	session.Pause()
	if session.Status != SessionPaused {
		t.Errorf("expected Paused, got %s", session.Status)
	}

	// Verify WaitIfPaused blocks (using a timer)
	done := make(chan bool)
	go func() {
		err := session.WaitIfPaused(context.Background())
		if err != nil {
			t.Errorf("WaitIfPaused returned error: %v", err)
		}
		done <- true
	}()

	select {
	case <-done:
		t.Fatal("WaitIfPaused did not block while paused")
	case <-time.After(50 * time.Millisecond):
		// OK, it blocked
	}

	// Resume
	session.Resume()
	if session.Status != SessionRunning {
		t.Errorf("expected Running, got %s", session.Status)
	}

	select {
	case <-done:
		// OK, it unblocked
	case <-time.After(50 * time.Millisecond):
		t.Fatal("WaitIfPaused did not unblock after Resume")
	}
}

func TestSession_Signals(t *testing.T) {
	session := NewSession(&GraphDefinition{}, nil)
	session.Start(context.Background())

	ch := session.GetSignalChannel("node-1")

	go func() {
		err := session.SendSignal("node-1", "data")
		if err != nil {
			t.Errorf("SendSignal failed: %v", err)
		}
	}()

	select {
	case val := <-ch:
		if val != "data" {
			t.Errorf("expected data, got %v", val)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Did not receive signal")
	}

	// Test Error cases
	err := session.SendSignal("non-existent", nil)
	if err == nil {
		t.Errorf("expected error for non-existent signal channel")
	}
}
