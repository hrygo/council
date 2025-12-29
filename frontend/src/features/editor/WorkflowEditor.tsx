import { type FC, useState, useCallback } from 'react';
import WorkflowCanvas from '../../components/workflow/WorkflowCanvas';
import { type ReactFlowInstance, type Node } from '@xyflow/react';
import { type BackendGraph, type BackendNode } from '../../utils/graphUtils';
import { PropertyPanel } from './components/PropertyPanel/PropertyPanel';
import type { WorkflowNode } from '../../types/workflow';
import { TemplateSidebar } from './components/TemplateSidebar';
import { SaveTemplateModal } from './components/SaveTemplateModal';
import { WizardMode } from './components/Wizard/WizardMode';
import type { Template } from '../../types/template';
import { Wand2, LayoutTemplate, Save } from 'lucide-react';
import { CostEstimator } from '../execution/components/CostEstimator';
import { useToast } from '../../components/ui/useToast';

export const WorkflowEditor: FC = () => {
    const [rfInstance, setRfInstance] = useState<ReactFlowInstance | null>(null);
    const [graph, setGraph] = useState<BackendGraph | null>(null);
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



    const { success, error } = useToast();

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
                node_id: n.id,
                type: (n.type as string) || 'agent', // Default/Fallback
                name: n.data.label as string,
                next_ids: nextIds,
                properties: n.data // Persist extra data
            };
        });

        // Find start node options
        const startNode = nodes.find(n => n.type === 'start');

        const payload = {
            workflow_uuid: graph?.workflow_uuid || undefined, // undefined to create new if not exists
            name: graph?.name || "Untitled Workflow",
            description: graph?.description || "Created via Builder",
            start_node_id: startNode ? startNode.id : (nodes[0]?.id || ""),
            nodes: backendNodes
        };

        try {
            // Determine method based on ID presence AND ensure empty string is treated as new
            const isNew = !payload.workflow_uuid || payload.workflow_uuid === '';
            const method = isNew ? 'POST' : 'PUT';
            const url = isNew ? '/api/v1/workflows' : `/api/v1/workflows/${payload.workflow_uuid}`;

            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (res.ok) {
                const saved = await res.json();
                setGraph(saved); // Update ID if new
                success("Workflow saved successfully!");
            } else {
                const errData = await res.json().catch(() => ({}));
                // Recover from 404 on Update by creating specific ID? 
                // Alternatively, just alert user.
                if (res.status === 404 && !isNew) {
                    error("Workflow not found on server. Try creating as new.");
                    // Optionally could retry as POST here if we want to "Upsert"
                } else {
                    error(`Save failed: ${errData.error || res.statusText}`);
                }
            }
        } catch (e) {
            console.error(e);
            error("Error saving workflow");
        }
    };

    // Template & Wizard State
    const [showTemplates, setShowTemplates] = useState(false);
    const [showSaveTemplate, setShowSaveTemplate] = useState(false);
    const [showWizard, setShowWizard] = useState(false);

    // Initial load wizard logic could go here if route param exists

    const handleApplyTemplate = (template: Template) => {
        // Load template graph into editor
        // IMPORTANT: Strip ID so it saves as a NEW workflow, not overwriting the template source ID
        const newGraph = {
            ...template.graph,
            workflow_uuid: '', // Reset ID for new creation
            name: `${template.graph.name} (Copy)`
        };
        setGraph(newGraph);
        // Also need to reset rfInstance nodes/edges via transform?
        // Actually WorkflowCanvas handles graph prop change via useEffect.
        setShowTemplates(false);
    };

    const handleWizardComplete = (generatedGraph: BackendGraph) => {
        // Ensure generated graph is treated as new
        const newGraph = {
            ...generatedGraph,
            workflow_uuid: '',
        };
        setGraph(newGraph);
        setShowWizard(false);
    };

    return (
        <div className="h-screen flex flex-col relative">
            <header className="h-14 border-b px-4 flex items-center justify-between bg-white dark:bg-gray-800 shrink-0 z-10 shadow-sm relative">
                <div className="flex items-center gap-4">
                    <h1 className="font-bold text-lg flex items-center gap-2">
                        Workflow Builder
                    </h1>

                    <div className="h-6 w-px bg-gray-200 dark:bg-gray-700 mx-2" />

                    <div className="flex gap-2">
                        <button
                            onClick={() => setShowWizard(true)}
                            className="flex items-center gap-1.5 px-3 py-1.5 bg-purple-50 text-purple-700 hover:bg-purple-100 dark:bg-purple-900/20 dark:text-purple-300 dark:hover:bg-purple-900/30 rounded-lg text-sm font-medium transition-colors"
                        >
                            <Wand2 size={16} />
                            Wizard
                        </button>

                        <button
                            onClick={() => setShowTemplates(!showTemplates)}
                            className={`flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${showTemplates
                                ? 'bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300'
                                : 'text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700'
                                }`}
                        >
                            <LayoutTemplate size={16} />
                            Templates
                        </button>
                    </div>
                </div>

                <div className="flex gap-2">
                    <button
                        onClick={() => setShowSaveTemplate(true)}
                        className="flex items-center gap-1.5 px-3 py-1.5 text-gray-600 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 rounded-lg text-sm font-medium transition-colors"
                    >
                        <Save size={16} />
                        Save as Template
                    </button>
                    <div className="h-6 w-px bg-gray-200 dark:bg-gray-700 mx-2" />
                    <button onClick={handleSave} className="px-4 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-lg text-sm font-medium transition-colors shadow-sm">
                        Save Workflow
                    </button>
                </div>
            </header>

            <div className="flex-1 overflow-hidden relative flex">
                {/* Template Sidebar */}
                {showTemplates && (
                    <div className="relative z-20">
                        <TemplateSidebar
                            open={showTemplates}
                            onClose={() => setShowTemplates(false)}
                            onApply={handleApplyTemplate}
                        />
                    </div>
                )}

                {/* Canvas */}
                <div className="flex-1 relative">
                    <WorkflowCanvas
                        graph={graph}
                        onInit={setRfInstance}
                        onNodeClick={handleNodeClick}
                        onPaneClick={handlePaneClick}
                    />

                    {/* Cost Estimator Widget */}
                    <div className="absolute top-4 left-4 z-10 pointer-events-none">
                        <div className="pointer-events-auto">
                            {rfInstance && <CostEstimator nodes={rfInstance.getNodes()} edges={rfInstance.getEdges()} />}
                        </div>
                    </div>

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

            {/* Modals */}
            <SaveTemplateModal
                open={showSaveTemplate}
                onClose={() => setShowSaveTemplate(false)}
                currentGraph={graph} // Passes current loaded graph structure (needs to be kept in sync or re-extracted from RF instance)
            />

            <WizardMode
                open={showWizard}
                onClose={() => setShowWizard(false)}
                onComplete={handleWizardComplete}
            />
        </div>
    );
};

