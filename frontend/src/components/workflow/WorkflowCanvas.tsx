import { useCallback, useEffect } from 'react';
import {
    ReactFlow,
    Background,
    Controls,
    useNodesState,
    useEdgesState,
    addEdge,
    type Connection,
    type Node,
    type Edge,
    type ReactFlowInstance,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useWorkflow } from '../../hooks/useWorkflow';
import { transformToReactFlow, type BackendGraph } from '../../utils/graphUtils';
import { useConnectStore } from '../../stores/useConnectStore';

interface WorkflowCanvasProps {
    readOnly?: boolean;
    fullscreen?: boolean;
    onExitFullscreen?: () => void;
    workflowId?: string;
    graph?: BackendGraph | null;
    onInit?: (instance: ReactFlowInstance) => void;
}

export default function WorkflowCanvas({
    readOnly,
    fullscreen,
    onExitFullscreen,
    workflowId,
    graph,
    onInit
}: WorkflowCanvasProps) {
    const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
    const { workflow: fetchedWorkflow, fetchWorkflow } = useWorkflow();
    const { lastMessage } = useConnectStore();

    // Determine effective workflow to display: prop > fetched
    const displayWorkflow = graph || fetchedWorkflow;

    // Fetch if ID provided and no external graph
    useEffect(() => {
        if (workflowId && !graph) {
            fetchWorkflow(workflowId);
        }
    }, [workflowId, graph, fetchWorkflow]);

    // Transform Backend Graph to React Flow
    useEffect(() => {
        if (displayWorkflow) {
            const { nodes: initNodes, edges: initEdges } = transformToReactFlow(displayWorkflow);
            setNodes(initNodes);
            setEdges(initEdges);
        }
    }, [displayWorkflow, setNodes, setEdges]);

    // Listen for WebSocket Events to Highlight Nodes
    useEffect(() => {
        // Safe access to lastMessage properties
        const msg = lastMessage as Record<string, unknown> | null;
        if (!msg || !readOnly) return;

        if (msg.type === 'node:started') {
            const data = msg.data as Record<string, unknown> | undefined;
            const runningNodeId = data?.node_id as string | undefined;
            if (runningNodeId) {
                setNodes(nds => nds.map(node => {
                    if (node.id === runningNodeId) {
                        return { ...node, style: { ...node.style, border: '2px solid blue', background: '#e0f2fe' } };
                    }
                    return node;
                }));
            }
        }

        if (msg.type === 'node:completed') {
            const data = msg.data as Record<string, unknown> | undefined;
            const completedNodeId = data?.node_id as string | undefined;
            if (completedNodeId) {
                setNodes(nds => nds.map(node => {
                    if (node.id === completedNodeId) {
                        return { ...node, style: { ...node.style, border: '2px solid green', background: '#dcfce7' } };
                    }
                    return node;
                }));
            }
        }

    }, [lastMessage, readOnly, setNodes]);

    const onConnect = useCallback(
        (params: Connection) => {
            if (!readOnly) {
                setEdges((eds) => addEdge(params, eds));
            }
        },
        [setEdges, readOnly],
    );

    return (
        <div className={`h-full w-full bg-gray-50 flex flex-col ${fullscreen ? 'fixed inset-0 z-50 bg-white' : ''}`}>
            {fullscreen && (
                <div className="absolute top-4 right-4 z-50">
                    <button onClick={onExitFullscreen} className="bg-white/80 p-2 rounded shadow text-sm">Exit Fullscreen</button>
                </div>
            )}
            <ReactFlow
                nodes={nodes}
                edges={edges}
                onNodesChange={readOnly ? undefined : onNodesChange}
                onEdgesChange={readOnly ? undefined : onEdgesChange}
                onConnect={readOnly ? undefined : onConnect}
                nodesDraggable={!readOnly}
                nodesConnectable={!readOnly}
                onInit={onInit}
                fitView
            >
                <Background color="#ccc" gap={20} />
                <Controls showInteractive={!readOnly} />
            </ReactFlow>
        </div>
    );
}
