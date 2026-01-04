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
import {
    AgentNode,
    VoteNode,
    LoopNode,
    FactCheckNode,
    HumanReviewNode,
    StartNode,
    EndNode,
    ParallelNode
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
    parallel: ParallelNode,
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
    const hasFittedViewRef = useRef(false); // Track if initial fitView has been done
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



    // Track previous node count to detect first load
    const prevNodeCountRef = useRef(0);

    // Track the last container size used for fitView to detect significant changes
    const lastFitSizeRef = useRef<{ width: number; height: number } | null>(null);
    const fitViewTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

    // CRITICAL: When nodes are first loaded (0 -> N), schedule fitView
    // This handles the case where onInit fires before nodes are loaded
    useEffect(() => {
        const prevCount = prevNodeCountRef.current;
        const currCount = displayedNodes.length;
        prevNodeCountRef.current = currCount;

        // First load: nodes went from 0 to > 0
        if (prevCount === 0 && currCount > 0) {
            // Wait a bit for ReactFlow to measure the nodes
            const timeoutId = setTimeout(() => {
                fitView({ padding: 0.2, duration: 200 });
                hasFittedViewRef.current = true;
            }, 300);
            return () => clearTimeout(timeoutId);
        }
    }, [displayedNodes.length, fitView]);

    // Reset fit state when nodes change (e.g., new workflow loaded)
    useEffect(() => {
        hasFittedViewRef.current = false;
        lastFitSizeRef.current = null;
    }, [displayedNodes.length]);

    // Debounced fitView that waits for layout to stabilize
    // Using 500ms delay to ensure ReactFlow has time to measure nodes internally
    const debouncedFitView = useCallback((width: number, height: number) => {
        // Clear any pending fitView
        if (fitViewTimeoutRef.current) {
            clearTimeout(fitViewTimeoutRef.current);
        }

        // Schedule fitView after layout stabilizes (500ms debounce - longer to ensure nodes are measured)
        fitViewTimeoutRef.current = setTimeout(() => {
            const lastSize = lastFitSizeRef.current;
            const sizeChanged = !lastSize ||
                Math.abs(lastSize.width - width) > 50 ||
                Math.abs(lastSize.height - height) > 50;

            if (sizeChanged) {
                fitView({ padding: 0.2, duration: 200 });
                lastFitSizeRef.current = { width, height };
                hasFittedViewRef.current = true;
            }
        }, 500);
    }, [fitView]);

    // Observer for container resize
    useEffect(() => {
        if (!wrapperRef.current || displayedNodes.length === 0) return;

        const resizeObserver = new ResizeObserver((entries) => {
            const entry = entries[0];
            const { width, height } = entry.contentRect;

            if (width > 0 && height > 0) {
                debouncedFitView(width, height);
            }
        });

        resizeObserver.observe(wrapperRef.current);
        return () => {
            resizeObserver.disconnect();
            if (fitViewTimeoutRef.current) {
                clearTimeout(fitViewTimeoutRef.current);
            }
        };
    }, [displayedNodes.length, debouncedFitView]);

    // CRITICAL: When nodesInitialized becomes true, ALWAYS fitView (regardless of hasFittedViewRef)
    // because the initial fitView while nodesInitialized=false is ineffective
    useEffect(() => {
        if (nodesInitialized && displayedNodes.length > 0) {
            // Use requestAnimationFrame to ensure DOM is ready
            requestAnimationFrame(() => {
                if (wrapperRef.current) {
                    const { clientWidth, clientHeight } = wrapperRef.current;
                    if (clientWidth > 0 && clientHeight > 0) {
                        fitView({ padding: 0.2, duration: 200 });
                        lastFitSizeRef.current = { width: clientWidth, height: clientHeight };
                        hasFittedViewRef.current = true;
                    }
                }
            });
        }
    }, [nodesInitialized, displayedNodes.length, fitView]);

    // Debug active nodes
    useEffect(() => {
        console.log('[Canvas] Active node IDs:', Array.from(activeIds));
    }, [activeIds]);

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
                // Force remount when nodes first become available to ensure proper initialization
                key={displayedNodes.length > 0 ? 'has-nodes' : 'empty'}
                nodes={displayedNodesWithStyle}
                edges={displayedEdges}
                nodeTypes={nodeTypes}
                onNodesChange={readOnly ? undefined : onNodesChange}
                onEdgesChange={readOnly ? undefined : onEdgesChange}
                onConnect={readOnly ? undefined : onConnect}
                nodesDraggable={!readOnly}
                nodesConnectable={!readOnly}
                onInit={(instance) => {
                    // Call user-provided onInit if any
                    if (onInit) {
                        // eslint-disable-next-line @typescript-eslint/no-explicit-any
                        (onInit as any)(instance);
                    }
                    // Schedule fitView after a delay to ensure nodes are measured
                    setTimeout(() => {
                        instance.fitView({ padding: 0.2 });
                    }, 100);
                }}
                onNodeClick={onNodeClick}
                onPaneClick={onPaneClick}
                minZoom={0.1}
                defaultViewport={{ x: 0, y: 0, zoom: 0.5 }}
                onlyRenderVisibleElements={false}
                fitView
                fitViewOptions={{ padding: 0.2 }}
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
