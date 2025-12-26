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
    type NodeTypes,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useConfigStore } from '../../stores/useConfigStore';
import { useWorkflow } from '../../hooks/useWorkflow';
import { transformToReactFlow, type BackendGraph } from '../../utils/graphUtils';
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore';
import {
    AgentNode,
    VoteNode,
    LoopNode,
    FactCheckNode,
    HumanReviewNode,
    StartNode,
    EndNode
} from './nodes/CustomNodes';
import { ExitFullscreenButton } from '../ui/ExitFullscreenButton';

const nodeTypes: NodeTypes = {
    agent: AgentNode,
    vote: VoteNode,
    loop: LoopNode,
    fact_check: FactCheckNode,
    human_review: HumanReviewNode,
    start: StartNode,
    end: EndNode,
};

interface WorkflowCanvasProps {
    readOnly?: boolean;
    mode?: 'edit' | 'run';
    fullscreen?: boolean;
    onExitFullscreen?: () => void;
    workflowId?: string;
    graph?: BackendGraph | null;
    onInit?: (instance: ReactFlowInstance) => void;
    onNodeClick?: (event: React.MouseEvent, node: Node) => void;
    onPaneClick?: (event: React.MouseEvent) => void;
}

export default function WorkflowCanvas({
    readOnly: propReadOnly,
    mode,
    fullscreen,
    onExitFullscreen,
    workflowId,
    graph,
    onInit,
    onNodeClick,
    onPaneClick
}: WorkflowCanvasProps) {
    const readOnly = propReadOnly || mode === 'run';
    const [nodes, setNodes, onNodesChange] = useNodesState<Node>([]);
    const [edges, setEdges, onEdgesChange] = useEdgesState<Edge>([]);
    const { workflow: fetchedWorkflow, fetchWorkflow } = useWorkflow();
    const theme = useConfigStore((state) => state.theme);

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

    useEffect(() => {
        if (readOnly && displayWorkflow) {
            const { nodes: initNodes, edges: initEdges } = transformToReactFlow(displayWorkflow);
            useWorkflowRunStore.getState().loadWorkflow(initNodes, initEdges);
        }
    }, [readOnly, displayWorkflow]);

    const storeNodes = useWorkflowRunStore(state => state.nodes);
    const storeEdges = useWorkflowRunStore(state => state.edges);
    const activeIds = useWorkflowRunStore(state => state.active_node_ids);

    const displayedNodes = readOnly ? storeNodes.map((node) => ({
        ...node,
        className: activeIds.has(node.id) ? 'node-active-pulse' : '',
        style: activeIds.has(node.id) ? { ...node.style, border: '2px solid #3B82F6' } : node.style
    })) : nodes;

    const displayedEdges = readOnly ? storeEdges : edges;

    const onConnect = useCallback(
        (params: Connection) => {
            if (!readOnly) {
                setEdges((eds) => addEdge(params, eds));
            }
        },
        [setEdges, readOnly],
    );

    // Determine grid color based on theme
    const isSystemDark = typeof window !== 'undefined' ? window.matchMedia('(prefers-color-scheme: dark)').matches : false;
    const isDark = theme === 'dark' || (theme === 'system' && isSystemDark);
    const gridColor = isDark ? '#374151' : '#ccc'; // gray-700 : gray-300

    return (
        <div className={`h-full w-full bg-gray-50 dark:bg-gray-900 flex flex-col ${fullscreen ? 'fixed inset-0 z-50' : ''}`}>
            {fullscreen && (
                <div className="absolute top-4 right-4 z-50">
                    <ExitFullscreenButton onClick={onExitFullscreen} />
                </div>
            )}
            <ReactFlow
                nodes={displayedNodes}
                edges={displayedEdges}
                nodeTypes={nodeTypes}
                onNodesChange={readOnly ? undefined : onNodesChange}
                onEdgesChange={readOnly ? undefined : onEdgesChange}
                onConnect={readOnly ? undefined : onConnect}
                nodesDraggable={!readOnly}
                nodesConnectable={!readOnly}
                onInit={onInit}
                onNodeClick={onNodeClick}
                onPaneClick={onPaneClick}
                fitView
            >
                <Background color={gridColor} gap={20} />
                <Controls showInteractive={!readOnly} />
            </ReactFlow>
        </div>
    );
}
