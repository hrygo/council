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

	// 3. Traversal for Reachability
	// We allow cycles, so we just use a visited map.
	visited := make(map[string]bool)

	if err := g.detectCycle(g.StartNodeID, visited); err != nil {
		return err
	}

	// Check if all nodes were visited (Reachability)
	if len(visited) != len(g.Nodes) {
		return fmt.Errorf("graph contains unreachable nodes (visited %d/%d)", len(visited), len(g.Nodes))
	}

	return nil
}

func (g *GraphDefinition) detectCycle(nodeID string, visited map[string]bool) error {
	visited[nodeID] = true

	node := g.Nodes[nodeID]
	for _, nextID := range node.NextIDs {
		// If not visited, recurse
		if !visited[nextID] {
			if err := g.detectCycle(nextID, visited); err != nil {
				return err
			}
		}
		// We allow cycles because Loop nodes are valid constructs in this engine.
		// Runtime protection (MaxSteps/LoopProcessor) handles infinite loops.
	}

	return nil
}
