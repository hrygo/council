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

	// Join mechanism for fan-in nodes (SPEC-1206)
	inDegree      map[string]int                      // Node in-degree count
	pendingInputs map[string][]map[string]interface{} // Pending inputs for join
	joinMu        sync.Mutex                          // Mutex for join operations
	MergeStrategy MergeStrategy                       // Pluggable merge strategy
}

// NewEngine creates a new workflow engine
func NewEngine(session *Session) *Engine {
	e := &Engine{
		Graph:         session.Graph,
		Status:        make(map[string]NodeStatus),
		StreamChannel: make(chan StreamEvent, 100), // Buffer for safety
		inputs:        session.Inputs,
		Session:       session,
		pendingInputs: make(map[string][]map[string]interface{}),
		MergeStrategy: &DefaultMergeStrategy{}, // Default strategy, can be overridden
		// Default Factory (can be overridden)
		NodeFactory: func(n *Node) (NodeProcessor, error) {
			return nil, fmt.Errorf("no factory configured for node type %s", n.Type)
		},
	}
	e.computeInDegrees()
	return e
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
	// Check for Pause
	if e.Session.Status == SessionPaused {
		e.StreamChannel <- StreamEvent{
			Type:      "execution:paused",
			Timestamp: time.Now(),
			// NodeID might be relevant if we pause BEFORE a specific node
			NodeID: nodeID,
			Data:   map[string]interface{}{"reason": "session_paused"},
		}
	}
	if err := e.Session.WaitIfPaused(ctx); err != nil {
		// Context cancelled or other error
		e.emitError(nodeID, err)
		return
	}

	e.mu.Lock()
	node, exists := e.Graph.Nodes[nodeID]
	if !exists {
		e.mu.Unlock()
		e.emitError(nodeID, fmt.Errorf("node not found"))
		return
	}
	e.Status[nodeID] = StatusRunning
	e.mu.Unlock()

	// Notify frontend immediately that node is running
	e.StreamChannel <- StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Data:      map[string]interface{}{"status": "running"},
	}

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

	// Inject Contextual IDs
	if input == nil {
		input = make(map[string]interface{})
	}
	input["session_id"] = e.Session.ID

	// Execute Processor
	output, err := processor.Process(ctx, input, e.StreamChannel)
	if err != nil {
		if err == ErrSuspended {
			e.updateStatus(nodeID, StatusSuspended)
			return // Suspended execution
		}
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

	// Deliver output to downstream nodes using Join mechanism (SPEC-1206)
	e.deliverToDownstream(ctx, nodeID, output)
}

func (e *Engine) handleParallel(ctx context.Context, node *Node, input map[string]interface{}) {
	e.StreamChannel <- StreamEvent{
		Type:      "node:parallel_start",
		Timestamp: time.Now(),
		NodeID:    node.ID,
		Data:      map[string]interface{}{"branches": node.NextIDs},
	}

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
	e.Status[nodeID] = status
	e.mu.Unlock()

	e.StreamChannel <- StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Data:      map[string]interface{}{"status": status},
	}
}

func (e *Engine) GetStatus(nodeID string) NodeStatus {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.Status[nodeID]
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

// ResumeNode resumes execution of a suspended node with provided output
func (e *Engine) ResumeNode(ctx context.Context, nodeID string, output map[string]interface{}) error {
	e.mu.Lock()
	status, exists := e.Status[nodeID]
	e.mu.Unlock()

	if !exists {
		return fmt.Errorf("node %s not found in execution status", nodeID)
	}
	if status != StatusSuspended {
		return fmt.Errorf("node %s is not suspended (status: %s)", nodeID, status)
	}

	// Update Status
	e.updateStatus(nodeID, StatusCompleted)

	e.StreamChannel <- StreamEvent{
		Type:      "node_resumed",
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Data:      output,
	}

	// Use Session's context for long-running workflow execution (Fix-A4)
	// The passed ctx is typically a short-lived HTTP request context
	// Using Session.Context() ensures workflow continues after request ends
	workflowCtx := e.Session.Context()
	if workflowCtx == nil {
		workflowCtx = context.Background()
	}

	// Use deliverToDownstream for consistency with Join mechanism (SPEC-1206)
	go e.deliverToDownstream(workflowCtx, nodeID, output)

	return nil
}

// computeInDegrees calculates the in-degree for each node in the graph.
// This is called during Engine initialization to prepare for join operations.
// Loop nodes' back-edges (first next_id) are excluded from in-degree count
// to prevent deadlock on first execution.
func (e *Engine) computeInDegrees() {
	e.inDegree = make(map[string]int)
	if e.Graph == nil || e.Graph.Nodes == nil {
		return
	}
	for _, node := range e.Graph.Nodes {
		for i, nextID := range node.NextIDs {
			// Skip Loop node's back-edge (first next_id is the continue/loop path)
			// This edge is conditionally triggered and bypasses in-degree check anyway
			if node.Type == NodeTypeLoop && i == 0 {
				continue
			}
			e.inDegree[nextID]++
		}
	}
}

// deliverToDownstream delivers output to downstream nodes.
// For nodes with in-degree > 1, it waits for all upstream nodes to complete
// before executing (join/barrier pattern).
// For Loop nodes, it implements conditional routing based on should_exit.
func (e *Engine) deliverToDownstream(ctx context.Context, nodeID string, output map[string]interface{}) {
	node := e.Graph.Nodes[nodeID]
	if node == nil {
		return
	}

	// Determine which downstream nodes to trigger
	targetNextIDs := node.NextIDs

	// Special handling for Loop nodes: conditional routing (SPEC-1206)
	// Loop nodes have 2 next_ids: [continue_path, exit_path]
	// Route based on should_exit flag in output
	if node.Type == NodeTypeLoop && len(node.NextIDs) >= 2 {
		shouldExit, _ := output["should_exit"].(bool)
		if shouldExit {
			// Exit path: only trigger the second next_id (typically "end")
			targetNextIDs = []string{node.NextIDs[1]}
		} else {
			// Continue path: only trigger the first next_id (loop back)
			targetNextIDs = []string{node.NextIDs[0]}
		}
	}

	// Check if this is a loop-back delivery (from Loop node to continue path)
	// Loop-back deliveries should bypass in-degree waiting to avoid deadlock
	isLoopBack := node.Type == NodeTypeLoop && len(targetNextIDs) == 1 && targetNextIDs[0] == node.NextIDs[0]

	for _, nextID := range targetNextIDs {
		e.joinMu.Lock()

		var mergedInput map[string]interface{}
		var ready bool

		if isLoopBack {
			// Loop-back: bypass in-degree check, directly use output as input
			// This prevents deadlock when looping back to a node that has multiple in-edges
			mergedInput = output
			ready = true
			// Clear any pending inputs for this node to reset state
			e.pendingInputs[nextID] = nil
		} else {
			// Normal path: collect and wait for all upstream nodes
			e.pendingInputs[nextID] = append(e.pendingInputs[nextID], output)

			// Check if all upstream nodes have delivered
			expectedInputs := e.inDegree[nextID]
			receivedInputs := len(e.pendingInputs[nextID])
			ready = receivedInputs >= expectedInputs

			if ready {
				// All upstreams delivered, merge and clear
				mergedInput = e.MergeStrategy.Merge(e.pendingInputs[nextID])
				e.pendingInputs[nextID] = nil
			}
		}

		e.joinMu.Unlock()

		if ready {
			// Execute downstream node with merged input
			e.executeNode(ctx, nextID, mergedInput)
		}
	}
}
