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
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useWorkflow } from '../../hooks/useWorkflow';
import { transformToReactFlow } from '../../utils/graphUtils';
import { useConnectStore } from '../../stores/useConnectStore';

interface WorkflowCanvasProps {
    readOnly?: boolean;
    fullscreen?: boolean;
    onExitFullscreen?: () => void;
    workflowId?: string; // Optional for now
}

export default function WorkflowCanvas({ readOnly, fullscreen, onExitFullscreen, workflowId }: WorkflowCanvasProps) {
    const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
    const { workflow, fetchWorkflow } = useWorkflow();
    const { lastMessage } = useConnectStore();

    // Fetch Workflow on Mount
    useEffect(() => {
        if (workflowId) {
            fetchWorkflow(workflowId);
        }
    }, [workflowId, fetchWorkflow]);

    // Transform Backend Graph to React Flow
    useEffect(() => {
        if (workflow) {
            const { nodes: initNodes, edges: initEdges } = transformToReactFlow(workflow);
            setNodes(initNodes);
            setEdges(initEdges);
        }
    }, [workflow, setNodes, setEdges]);

    // Listen for WebSocket Events to Highlight Nodes
    useEffect(() => {
        if (!lastMessage || !readOnly) return; // Only highlight in readOnly/Run mode

        // Event: "node:started", "node:completed"
        // Data: { node_id: "..." }

        if (lastMessage.type === 'node:started') {
            const runningNodeId = lastMessage.data?.node_id;
            if (runningNodeId) {
                setNodes(nds => nds.map(node => {
                    if (node.id === runningNodeId) {
                        return { ...node, style: { ...node.style, border: '2px solid blue', background: '#e0f2fe' } };
                    }
                    return node;
                }));
            }
        }

        if (lastMessage.type === 'node:completed') {
            const completedNodeId = lastMessage.data?.node_id;
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
                fitView
            >
                <Background color="#ccc" gap={20} />
                <Controls showInteractive={!readOnly} />
            </ReactFlow>
        </div>
    );
}
