import { type FC, useState } from 'react';
import WorkflowCanvas from '../workflow/WorkflowCanvas';
import { type ReactFlowInstance } from '@xyflow/react';
import { type BackendGraph, type BackendNode } from '../../utils/graphUtils';

export const WorkflowEditor: FC = () => {
    const [rfInstance, setRfInstance] = useState<ReactFlowInstance | null>(null);
    const [graph, setGraph] = useState<BackendGraph | null>(null);
    const [prompt, setPrompt] = useState("");
    const [isGenerating, setIsGenerating] = useState(false);

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
        // This is a simplification. Ideally we use a util.
        // For MVP, we just verify we can capture the flow.
        // We actually need to map edges back to 'next_ids'.

        const nodes = rfInstance.getNodes();
        const edges = rfInstance.getEdges();

        const backendNodes: Record<string, BackendNode> = {};
        nodes.forEach(n => {
            const nextIds = edges
                .filter(e => e.source === n.id)
                .map(e => e.target);

            backendNodes[n.id] = {
                id: n.id,
                type: (n.type === 'input' ? 'start' : (n.type === 'output' ? 'end' : 'agent')), // Simple mapping
                name: n.data.label as string,
                next_ids: nextIds,
                properties: n.data // Persist extra data
            };
        });

        // Find start node (node with no incoming edges? or explicit type 'input'/'start')
        const startNode = nodes.find(n => n.type === 'input' || n.type === 'start');

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
        <div className="h-screen flex flex-col">
            <header className="h-14 border-b px-4 flex items-center justify-between bg-white dark:bg-gray-800">
                <div className="flex items-center gap-4">
                    <h1 className="font-bold">Workflow Builder</h1>
                    <div className="flex gap-2">
                        <input
                            className="border rounded px-2 py-1 text-sm w-64"
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
            <div className="flex-1 overflow-hidden">
                <WorkflowCanvas
                    graph={graph}
                    onInit={setRfInstance}
                />
            </div>
        </div>
    );
};
