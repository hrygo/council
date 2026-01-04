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
	Graph  *GraphDefinition
	Status map[string]NodeStatus
	// NodeFactory creates processors for graph nodes
	NodeFactory   NodeFactory
	StreamChannel chan StreamEvent
	// Public Mutex for state access
	Mu     sync.RWMutex
	inputs map[string]interface{}

	// Middleware hooks
	Middlewares []Middleware
	Session     *Session // Reference to the session state

	// Join mechanism for fan-in nodes (SPEC-1206)
	inDegree      map[string]int                      // Node in-degree count
	pendingInputs map[string][]map[string]interface{} // Pending inputs for join
	joinMu        sync.Mutex                          // Mutex for join operations
	MergeStrategy MergeStrategy                       // Pluggable merge strategy
	SessionRepo   SessionRepository                   // Injected persistence
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
		NodeFactory:   &DefaultNodeFactory{},   // Default Factory
	}
	// Resume status from session if available
	if session.NodeStatuses != nil {
		for k, v := range session.NodeStatuses {
			e.Status[k] = v
		}
	}
	e.computeInDegrees()
	return e
}

func (e *Engine) SetSessionRepository(repo SessionRepository) {
	e.SessionRepo = repo
}

// Run executes the workflow from the start node
func (e *Engine) Run(ctx context.Context) error {
	// 1. Basic Validation
	if err := e.Graph.Validate(); err != nil {
		e.emitError("validation_failed", err)
		return err
	}

	startNodeID := e.Graph.StartNodeID
	if startNodeID == "" {
		return fmt.Errorf("no start node defined")
	}

	// Start from the defined start node
	// If resuming (some nodes completed), we might need logic?
	// Run() usually implies fresh start. Resume() is manual?
	// If we have statuses, we don't re-run completed nodes?
	// Current Run() just starts at StartNode.
	// Assuming idempotent execution or fresh start.
	return e.executeNode(ctx, startNodeID, e.inputs)
}

