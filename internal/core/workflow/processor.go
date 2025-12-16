package workflow

import (
	"context"
)

// NodeProcessor defines the interface that all node types must implement
type NodeProcessor interface {
	// Process executes the node's logic.
	// input: Output from previous nodes (or initial input for StartNode)
	// stream: Channel to push real-time events
	// It returns the output map for downstream nodes.
	Process(ctx context.Context, input map[string]interface{}, stream chan<- StreamEvent) (output map[string]interface{}, err error)
}
