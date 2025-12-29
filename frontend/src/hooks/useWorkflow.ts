import { useState, useCallback } from 'react';
import { type BackendGraph } from '../utils/graphUtils';

interface UseWorkflowReturn {
    workflow: BackendGraph | null;
    loading: boolean;
    error: string | null;
    fetchWorkflow: (id: string) => Promise<void>;
    executeWorkflow: (id: string, input: string) => Promise<void>;
}

export const useWorkflow = (): UseWorkflowReturn => {
    const [workflow, setWorkflow] = useState<BackendGraph | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchWorkflow = useCallback(async (id: string) => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch(`/api/v1/workflows/${id}`);
            if (!response.ok) {
                throw new Error(`Failed to fetch workflow: ${response.statusText}`);
            }
            const data = await response.json();
            // Assuming API returns { id, name, ..., graph_definition } 
            // or maybe the graph is nested.
            // WorkflowMgmtHandler.Get returns the Workflow struct which has GraphDefinition.
            // Let's assume the API returns the GraphDefinition directly or we extract it.
            // Actually, the Workflow struct in Go has `GraphDefinition JSONB`. 
            // The API likely returns the full Workflow object.
            // We need to verify what the backend returns.
            // Based on `workflowMgmtHandler.Get`, it likely returns JSON of the workflow.
            // So `data.graph_definition` is likely what we need, which matches `BackendGraph`.

            // Adjust based on actual API response structure
            const graphDef = data.graph_definition || data;

            // Map to our BackendGraph interface
            // Ensure fields match
            setWorkflow(graphDef);
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Unknown error';
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    }, []);

    const executeWorkflow = useCallback(async (id: string, input: string) => {
        setLoading(true);
        setError(null);
        try {
            const response = await fetch('/api/v1/workflows/execute', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    workflow_uuid: id,
                    input: { initial_prompt: input }
                }),
            });

            if (!response.ok) {
                const errData = await response.json();
                throw new Error(errData.error || 'Execution failed');
            }
            // Execution started successfully
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Unknown error';
            setError(errorMessage);
            // Re-throw if caller needs to handle it
            throw err;
        } finally {
            setLoading(false);
        }
    }, []);

    return {
        workflow,
        loading,
        error,
        fetchWorkflow,
        executeWorkflow,
    };
};