// executeNodeWithoutDownstream executes a single node without automatically triggering downstream nodes
// Used for parallel branches where downstream execution is coordinated by the parallel node itself
func (e *Engine) executeNodeWithoutDownstream(ctx context.Context, nodeID string, input map[string]interface{}) (map[string]interface{}, error) {
	// Check for Pause
	if e.Session.Status == SessionPaused {
		e.StreamChannel <- StreamEvent{
			Type:      "execution:paused",
			Timestamp: time.Now(),
			NodeID:    nodeID,
			Data:      map[string]interface{}{"reason": "session_paused"},
		}
	}
	if err := e.Session.WaitIfPaused(ctx); err != nil {
		e.emitError(nodeID, err)
		return nil, err
	}

	e.Mu.RLock()
	node, exists := e.Graph.Nodes[nodeID]
	e.Mu.RUnlock()

	if !exists {
		e.emitError(nodeID, fmt.Errorf("node not found"))
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	// Parallel nodes should not be executed via this function
	if node.Type == NodeTypeParallel {
		return nil, fmt.Errorf("parallel nodes should be handled by handleParallel")
	}

	// Update status
	e.updateStatus(nodeID, StatusRunning)

	// Middleware: Before
	for _, mw := range e.Middlewares {
		if err := mw.BeforeNodeExecution(ctx, e.Session, node); err != nil {
			e.emitError(nodeID, fmt.Errorf("middleware %s blocked execution: %w", mw.Name(), err))
			return nil, err
		}
	}

	// Inject Contextual IDs
	if input == nil {
		input = make(map[string]interface{})
	}
	input["session_id"] = e.Session.ID

	// Standard Processing using Factory
	processor, err := e.NodeFactory.CreateNode(node, FactoryDeps{Session: e.Session})
	if err != nil {
		e.emitError(nodeID, err)
		e.updateStatus(nodeID, StatusFailed)
		return nil, err
	}

	output, err := processor.Process(ctx, input, e.StreamChannel)
	if err != nil {
		if err == ErrSuspended {
			e.updateStatus(nodeID, StatusSuspended)
			return nil, nil // Suspended execution
		}
		e.updateStatus(nodeID, StatusFailed)
		e.emitError(nodeID, err)
		return nil, err
	}

	// Middleware: After Execution
	for _, mw := range e.Middlewares {
		var mwErr error
		output, mwErr = mw.AfterNodeExecution(ctx, e.Session, node, output)
		if mwErr != nil {
			e.emitError(nodeID, fmt.Errorf("middleware %s failed post-processing: %w", mw.Name(), mwErr))
			return nil, mwErr
		}
	}

	e.updateStatus(nodeID, StatusCompleted)

	// NOTE: Unlike executeNode, we DON'T call deliverToDownstream here
	// The caller is responsible for handling downstream execution
	return output, nil
}

// executeNode runs a single node and recursively its children
func (e *Engine) executeNode(ctx context.Context, nodeID string, input map[string]interface{}) error {
	// e.Mu.RLock()
	// status, known := e.Status[nodeID]
	// e.Mu.RUnlock()

	// If node is already completed, skip execution (Idempotency)
	// if known && status == StatusCompleted {
	// 	// Logic to propel downstream?
	// 	// ...
	// }

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
		return err
	}

	e.Mu.RLock()
	node, exists := e.Graph.Nodes[nodeID]
	e.Mu.RUnlock()

	if !exists {
		e.emitError(nodeID, fmt.Errorf("node not found"))
		return fmt.Errorf("node %s not found", nodeID)
	}

	// Update status
	e.updateStatus(nodeID, StatusRunning)

	// Middleware: Before
	for _, mw := range e.Middlewares {
		if err := mw.BeforeNodeExecution(ctx, e.Session, node); err != nil {
			e.emitError(nodeID, fmt.Errorf("middleware %s blocked execution: %w", mw.Name(), err))
			return err
		}
	}

	// Inject Contextual IDs
	if input == nil {
		input = make(map[string]interface{})
	}
	input["session_id"] = e.Session.ID

	var output map[string]interface{}
	var err error

	// Special Handling for Control Flow Nodes
	if node.Type == NodeTypeParallel {
		e.handleParallel(ctx, node, input)
		return nil
	}

	// Standard Processing using Factory
	processor, err := e.NodeFactory.CreateNode(node, FactoryDeps{Session: e.Session})
	if err != nil {
		e.emitError(nodeID, err)
		e.updateStatus(nodeID, StatusFailed)
		return err
	}

	output, err = processor.Process(ctx, input, e.StreamChannel)
	if err != nil {
		if err == ErrSuspended {
			e.updateStatus(nodeID, StatusSuspended)
			return nil // Suspended execution
		}
		e.updateStatus(nodeID, StatusFailed)
		e.emitError(nodeID, err)
		return err
	}

	// Middleware: After Execution
	for _, mw := range e.Middlewares {
		var mwErr error
		output, mwErr = mw.AfterNodeExecution(ctx, e.Session, node, output)
		if mwErr != nil {
			e.emitError(nodeID, fmt.Errorf("middleware %s failed post-processing: %w", mw.Name(), mwErr))
			return mwErr
		}
	}

	e.updateStatus(nodeID, StatusCompleted)

	// Determine Routing
	var nextIDs []string
	if router, ok := processor.(ConditionalRouter); ok {
		nextIDs, err = router.GetNextNodes(ctx, output, node.NextIDs)
		if err != nil {
			e.emitError(nodeID, fmt.Errorf("routing failed: %w", err))
			return err
		}
	} else {
		nextIDs = node.NextIDs
	}

	// Deliver output to downstream nodes using Join mechanism (SPEC-1206)
	e.deliverToDownstream(ctx, nodeID, output, nextIDs)
	return nil
}

func (e *Engine) handleParallel(ctx context.Context, node *Node, input map[string]interface{}) {
	// Set parallel node status to running
	log.Printf("[Engine] Setting parallel node %s to RUNNING", node.ID)
	e.updateStatus(node.ID, StatusRunning)

	e.StreamChannel <- StreamEvent{
		Type:      "node:parallel_start",
		Timestamp: time.Now(),
		NodeID:    node.ID,
		Data:      map[string]interface{}{"branches": node.NextIDs},
	}

	// Collect outputs from all branches for potential Join logic
	type branchResult struct {
		nodeID string
		output map[string]interface{}
		err    error
	}
	resultsChan := make(chan branchResult, len(node.NextIDs))

	var wg sync.WaitGroup
	for _, nextID := range node.NextIDs {
		wg.Add(1)
		go func(cid string) {
			defer wg.Done()

			// Clone input map to avoid Data Race
			clonedInput := make(map[string]interface{}, len(input))
			for k, v := range input {
				clonedInput[k] = v
			}

			// Execute branch node WITHOUT automatically continuing to downstream
			// We'll handle downstream execution after all branches complete
			output, err := e.executeNodeWithoutDownstream(ctx, cid, clonedInput)
			resultsChan <- branchResult{nodeID: cid, output: output, err: err}

			if err != nil {
				log.Printf("[Engine] Parallel branch %s failed: %v", cid, err)
			}
		}(nextID)
	}

	log.Printf("[Engine] Waiting for parallel branches of %s to complete...", node.ID)
	wg.Wait()
	close(resultsChan)

	log.Printf("[Engine] All parallel branches of %s completed. Setting to COMPLETED", node.ID)
	e.updateStatus(node.ID, StatusCompleted)

	// Now deliver to downstream nodes (Join point)
	var successfulBranches []branchResult
	for result := range resultsChan {
		if result.err == nil && result.output != nil {
			successfulBranches = append(successfulBranches, result)
		}
	}

	// Now deliver to downstream nodes (Deferred Execution)
	// We iterate through all branch results and deliver them individually.
	// This ensures that:
	// 1. The Parallel Node is already COMPLETED (User Requirement)
	// 2. The downstream Join Node receives the correct number of inputs to satisfy InDegree (System Requirement)

	// Prepare delivery tasks to avoid holding RLock during downstream execution (Deadlock prevention)
	type deliveryTask struct {
		nodeID  string
		output  map[string]interface{}
		nextIDs []string
	}
	var tasks []deliveryTask

	e.Mu.RLock()
	for _, res := range successfulBranches {
		if branch, ok := e.Graph.Nodes[res.nodeID]; ok {
			tasks = append(tasks, deliveryTask{
				nodeID:  res.nodeID,
				output:  res.output,
				nextIDs: branch.NextIDs,
			})
		}
	}
	e.Mu.RUnlock()

	// Execute deliveries
	for _, task := range tasks {
		e.deliverToDownstream(ctx, task.nodeID, task.output, task.nextIDs)
	}
}

