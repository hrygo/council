import React, { useState, useEffect } from 'react';
import { DollarSign, PieChart, Loader2 } from 'lucide-react';
import type { Node, Edge } from '@xyflow/react';

interface CostEstimatorProps {
    nodes: Node[];
    edges: Edge[];
}

interface CostEstimate {
    total_cost_usd: number;
    total_tokens: number;
    agent_breakdown: Record<string, number>;
}

export const CostEstimator: React.FC<CostEstimatorProps> = ({ nodes, edges }) => {
    const [estimate, setEstimate] = useState<CostEstimate | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchEstimate = React.useCallback(async () => {
        if (nodes.length === 0) return;

        setLoading(true);
        setError(null);
        try {
            // Transform react-flow nodes to GraphDefinition structure expected by backend
            // Note: Backend expects simplified Node struct. We need to map it.
            // But wait, the backend `GraphDefinition` expects `Nodes map[string]*Node`.
            // The frontend has React Flow nodes.
            // We need to mirror the `WorkflowEditor` save logic.

            const graphNodes: Record<string, unknown> = {};
            let startNodeId = '';

            nodes.forEach(n => {
                if (n.type === 'start') startNodeId = n.id;

                // Collect NextIDs from edges
                const nextIds = edges
                    .filter(e => e.source === n.id)
                    .map(e => e.target);

                graphNodes[n.id] = {
                    id: n.id,
                    type: n.type || 'default',
                    name: n.data.label || n.id,
                    next_ids: nextIds,
                    properties: n.data, // Pass all data as properties
                };
            });

            const payload = {
                id: 'draft',
                name: 'draft',
                description: 'draft',
                start_node_id: startNodeId,
                nodes: graphNodes,
            };

            const response = await fetch('/api/v1/workflows/estimate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload),
            });

            if (!response.ok) throw new Error('Failed to fetch estimate');
            const data = await response.json();
            setEstimate(data);
        } catch (err) {
            console.error(err);
            setError('Estimation failed');
        } finally {
            setLoading(false);
        }
    }, [nodes, edges]);

    // Debounce estimation to avoid spamming while dragging
    useEffect(() => {
        const timer = setTimeout(() => {
            fetchEstimate();
        }, 1000);
        return () => clearTimeout(timer);
    }, [fetchEstimate]); // Re-estimate when nodes change

    if (!estimate && !loading) return null;

    return (
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-4 w-64 text-sm">
            <div className="flex items-center justify-between mb-3">
                <h3 className="font-semibold text-gray-700 dark:text-gray-200 flex items-center gap-2">
                    <DollarSign className="w-4 h-4" />
                    Cost Estimate
                </h3>
                {loading && <Loader2 className="w-3 h-3 animate-spin text-gray-400" />}
            </div>

            {error ? (
                <div className="text-red-500 text-xs">{error}</div>
            ) : estimate ? (
                <div className="space-y-3">
                    <div className="grid grid-cols-2 gap-2">
                        <div className="bg-green-50 dark:bg-green-900/20 p-2 rounded">
                            <div className="text-xs text-green-600 dark:text-green-400">Total Cost</div>
                            <div className="font-bold text-green-700 dark:text-green-300">
                                ${estimate.total_cost_usd.toFixed(4)}
                            </div>
                        </div>
                        <div className="bg-blue-50 dark:bg-blue-900/20 p-2 rounded">
                            <div className="text-xs text-blue-600 dark:text-blue-400">Tokens</div>
                            <div className="font-bold text-blue-700 dark:text-blue-300">
                                {(estimate.total_tokens / 1000).toFixed(1)}k
                            </div>
                        </div>
                    </div>

                    {Object.keys(estimate.agent_breakdown).length > 0 && (
                        <div>
                            <div className="text-xs font-medium text-gray-500 mb-1 flex items-center gap-1">
                                <PieChart className="w-3 h-3" /> Breakdown
                            </div>
                            <ul className="space-y-1">
                                {Object.entries(estimate.agent_breakdown).map(([agent, cost]) => (
                                    <li key={agent} className="flex justify-between text-xs">
                                        <span className="text-gray-600 truncate max-w-[100px]">{agent}</span>
                                        <span className="text-gray-800 font-mono">${cost.toFixed(4)}</span>
                                    </li>
                                ))}
                            </ul>
                        </div>
                    )}
                </div>
            ) : null}
        </div>
    );
};
