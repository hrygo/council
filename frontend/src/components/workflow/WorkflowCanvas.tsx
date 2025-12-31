import { useCallback, useEffect, useRef } from 'react';
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
    ReactFlowProvider,
    useReactFlow,
    useNodesInitialized,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import { useConfigStore } from '../../stores/useConfigStore';
import { useWorkflow } from '../../hooks/useWorkflow';
import { transformToReactFlow, type BackendGraph } from '../../utils/graphUtils';
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore';
import { useLayoutStore } from '../../stores/useLayoutStore';
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
    layoutOptions?: {
        direction?: 'vertical' | 'horizontal';
        spacingX?: number;
        spacingY?: number;
    };
}

function WorkflowCanvasInner({
    readOnly: propReadOnly,
    mode,
    fullscreen,
    onExitFullscreen,
    workflowId,
    graph,
    onInit,
    onNodeClick,
    onPaneClick,
    layoutOptions
}: WorkflowCanvasProps) {
    const { fitView } = useReactFlow();
    const nodesInitialized = useNodesInitialized();
    const wrapperRef = useRef<HTMLDivElement>(null);
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

    // Use primitive values for dependencies to avoid infinite loops with object literals
    const { direction, spacingX, spacingY } = layoutOptions || {};

    // Transform Backend Graph to React Flow
    useEffect(() => {
        if (displayWorkflow) {
            const { nodes: initNodes, edges: initEdges } = transformToReactFlow(displayWorkflow, { direction, spacingX, spacingY });
            setNodes(initNodes);
            setEdges(initEdges);
        }
    }, [displayWorkflow, setNodes, setEdges, direction, spacingX, spacingY]);

    useEffect(() => {
        if (readOnly && displayWorkflow) {
            const { nodes: initNodes, edges: initEdges } = transformToReactFlow(displayWorkflow, { direction, spacingX, spacingY });
            useWorkflowRunStore.getState().loadWorkflow(initNodes, initEdges);
        }
    }, [readOnly, displayWorkflow, direction, spacingX, spacingY]);

    const storeNodes = useWorkflowRunStore(state => state.nodes);
    const storeEdges = useWorkflowRunStore(state => state.edges);
    const activeIds = useWorkflowRunStore(state => state.active_node_ids);

    const displayedNodes = readOnly ? storeNodes : nodes;

    // Debug Log
    useEffect(() => {
        console.log('[WorkflowCanvas] Render:', {
            readOnly,
            localNodes: nodes.length,
            storeNodes: storeNodes.length,
            displayedNodes: displayedNodes.length,
            layout: { direction, spacingX, spacingY }
        });
    }, [readOnly, nodes.length, storeNodes.length, displayedNodes.length, direction, spacingX, spacingY]);

    const panelSizes = useLayoutStore(state => state.panelSizes);

    // Observer for container resize to trigger fitView
    useEffect(() => {
        if (!wrapperRef.current || displayedNodes.length === 0) return;

        const resizeObserver = new ResizeObserver((entries) => {
            const entry = entries[0];
            const { width, height } = entry.contentRect;

            console.log('[WorkflowCanvas] Resize:', width, height);

            if (width > 0 && height > 0) {
                // Debounce fitView for resize events
                window.requestAnimationFrame(() => {
                    fitView({ padding: 0.2, duration: 200 });
                });
            }
        });

        resizeObserver.observe(wrapperRef.current);
        return () => resizeObserver.disconnect();
    }, [displayedNodes.length, fitView]);

    // Aggressive FitView Strategy:
    // When nodes are initialized, we poll fitView for a short duration to ensure 
    // it catches any layout settlements (panel expansion, animations, etc).
    useEffect(() => {
        if (nodesInitialized && displayedNodes.length > 0) {
            console.log('[WorkflowCanvas] Nodes Initialized -> Starting FitView Poller');

            // Immediate attempt
            fitView({ padding: 0.2, duration: 0 });

            // Poll every 250ms for 2 seconds
            const interval = setInterval(() => {
                if (wrapperRef.current && wrapperRef.current.clientWidth > 0) {
                    console.log('[WorkflowCanvas] Polling FitView...');
                    fitView({ padding: 0.2, duration: 300 }); // Simpler duration
                }
            }, 250);

            // Cleanup after 2 seconds
            const timeout = setTimeout(() => {
                clearInterval(interval);
                console.log('[WorkflowCanvas] Stopped FitView Poller');
            }, 2000);

            return () => {
                clearInterval(interval);
                clearTimeout(timeout);
            };
        }
    }, [nodesInitialized, displayedNodes.length, fitView, panelSizes, fullscreen]);

    const displayedNodesWithStyle = displayedNodes.map((node) => ({
        ...node,
        className: activeIds.has(node.id) ? 'node-active-pulse' : '',
        style: activeIds.has(node.id) ? { ...node.style, border: '2px solid #3B82F6' } : node.style
    }));

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
        <div ref={wrapperRef} className={`h-full w-full bg-gray-50 dark:bg-gray-900 flex flex-col ${fullscreen ? 'fixed inset-0 z-50' : ''}`}>
            {fullscreen && (
                <div className="absolute top-4 right-4 z-50">
                    <ExitFullscreenButton onClick={onExitFullscreen} />
                </div>
            )}
            <ReactFlow
                nodes={displayedNodesWithStyle}
                edges={displayedEdges}
                nodeTypes={nodeTypes}
                onNodesChange={readOnly ? undefined : onNodesChange}
                onEdgesChange={readOnly ? undefined : onEdgesChange}
                onConnect={readOnly ? undefined : onConnect}
                nodesDraggable={!readOnly}
                nodesConnectable={!readOnly}
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                onInit={onInit as any}
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

export default function WorkflowCanvas(props: WorkflowCanvasProps) {
    return (
        <ReactFlowProvider>
            <WorkflowCanvasInner {...props} />
        </ReactFlowProvider>
    );
}
