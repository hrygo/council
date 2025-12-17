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
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore';

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

    // Listen for WebSocket Events via hook (should be called in parent or here? better in parent or high level layout)
    // But since this is a presentation component, maybe the hook should be used in the page.
    // However, for highlighting, we need to merge passed `nodes` with `activeNodeIds` or usage `runNodes` if in readOnly mode.

    // Strategy:
    // If readOnly (Run Mode), we use `runNodes` from store which are synchronized with execution state.
    // If not readOnly (Edit Mode), we use local `nodes`.

    // Actually, `runNodes` in store are populated via `loadWorkflow`.
    // Let's assume onInit or prop update triggers `loadWorkflow`.

    useEffect(() => {
        if (readOnly && displayWorkflow) {
            const { nodes: initNodes, edges: initEdges } = transformToReactFlow(displayWorkflow);
            useWorkflowRunStore.getState().loadWorkflow(initNodes, initEdges);
        }
    }, [readOnly, displayWorkflow]);

    // Use store nodes if readOnly, otherwise local state
    // But wait, React Flow needs `nodes` passed to it.
    // If readOnly, we should sync `nodes` from store or just derive styles?
    // Using store nodes directly might be cleaner for ReadOnly mode.

    const storeNodes = useWorkflowRunStore(state => state.nodes);
    const storeEdges = useWorkflowRunStore(state => state.edges);
    const activeIds = useWorkflowRunStore(state => state.activeNodeIds);

    // Merge styles for active nodes if using local nodes (hybrid approach) OR just use store nodes.
    // Let's use storeNodes for readOnly mode.

    // We need to type cast or ensure runtime nodes are compatible with ReactFlow Node type
    const displayedNodes = readOnly ? storeNodes.map((node) => ({
        ...node,
        className: activeIds.has(node.id) ? 'node-active-pulse' : '',
        style: activeIds.has(node.id) ? { ...node.style, border: '2px solid #3B82F6' } : node.style
    })) : nodes;

    const displayedEdges = readOnly ? storeEdges : edges;

    // Remove old WebSocket effect since we use useWorkflowEvents hook globally or at page level
    // But wait, the previous code had local effect. We are replacing it.

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
                nodes={displayedNodes}
                edges={displayedEdges}
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
