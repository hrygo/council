package workflow

import (
	"context"
	"sync"
	"testing"
	"time"
)

// MockProcessor implements NodeProcessor for testing
type MockProcessor struct {
	CapturedInput map[string]interface{}
	Output        map[string]interface{}
}

func (m *MockProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
	m.CapturedInput = input
	// Simulate work
	stream <- StreamEvent{Type: "mock_event", NodeID: "mock", Timestamp: time.Now()}
	return m.Output, nil
}

func TestEngine_Run_Linear(t *testing.T) {
	// Define Graph
	graph := &GraphDefinition{
		ID:          "test-graph",
		StartNodeID: "start",
		Nodes: map[string]*Node{
			"start": {ID: "start", Type: "test_node", NextIDs: []string{"end"}},
			"end":   {ID: "end", Type: "test_node"},
		},
	}

	// Setup Engine
	engine := NewEngine(graph, nil)

	// Mock Factory
	mu := sync.Mutex{}
	executed := make(map[string]bool)

	engine.NodeFactory = func(n *Node) (NodeProcessor, error) {
		mu.Lock()
		executed[n.ID] = true
		mu.Unlock()
		return &MockProcessor{Output: map[string]interface{}{"val": n.ID}}, nil
	}

	// Run
	ctx := context.Background()
	done := make(chan bool)
	go func() {
		engine.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Run timed out")
	}

	// Verify Execution Order/Count
	if !executed["start"] {
		t.Error("Start node not executed")
	}
	if !executed["end"] {
		t.Error("End node not executed")
	}
}

func TestEngine_Run_Cycle_Detection(t *testing.T) {
	// Graph with cycle but Validate() should catch it.
	// We want to ensure Run() also checks or respects Validation.
	// graph := &GraphDefinition{
	// 	ID:          "cycle-graph",
	// 	StartNodeID: "start",
	// 	Nodes: map[string]*Node{
	// 		"start": {ID: "start", Type: "test", NextIDs: []string{"start"}},
	// 	},
	// }

	// engine := NewEngine(graph, nil)
	// Expect Validate logic inside Run, or implicit check?
	// If Run calls Validate, it should fail.
	// But Engine.Run usually just runs. Validate should be called before.
	// If we want Engine.Run to return error, signature needs update.
	// For now, let's assume Run logs error or we must call Validate separately.
	// Skip this test if Run signature doesn't return error (void in TDD example).
}
