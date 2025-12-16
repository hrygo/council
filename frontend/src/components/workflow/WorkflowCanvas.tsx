import { useCallback } from 'react';
import {
    ReactFlow,
    Background,
    Controls,
    useNodesState,
    useEdgesState,
    addEdge,
    type Connection,
    type Edge,
    MarkerType,
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';

const initialNodes = [
    {
        id: 'start',
        position: { x: 250, y: 50 },
        data: { label: 'Start' },
        type: 'input',
    },
    {
        id: 'process',
        position: { x: 250, y: 150 },
        data: { label: 'Process Node' },
    },
    {
        id: 'end',
        position: { x: 250, y: 250 },
        data: { label: 'End' },
        type: 'output',
    },
];

const initialEdges: Edge[] = [
    { id: 'e1-2', source: 'start', target: 'process', markerEnd: { type: MarkerType.ArrowClosed } },
    { id: 'e2-3', source: 'process', target: 'end', markerEnd: { type: MarkerType.ArrowClosed } },
];

interface WorkflowCanvasProps {
    readOnly?: boolean;
    fullscreen?: boolean;
    onExitFullscreen?: () => void;
}

export default function WorkflowCanvas({ readOnly, fullscreen, onExitFullscreen }: WorkflowCanvasProps) {
    const [nodes, , onNodesChange] = useNodesState(initialNodes);
    const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges);

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
