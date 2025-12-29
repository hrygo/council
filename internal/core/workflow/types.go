package workflow

import (
	"context"
	"fmt"
)

// NodeStatus defines the execution state of a node
type NodeStatus string

const (
	StatusPending   NodeStatus = "pending"
	StatusRunning   NodeStatus = "running"
	StatusCompleted NodeStatus = "completed"
	StatusFailed    NodeStatus = "failed"
	StatusSkipped   NodeStatus = "skipped"
	StatusSuspended NodeStatus = "suspended"
)

var (
	// ErrSuspended indicates that execution is suspended at this node
	ErrSuspended = fmt.Errorf("execution suspended")
)

// NodeType enum for supported node types
type NodeType string

const (
	NodeTypeStart           NodeType = "start"
	NodeTypeEnd             NodeType = "end"
	NodeTypeAgent           NodeType = "agent"
	NodeTypeLLM             NodeType = "llm"              // Direct LLM call
	NodeTypeTool            NodeType = "tool"             // Search, etc.
	NodeTypeParallel        NodeType = "parallel"         // Logic node: Parallel branch
	NodeTypeSequence        NodeType = "sequence"         // Logic node: Sequential steps
	NodeTypeVote            NodeType = "vote"             // Logic node: Vote
	NodeTypeLoop            NodeType = "loop"             // Logic node: Loop
	NodeTypeFactCheck       NodeType = "fact_check"       // Logic node: Fact Check
	NodeTypeHumanReview     NodeType = "human_review"     // Logic node: Human Review
	NodeTypeMemoryRetrieval NodeType = "memory_retrieval" // Logic node: Memory Retrieval
	NodeTypeContextSynth    NodeType = "context_synth"    // Logic node: Context Synthesizer
)

// GraphDefinition represents the static definition of a workflow
type GraphDefinition struct {
	ID          string           `json:"workflow_uuid"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Nodes       map[string]*Node `json:"nodes"`
	StartNodeID string           `json:"start_node_id"`
}

// Node represents a single step in the workflow
type Node struct {
	ID         string                 `json:"node_id"`
	Type       NodeType               `json:"type"`
	Name       string                 `json:"name"`
	NextIDs    []string               `json:"next_ids,omitempty"` // Adjacency list for next steps
	Properties map[string]interface{} `json:"properties"`         // Node-specific config (e.g. Prompt, Model)
}

// Middleware allows intercepting node execution for safety and observability
type Middleware interface {
	Name() string
	BeforeNodeExecution(ctx context.Context, session *Session, node *Node) error
	AfterNodeExecution(ctx context.Context, session *Session, node *Node, output map[string]interface{}) (map[string]interface{}, error)
}
