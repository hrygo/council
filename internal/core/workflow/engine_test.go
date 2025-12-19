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
	session := NewSession(graph, nil)
	engine := NewEngine(session)

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

func TestEngine_ResumeNode(t *testing.T) {
	graph := &GraphDefinition{
		ID:          "test-graph",
		StartNodeID: "node1",
		Nodes: map[string]*Node{
			"node1": {ID: "node1", Type: "suspending_node", NextIDs: []string{"node2"}},
			"node2": {ID: "node2", Type: "test_node"},
		},
	}

	session := NewSession(graph, nil)
	engine := NewEngine(session)

	engine.NodeFactory = func(n *Node) (NodeProcessor, error) {
		if n.Type == "suspending_node" {
			return &suspendingProcessor{}, nil
		}
		return &MockProcessor{Output: map[string]interface{}{"val": n.ID}}, nil
	}

	// 1. Run until suspended
	ctx := context.Background()
	session.Start(ctx)
	engine.Run(ctx)

	if engine.GetStatus("node1") != StatusSuspended {
		t.Errorf("expected node1 to be suspended, got %s", engine.GetStatus("node1"))
	}

	// 2. Resume
	err := engine.ResumeNode(ctx, "node1", map[string]interface{}{"resumed": true})
	if err != nil {
		t.Errorf("ResumeNode failed: %v", err)
	}

	// Wait for node2 (ResumeNode runs in background)
	time.Sleep(50 * time.Millisecond)

	if engine.GetStatus("node1") != StatusCompleted {
		t.Errorf("expected node1 to be completed, got %s", engine.GetStatus("node1"))
	}
	if engine.GetStatus("node2") != StatusCompleted {
		t.Errorf("expected node2 to be completed, got %s", engine.GetStatus("node2"))
	}
}

type suspendingProcessor struct{}

func (p *suspendingProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
	return nil, ErrSuspended
}

func TestEngine_Parallel(t *testing.T) {
	graph := &GraphDefinition{
		ID:          "parallel-graph",
		StartNodeID: "p1",
		Nodes: map[string]*Node{
			"p1": {ID: "p1", Type: NodeTypeParallel, NextIDs: []string{"b1", "b2"}},
			"b1": {ID: "b1", Type: "test"},
			"b2": {ID: "b2", Type: "test"},
		},
	}

	session := NewSession(graph, nil)
	engine := NewEngine(session)

	executed := make(map[string]int)
	mu := sync.Mutex{}
	engine.NodeFactory = func(n *Node) (NodeProcessor, error) {
		mu.Lock()
		executed[n.ID]++
		mu.Unlock()
		return &MockProcessor{}, nil
	}

	session.Start(context.Background())
	engine.Run(context.Background())

	if executed["b1"] != 1 || executed["b2"] != 1 {
		t.Errorf("expected branches to execute exactly once, got b1:%d, b2:%d", executed["b1"], executed["b2"])
	}
}
