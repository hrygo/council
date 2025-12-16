package workflow

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

// Engine orchestrates the workflow execution
type Engine struct {
	Graph         *GraphDefinition
	Status        map[string]NodeStatus
	NodeFactory   func(node *Node) (NodeProcessor, error)
	StreamChannel chan StreamEvent
	mu            sync.RWMutex
	inputs        map[string]interface{}

	// Middleware hooks
	Middlewares []Middleware
	Session     *Session // Reference to the session state
}

// NewEngine creates a new workflow engine
func NewEngine(session *Session) *Engine {
	return &Engine{
		Graph:         session.Graph,
		Status:        make(map[string]NodeStatus),
		StreamChannel: make(chan StreamEvent, 100), // Buffer for safety
		inputs:        session.Inputs,
		Session:       session,
		// Default Factory (can be overridden)
		NodeFactory: func(n *Node) (NodeProcessor, error) {
			return nil, fmt.Errorf("no factory configured for node type %s", n.Type)
		},
	}
}

// Run executes the workflow from the start node
func (e *Engine) Run(ctx context.Context) {
	// 1. Basic Validation
	if err := e.Graph.Validate(); err != nil {
		e.emitError("validation_failed", err)
		return
	}

	// 2. Start Execution
	e.executeNode(ctx, e.Graph.StartNodeID, e.inputs)
}

func (e *Engine) executeNode(ctx context.Context, nodeID string, input map[string]interface{}) {
	e.mu.Lock()
	node, exists := e.Graph.Nodes[nodeID]
	if !exists {
		e.mu.Unlock()
		e.emitError(nodeID, fmt.Errorf("node not found"))
		return
	}
	e.Status[nodeID] = StatusRunning
	e.mu.Unlock()

	// Special Handling for Control Flow Nodes
	if node.Type == NodeTypeParallel {
		e.handleParallel(ctx, node, input)
		return
	}

	// Standard Processing
	processor, err := e.NodeFactory(node)
	if err != nil {
		e.emitError(nodeID, err)
		return
	}

	// Middleware: Before Execution
	for _, mw := range e.Middlewares {
		if err := mw.BeforeNodeExecution(ctx, e.Session, node); err != nil {
			e.emitError(nodeID, fmt.Errorf("middleware %s blocked execution: %w", mw.Name(), err))
			return
		}
	}

	// Execute Processor
	output, err := processor.Process(ctx, input, e.StreamChannel)
	if err != nil {
		e.updateStatus(nodeID, StatusFailed)
		e.emitError(nodeID, err)
		return
	}

	// Middleware: After Execution
	for _, mw := range e.Middlewares { // Execute in order (or reverse? usually order is fine for transform)
		var mwErr error
		output, mwErr = mw.AfterNodeExecution(ctx, e.Session, node, output)
		if mwErr != nil {
			e.emitError(nodeID, fmt.Errorf("middleware %s failed post-processing: %w", mw.Name(), mwErr))
			return
		}
	}

	e.updateStatus(nodeID, StatusCompleted)

	// Propagate to Next Nodes
	// For simple graphs, input -> output flow
	// If multiple next nodes, we execute them concurrently or sequentially?
	// TDD said: "for _, nextID := range node.NextIDs { e.executeNode... }"
	// This implies parallel execution for branching.

	var wg sync.WaitGroup
	for _, nextID := range node.NextIDs {
		wg.Add(1)
		go func(nid string) {
			defer wg.Done()
			e.executeNode(ctx, nid, output)
		}(nextID)
	}
	wg.Wait()
}

func (e *Engine) handleParallel(ctx context.Context, node *Node, input map[string]interface{}) {
	var wg sync.WaitGroup
	for _, nextID := range node.NextIDs {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()
			e.executeNode(ctx, cid, input) // Pass same input to all branches
		}(nextID)
	}
	wg.Wait()
	e.updateStatus(node.ID, StatusCompleted)
}

func (e *Engine) updateStatus(nodeID string, status NodeStatus) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.Status[nodeID] = status
}

func (e *Engine) emitError(nodeID string, err error) {
	log.Printf("Error in node %s: %v", nodeID, err)
	e.StreamChannel <- StreamEvent{
		Type:      "error",
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Data:      map[string]interface{}{"error": err.Error()},
	}
	e.updateStatus(nodeID, StatusFailed)
}
