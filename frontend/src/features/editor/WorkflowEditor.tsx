import { type FC, useState, useCallback } from 'react';
import WorkflowCanvas from '../../components/workflow/WorkflowCanvas';
import { type ReactFlowInstance, type Node } from '@xyflow/react';
import { type BackendGraph, type BackendNode } from '../../utils/graphUtils';
import { PropertyPanel } from './components/PropertyPanel/PropertyPanel';
import type { WorkflowNode } from '../../types/workflow';

export const WorkflowEditor: FC = () => {
    const [rfInstance, setRfInstance] = useState<ReactFlowInstance | null>(null);
    const [graph, setGraph] = useState<BackendGraph | null>(null);
    const [prompt, setPrompt] = useState("");
    const [isGenerating, setIsGenerating] = useState(false);

    // Property Panel State
    const [selectedNodeId, setSelectedNodeId] = useState<string | null>(null);
    const [selectedNode, setSelectedNode] = useState<WorkflowNode | null>(null);

    // Handle Selection
    const handleNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
        setSelectedNodeId(node.id);
        setSelectedNode(node as unknown as WorkflowNode);
    }, []);

    const handlePaneClick = useCallback(() => {
        setSelectedNodeId(null);
        setSelectedNode(null);
    }, []);

    // Handle Node Updates
    const handleNodeUpdate = useCallback((nodeId: string, newData: Record<string, unknown>) => {
        if (!rfInstance) return;

        rfInstance.setNodes((nodes) =>
            nodes.map((node) => {
                if (node.id === nodeId) {
                    const updatedNode = {
                        ...node,
                        data: { ...node.data, ...newData }
                    };
                    // Update local selected node state if it's the one being edited
                    if (selectedNodeId === nodeId) {
                        setSelectedNode(updatedNode as unknown as WorkflowNode);
                    }
                    return updatedNode;
                }
                return node;
            })
        );
    }, [rfInstance, selectedNodeId]);

    const handleNodeDelete = useCallback((nodeId: string) => {
        if (!rfInstance) return;

        rfInstance.setNodes((nodes) => nodes.filter((n) => n.id !== nodeId));
        rfInstance.setEdges((edges) => edges.filter((e) => e.source !== nodeId && e.target !== nodeId));

        setSelectedNodeId(null);
        setSelectedNode(null);
    }, [rfInstance]);


    // We can use a custom hook or fetch directly for Generation since it's specific
    const handleGenerate = async () => {
        setIsGenerating(true);
        try {
            const res = await fetch('/api/v1/workflows/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ prompt })
            });
            const data = await res.json();
            if (data.graph) {
                setGraph(data.graph);
            }
        } catch (e) {
            console.error(e);
            alert("Generation failed");
        } finally {
            setIsGenerating(false);
        }
    };

    const handleSave = async () => {
        if (!rfInstance) return;

        // Convert RF nodes/edges back to Backend Graph
        const nodes = rfInstance.getNodes();
        const edges = rfInstance.getEdges();

        const backendNodes: Record<string, BackendNode> = {};
        nodes.forEach(n => {
            const nextIds = edges
                .filter(e => e.source === n.id)
                .map(e => e.target);

            backendNodes[n.id] = {
                id: n.id,
                type: (n.type as string) || 'agent', // Default/Fallback
                name: n.data.label as string,
                next_ids: nextIds,
                properties: n.data // Persist extra data
            };
        });

        // Find start node options
        const startNode = nodes.find(n => n.type === 'start');

        const payload = {
            id: graph?.id || undefined, // undefined to create new if not exists
            name: graph?.name || "Untitled Workflow",
            description: graph?.description || "Created via Builder",
            start_node_id: startNode ? startNode.id : (nodes[0]?.id || ""),
            nodes: backendNodes
        };

        try {
            const method = payload.id ? 'PUT' : 'POST';
            const url = payload.id ? `/api/v1/workflows/${payload.id}` : '/api/v1/workflows';

            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                const saved = await res.json();
                setGraph(saved); // Update ID if new
                alert("Saved successfully!");
            } else {
                alert("Save failed");
            }
        } catch (e) {
            console.error(e);
            alert("Error saving");
        }
    };

    return (
        <div className="h-screen flex flex-col relative">
            <header className="h-14 border-b px-4 flex items-center justify-between bg-white dark:bg-gray-800 shrink-0 z-10">
                <div className="flex items-center gap-4">
                    <h1 className="font-bold">Workflow Builder</h1>
                    <div className="flex gap-2">
                        <input
                            className="border rounded px-2 py-1 text-sm w-64 dark:bg-gray-700 dark:border-gray-600"
                            placeholder="Describe workflow to generate..."
                            value={prompt}
                            onChange={e => setPrompt(e.target.value)}
                        />
                        <button
                            onClick={handleGenerate}
                            disabled={!prompt || isGenerating}
                            className="px-3 py-1 bg-purple-600 text-white rounded text-sm disabled:opacity-50"
                        >
                            {isGenerating ? "Generating..." : "Generate AI"}
                        </button>
                    </div>
                </div>
                <div className="flex gap-2">
                    <button onClick={handleSave} className="px-3 py-1 bg-blue-500 text-white rounded text-sm">Save</button>
                    <button className="px-3 py-1 bg-green-500 text-white rounded text-sm">Run Session</button>
                </div>
            </header>
            <div className="flex-1 overflow-hidden relative">
                <WorkflowCanvas
                    graph={graph}
                    onInit={setRfInstance}
                    onNodeClick={handleNodeClick}
                    onPaneClick={handlePaneClick}
                />

                {selectedNodeId && selectedNode && (
                    <PropertyPanel
                        node={selectedNode}
                        onUpdate={handleNodeUpdate}
                        onDelete={handleNodeDelete}
                        onClose={() => {
                            setSelectedNodeId(null);
                            setSelectedNode(null);
                        }}
                    />
                )}
            </div>
        </div>
    );
};

