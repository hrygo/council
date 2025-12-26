import { type Node, type Edge, MarkerType } from '@xyflow/react';

// Types mirroring Backend GraphDefinition
export interface BackendNode {
    node_id: string;
    type: string; // "start", "end", "agent", "llm", "tool", "parallel", "sequence"
    name: string;
    next_ids?: string[];
    properties?: Record<string, unknown>;
}

export interface BackendGraph {
    workflow_id: string;
    name: string;
    description: string;
    nodes: Record<string, BackendNode>;
    start_node_id: string;
}

export const transformToReactFlow = (graph: BackendGraph): { nodes: Node[]; edges: Edge[] } => {
    const nodes: Node[] = [];
    const edges: Edge[] = [];
    const visited = new Set<string>();

    // Simple BFS/Level-based layout calculation
    const levels: Record<string, number> = {};
    const queue: { id: string; level: number }[] = [];

    if (graph.start_node_id) {
        queue.push({ id: graph.start_node_id, level: 0 });
        visited.add(graph.start_node_id);
    }

    // Nodes processing & Edge creation
    // We do a first pass to traverse and establish edges and levels
    // Note: This assumes a DAG. Cycles might need special handling but for now BFS works for levels.

    // Safety check for nodes
    if (!graph.nodes) {
        return { nodes: [], edges: [] };
    }

    const nodeIds = Object.keys(graph.nodes);
    // If start_node_id is missing or invalid, just map all nodes linearly or loosely?
    // Let's assume start_node_id gives us the entry.

    let head = 0;
    while (head < queue.length) {
        const { id, level } = queue[head++];
        levels[id] = level;

        const node = graph.nodes[id];
        if (node && node.next_ids) {
            node.next_ids.forEach(nextId => {
                // Create Edge
                edges.push({
                    id: `e-${id}-${nextId}`,
                    source: id,
                    target: nextId,
                    type: 'smoothstep',
                    markerEnd: { type: MarkerType.ArrowClosed },
                    animated: false,
                });

                if (!visited.has(nextId)) {
                    visited.add(nextId);
                    queue.push({ id: nextId, level: level + 1 });
                }
            });
        }
    }

    // Handle disconnected nodes (if any)
    nodeIds.forEach(id => {
        if (!visited.has(id)) {
            levels[id] = 0; // Default to level 0 for disconnected
            // We might want to add them to nodes list
        }
    });

    // Create React Flow Nodes with calculated positions
    // Simple layout: x = level * 250, y = indexInLevel * 100
    const levelCounts: Record<number, number> = {};

    nodeIds.forEach(id => {
        const node = graph.nodes[id];
        const level = levels[id] || 0;

        if (!levelCounts[level]) levelCounts[level] = 0;
        const indexInLevel = levelCounts[level]++;

        nodes.push({
            id: node.node_id,
            type: mapNodeType(node.type), // Map backend type to RF type
            position: { x: indexInLevel * 200 + 100, y: level * 150 + 50 }, // Vertical layout? Or horizontal?
            // Vertical Layout:
            // x: indexInLevel * 200
            // y: level * 150
            data: { label: node.name, ...node.properties },
            // 'input' type for start, 'output' for end?
            // Actually React Flow 'input'/'output' logic is about handles. 
            // We can use default for all or customize.
            // Let's stick to default for now, or 'input' for start.
        });
    });

    return { nodes, edges };
};

const mapNodeType = (backendType: string): string => {
    // Map backend types directly to registered custom node types
    // 'agent', 'vote', 'loop', 'fact_check', 'human_review' match 1:1
    switch (backendType) {
        // Core/Flow types
        case 'start': return 'start';
        case 'end': return 'end';

        // Task types (pass through or explicit map if keys differ)
        case 'agent': return 'agent';
        case 'vote': return 'vote';
        case 'loop': return 'loop';
        case 'fact_check': return 'fact_check';
        case 'human_review': return 'human_review';

        // Fallback for visual safety, though 'default' brings the white-node issue.
        // Better to fallback to 'agent' if unknown to get the BaseNode style?
        // Or keep 'default' but we must ensure we don't use it for known types.
        default: return 'agent'; // Fallback to agent (generic node) to ensure styles
    }
};
