import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { Template, CreateTemplateInput } from '../types/template';

// Mock Data
const MOCK_TEMPLATES: Template[] = [
    {
        id: 'sys-1',
        name: 'Code Review Workflow',
        description: '3 Agents: Reviewer, Security, Architect. Strict process.',
        category: 'code_review',
        is_system: true,
        graph: {
            id: 'g-1',
            name: 'Code Review',
            description: '',
            start_node_id: 'start',
            nodes: {
                'start': { id: 'start', type: 'start', name: 'Start' },
                'agent-1': { id: 'agent-1', type: 'agent', name: 'Reviewer' }
            }
        }
    },
    {
        id: 'sys-2',
        name: 'Business Plan Stress Test',
        description: 'Analyze market fit, financials, and risks.',
        category: 'business_plan',
        is_system: true,
        graph: {
            id: 'g-2',
            name: 'Business Plan',
            description: '',
            start_node_id: 'start',
            nodes: {}
        }
    }
];

export function useTemplates() {
    return useQuery({
        queryKey: ['templates'],
        queryFn: async (): Promise<Template[]> => {
            // Simulate API
            await new Promise(r => setTimeout(r, 500));
            // In real app: return fetch('/api/v1/templates').then(r => r.json())
            return MOCK_TEMPLATES;
        }
    });
}

export function useCreateTemplate() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (data: CreateTemplateInput) => {
            // Simulate API
            await new Promise(r => setTimeout(r, 800));
            console.log("Creating template:", data);
            // In real app: POST /api/v1/templates
            return { id: 'new-id', ...data, is_system: false };
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['templates'] });
        }
    });
}

export function useDeleteTemplate() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: async (id: string) => {
            console.log("Deleting template:", id);
        },
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['templates'] });
        }
    });
}