func (e *Engine) updateStatus(nodeID string, status NodeStatus) {
	e.Mu.Lock()
	e.Status[nodeID] = status
	e.Mu.Unlock()

	log.Printf("[Engine] updateStatus: %s -> %s", nodeID, status)

	// Persist status if repo is injected
	if e.SessionRepo != nil {
		go func() {
			if err := e.SessionRepo.UpdateNodeStatus(context.Background(), e.Session.ID, nodeID, status); err != nil {
				log.Printf("Failed to persist status for node %s: %v", nodeID, err)
			}
		}()
	}

	event := StreamEvent{
		Type:      "node_state_change",
		Timestamp: time.Now(),
		NodeID:    nodeID,
		Data:      map[string]interface{}{"status": status},
	}
	log.Printf("[Engine] Sending WebSocket event: type=%s, node_id=%s, status=%s", event.Type, nodeID, status)
	e.StreamChannel <- event
}

func (e *Engine) GetStatus(nodeID string) NodeStatus {
	e.Mu.RLock()
	defer e.Mu.RUnlock()
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
	e.Mu.Lock()
	status, exists := e.Status[nodeID]
	node, nodeExists := e.Graph.Nodes[nodeID]
	e.Mu.Unlock()

	if !exists {
		return fmt.Errorf("node %s not found in execution status", nodeID)
	}
	if !nodeExists {
		return fmt.Errorf("node %s not found in graph", nodeID)
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
	workflowCtx := e.Session.Context()
	if workflowCtx == nil {
		workflowCtx = context.Background()
	}

	// Determine Routing logic (SPEC-1304)
	var nextIDs []string
	processor, err := e.NodeFactory.CreateNode(node, FactoryDeps{Session: e.Session})
	if err != nil {
		// If factory failed, we fallback to default nextIDs? Or return error?
		// Suspended node implies it WAS created before.
		// Fallback to default check:
		nextIDs = node.NextIDs
	} else {
		if router, ok := processor.(ConditionalRouter); ok {
			var errRoute error
			nextIDs, errRoute = router.GetNextNodes(workflowCtx, output, node.NextIDs)
			if errRoute != nil {
				e.emitError(nodeID, fmt.Errorf("routing failed on resume: %w", errRoute))
				return errRoute
			}
		} else {
			nextIDs = node.NextIDs
		}
	}

	// Use deliverToDownstream for consistency with Join mechanism (SPEC-1206)
	go e.deliverToDownstream(workflowCtx, nodeID, output, nextIDs)

	return nil
}

// computeInDegrees calculates the in-degree for each node in the graph.
func (e *Engine) computeInDegrees() {
	e.inDegree = make(map[string]int)
	if e.Graph == nil || e.Graph.Nodes == nil {
		return
	}
	for _, node := range e.Graph.Nodes {
		for i, nextID := range node.NextIDs {
			// Skip Loop node's back-edge (first next_id is the continue/loop path)
			if node.Type == NodeTypeLoop && i == 0 {
				continue
			}
			e.inDegree[nextID]++
		}
	}
}

// deliverToDownstream delivers output to downstream nodes.
func (e *Engine) deliverToDownstream(ctx context.Context, nodeID string, output map[string]interface{}, targetNextIDs []string) {
	node := e.Graph.Nodes[nodeID]
	if node == nil {
		return
	}

	// Logic routed via targetNextIDs passed from executeNode/ResumeNode

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
			if err := e.executeNode(ctx, nextID, mergedInput); err != nil {
				log.Printf("[Engine] Downstream node %s failed: %v", nextID, err)
			}
		}
	}
}
