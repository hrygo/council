package workflow

import (
	"context"
	"time"
)

// StreamEvent represents a real-time event sent to the client
type StreamEvent struct {
	Type      string                 `json:"type"` // e.g. "node_start", "token", "node_end", "error"
	Timestamp time.Time              `json:"timestamp"`
	NodeID    string                 `json:"node_id,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
}

// WorkflowContext handles state and data flow during execution
type WorkflowContext struct {
	ctx       context.Context
	SessionID string
	inputs    map[string]interface{} // Global read-only inputs
}

// NewWorkflowContext creates a base context wrapper
func NewWorkflowContext(ctx context.Context, sessionID string, inputs map[string]interface{}) *WorkflowContext {
	return &WorkflowContext{
		ctx:       ctx,
		SessionID: sessionID,
		inputs:    inputs,
	}
}

// Context returns the underlying standard context
func (wc *WorkflowContext) Context() context.Context {
	return wc.ctx
}
