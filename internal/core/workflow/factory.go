package workflow

import (
	"context"
	"fmt"
)

// ConditionalRouter extends NodeProcessor for dynamic routing.
// Nodes that implement this interface can decide which next nodes to execute
// based on their execution output.
type ConditionalRouter interface {
	GetNextNodes(ctx context.Context, output map[string]interface{}, defaultNextIDs []string) ([]string, error)
}

// FactoryDeps holds dependencies injected into nodes during creation.
type FactoryDeps struct {
	Session *Session
}

// NodeFactory defines the interface for creating node processors.
// Applications should implement this to provide specific node implementations.
type NodeFactory interface {
	CreateNode(node *Node, deps FactoryDeps) (NodeProcessor, error)
}

// DefaultNodeFactory is a base implementation of NodeFactory.
// It can be embedded in application-specific factories or used for testing.
type DefaultNodeFactory struct{}

// CreateNode returns an error for all node types, as the default factory has no built-in processors.
func (f *DefaultNodeFactory) CreateNode(node *Node, deps FactoryDeps) (NodeProcessor, error) {
	return nil, fmt.Errorf("default factory does not support node type: %s", node.Type)
}
