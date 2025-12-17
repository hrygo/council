import { useState } from 'react';
import type { BackendGraph } from '../utils/graphUtils';
import type { Template } from '../types/template';

interface GenerateWorkflowResponse {
    graph: BackendGraph;
    similar_templates?: Template[];
    confidence?: number;
}

export function useGenerateWorkflow() {
    const [isGenerating, setIsGenerating] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const generate = async (prompt: string): Promise<GenerateWorkflowResponse | null> => {
        setIsGenerating(true);
        setError(null);
        try {
            const res = await fetch('/api/v1/workflows/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ prompt })
            });

            if (!res.ok) {
                throw new Error('Generation failed');
            }

            const data = await res.json();

            // Normalize response if API structure varies
            return {
                graph: data.graph || data,
                similar_templates: data.similar_templates || [], // Mock or Real
                confidence: data.confidence || 0.85
            };
        } catch (e: unknown) {
            console.error(e);
            setError(e instanceof Error ? e.message : 'Unknown error');
            return null;
        } finally {
            setIsGenerating(false);
        }
    };

    return { generate, isGenerating, error };
}
