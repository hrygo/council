package workflow

import (
	"errors"
	"fmt"
)

// Validate checks if the graph definition is valid
// 1. Start node exists
// 2. All next_ids point to existing nodes
// 3. No cycles
// 4. All nodes are reachable from Start
func (g *GraphDefinition) Validate() error {
	if g.Nodes == nil {
		return errors.New("nodes map is nil")
	}

	// 1. Check Start Node
	if _, ok := g.Nodes[g.StartNodeID]; !ok {
		return fmt.Errorf("start node %s not found", g.StartNodeID)
	}

	// 2. Check all links
	for id, node := range g.Nodes {
		if node.ID != id {
			return fmt.Errorf("node ID mismatch: map key %s vs node.ID %s", id, node.ID)
		}
		for _, nextID := range node.NextIDs {
			if _, ok := g.Nodes[nextID]; !ok {
				return fmt.Errorf("node %s points to non-existent node %s", id, nextID)
			}
		}
	}

	// 3. Cycle Detection & 4. Reachability
	// We can do both with a traversal
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)

	if err := g.detectCycle(g.StartNodeID, visited, recursionStack); err != nil {
		return err
	}

	// Check if all nodes were visited (Reachability)
	if len(visited) != len(g.Nodes) {
		return fmt.Errorf("graph contains unreachable nodes (visited %d/%d)", len(visited), len(g.Nodes))
	}

	return nil
}

func (g *GraphDefinition) detectCycle(nodeID string, visited, recursionStack map[string]bool) error {
	visited[nodeID] = true
	recursionStack[nodeID] = true

	node := g.Nodes[nodeID]
	for _, nextID := range node.NextIDs {
		// If not visited, recurse
		if !visited[nextID] {
			if err := g.detectCycle(nextID, visited, recursionStack); err != nil {
				return err
			}
		} else if recursionStack[nextID] {
			// If already in recursion stack, it's a cycle
			return fmt.Errorf("cycle detected involving node %s", nextID)
		}
	}

	recursionStack[nodeID] = false
	return nil
}
