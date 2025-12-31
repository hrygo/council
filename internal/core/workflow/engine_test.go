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

// SimpleFuncNodeFactory adapts a simple function to the NodeFactory interface
type SimpleFuncNodeFactory func(node *Node) (NodeProcessor, error)

func (f SimpleFuncNodeFactory) CreateNode(node *Node, deps FactoryDeps) (NodeProcessor, error) {
	return f(node)
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

	engine.NodeFactory = SimpleFuncNodeFactory(func(n *Node) (NodeProcessor, error) {
		mu.Lock()
		executed[n.ID] = true
		mu.Unlock()
		return &MockProcessor{Output: map[string]interface{}{"val": n.ID}}, nil
	})

	// Run
	ctx := context.Background()
	done := make(chan bool)
	go func() {
		_ = engine.Run(ctx)
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

	engine.NodeFactory = SimpleFuncNodeFactory(func(n *Node) (NodeProcessor, error) {
		if n.Type == "suspending_node" {
			return &suspendingProcessor{}, nil
		}
		return &MockProcessor{Output: map[string]interface{}{"val": n.ID}}, nil
	})

	// 1. Run until suspended
	ctx := context.Background()
	session.Start(ctx)
	_ = engine.Run(ctx)

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
	engine.NodeFactory = SimpleFuncNodeFactory(func(n *Node) (NodeProcessor, error) {
		mu.Lock()
		executed[n.ID]++
		mu.Unlock()
		return &MockProcessor{}, nil
	})

	session.Start(context.Background())
	_ = engine.Run(context.Background())

	if executed["b1"] != 1 || executed["b2"] != 1 {
		t.Errorf("expected branches to execute exactly once, got b1:%d, b2:%d", executed["b1"], executed["b2"])
	}
}

// TestEngine_JoinMechanism tests the fan-in/join behavior (SPEC-1206)
// When multiple upstream nodes point to the same downstream node,
// the downstream should execute only once with merged inputs.
func TestEngine_JoinMechanism(t *testing.T) {
	// Graph: start -> parallel -> [branch1, branch2] -> join_node -> end
	// join_node has in-degree=2, should wait for both branches
	graph := &GraphDefinition{
		ID:          "join-test-graph",
		StartNodeID: "start",
		Nodes: map[string]*Node{
			"start":     {ID: "start", Type: NodeTypeStart, NextIDs: []string{"parallel"}},
			"parallel":  {ID: "parallel", Type: NodeTypeParallel, NextIDs: []string{"branch1", "branch2"}},
			"branch1":   {ID: "branch1", Type: "test", NextIDs: []string{"join_node"}},
			"branch2":   {ID: "branch2", Type: "test", NextIDs: []string{"join_node"}},
			"join_node": {ID: "join_node", Type: "test", NextIDs: []string{"end"}},
			"end":       {ID: "end", Type: "test"},
		},
	}

	session := NewSession(graph, nil)
	engine := NewEngine(session)

	executed := make(map[string]int)
	receivedInputs := make(map[string][]map[string]interface{})
	mu := sync.Mutex{}

	engine.NodeFactory = SimpleFuncNodeFactory(func(n *Node) (NodeProcessor, error) {
		return &MockProcessor{
			Output: map[string]interface{}{
				"source":       n.ID,
				"agent_output": "Output from " + n.ID,
			},
		}, nil
	})

	// Override with a tracking processor for join_node
	originalFactory := engine.NodeFactory
	engine.NodeFactory = SimpleFuncNodeFactory(func(n *Node) (NodeProcessor, error) {
		mu.Lock()
		executed[n.ID]++
		mu.Unlock()

		if n.ID == "join_node" {
			return &InputCapturingProcessor{
				OnProcess: func(input map[string]interface{}) {
					mu.Lock()
					receivedInputs[n.ID] = append(receivedInputs[n.ID], input)
					mu.Unlock()
				},
			}, nil
		}
		return originalFactory.CreateNode(n, FactoryDeps{})
	})

	session.Start(context.Background())
	_ = engine.Run(context.Background())

	// Verify join_node executed exactly once
	if executed["join_node"] != 1 {
		t.Errorf("expected join_node to execute exactly once, got %d", executed["join_node"])
	}

	// Verify join_node received merged input with branch data
	if len(receivedInputs["join_node"]) != 1 {
		t.Errorf("expected join_node to receive 1 merged input, got %d", len(receivedInputs["join_node"]))
	}

	if len(receivedInputs["join_node"]) > 0 {
		mergedInput := receivedInputs["join_node"][0]
		// Should have branch_0 and branch_1 from merge
		if _, hasBranch0 := mergedInput["branch_0"]; !hasBranch0 {
			t.Error("expected merged input to have branch_0")
		}
		if _, hasBranch1 := mergedInput["branch_1"]; !hasBranch1 {
			t.Error("expected merged input to have branch_1")
		}
	}
}

// InputCapturingProcessor captures input for testing
type InputCapturingProcessor struct {
	OnProcess func(input map[string]interface{})
}

func (p *InputCapturingProcessor) Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (map[string]interface{}, error) {
	if p.OnProcess != nil {
		p.OnProcess(input)
	}
	stream <- StreamEvent{Type: "mock_event", NodeID: "capture", Timestamp: time.Now()}
	return map[string]interface{}{"captured": true}, nil
}
