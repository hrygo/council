package middleware

import (
	"context"

	"github.com/hrygo/council/internal/core/workflow"
)

// CircuitBreaker intercepts nodes to enforce safety limits
type CircuitBreaker struct {
	MaxRecursionDepth int
}

func NewCircuitBreaker(maxDepth int) *CircuitBreaker {
	return &CircuitBreaker{MaxRecursionDepth: maxDepth}
}

func (cb *CircuitBreaker) Name() string {
	return "LogicCircuitBreaker"
}

func (cb *CircuitBreaker) BeforeNodeExecution(ctx context.Context, session *workflow.Session, node *workflow.Node) error {
	// Simple depth check: For Loop nodes, we need to track iteration count.
	// Since Session doesn't expose iteration map cleanly yet, we can check a generic depth counter if we had one.
	// For MVP: We assume Loop nodes update a counter in Session.Inputs or a special context key?
	// TDD says: "Recursion Depth > 10".
	// Let's defer strict implementation until Loop node is fully defined.
	// Current logic: Check strict loop count if available.

	// Real implementation would check session.RecursionMap[node.ID]
	// Since session.go is basic, we skip detailed logic for now but provide the hook.
	return nil
}

func (cb *CircuitBreaker) AfterNodeExecution(ctx context.Context, session *workflow.Session, node *workflow.Node, output map[string]interface{}) (map[string]interface{}, error) {
	// Post-execution check: Token Velocity?
	// If output contains token usage, we can check it.
	return output, nil
}
