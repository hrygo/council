import { type Node, type Edge, MarkerType } from '@xyflow/react';
import dagre from 'dagre';

// Types mirroring Backend GraphDefinition
export interface BackendNode {
    node_id: string;
    type: string; // "start", "end", "agent", "llm", "tool", "parallel", "sequence"
    name: string;
    next_ids?: string[];
    properties?: Record<string, unknown>;
}

export interface BackendGraph {
    workflow_uuid: string;
    name: string;
    description: string;
    nodes: Record<string, BackendNode>;
    start_node_id: string;
}

interface LayoutOptions {
    direction?: 'vertical' | 'horizontal';
    spacingX?: number; // Spacing between nodes horizontally (nodesep)
    spacingY?: number; // Spacing between nodes vertically (ranksep)
}

export const transformToReactFlow = (
    graph: BackendGraph,
    options: LayoutOptions = {}
): { nodes: Node[]; edges: Edge[] } => {
    // Safety check for nodes
    if (!graph.nodes) {
        return { nodes: [], edges: [] };
    }

    const {
        direction = 'vertical',
        spacingX = 120, // Horizontal separation
        spacingY = 150  // Vertical separation
    } = options;

    const dagreGraph = new dagre.graphlib.Graph();
    dagreGraph.setDefaultEdgeLabel(() => ({}));

    const isHorizontal = direction === 'horizontal';
    dagreGraph.setGraph({
        rankdir: isHorizontal ? 'LR' : 'TB',
        // align: 'UL', // Removed to allow center alignment
        nodesep: spacingX,
        ranksep: spacingY,
    });

    const nodeWidth = 250;
    const nodeHeight = 150; // Approximate height of our cards

    const nodes: Node[] = [];
    const edges: Edge[] = [];

    // Add nodes to dagre
    Object.values(graph.nodes).forEach((node) => {
        dagreGraph.setNode(node.node_id, { width: nodeWidth, height: nodeHeight });
    });

    // Add edges to dagre
    Object.values(graph.nodes).forEach((node) => {
        if (node.next_ids) {
            node.next_ids.forEach((nextId) => {
                // Ensure target exists
                if (graph.nodes[nextId]) {
                    dagreGraph.setEdge(node.node_id, nextId);

                    edges.push({
                        id: `e-${node.node_id}-${nextId}`,
                        source: node.node_id,
                        target: nextId,
                        type: 'default', // Bezier curves for smoother look
                        markerEnd: { type: MarkerType.ArrowClosed },
                        animated: false,
                        style: { strokeWidth: 2, stroke: '#94a3b8' } // Slate-400
                    });
                }
            });
        }
    });

    // Calculate layout
    dagre.layout(dagreGraph);

    // Create React Flow Nodes
    Object.values(graph.nodes).forEach((node) => {
        const nodeWithPosition = dagreGraph.node(node.node_id);

        // We need to pass a slightly different position in order to notify react flow about the center of the node
        // Dagre layout is based on the center of the node, but RF uses top-left.
        // Wait, Dagre uses Center-Center by default for calculations, but returns Top-Left `x,y`?
        // No, Dagre returns Center x,y.
        // React Flow positions are Top-Left.
        // So: x = center_x - width / 2

        nodes.push({
            id: node.node_id,
            type: mapNodeType(node.type),
            position: {
                x: nodeWithPosition.x - nodeWidth / 2,
                y: nodeWithPosition.y - nodeHeight / 2,
            },
            data: { label: node.name, ...node.properties },
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
