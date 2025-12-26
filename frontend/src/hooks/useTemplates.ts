import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { Template, CreateTemplateInput } from '../types/template';

const API_BASE = '/api/v1/templates';

async function fetchTemplates(): Promise<Template[]> {
    const res = await fetch(API_BASE);
    if (!res.ok) {
        throw new Error('Failed to fetch templates');
    }
    return res.json();
}

async function createTemplate(data: CreateTemplateInput): Promise<Template> {
    const res = await fetch(API_BASE, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to create template');
    }
    return res.json();
}

async function deleteTemplate(template_uuid: string): Promise<void> {
    const res = await fetch(`${API_BASE}/${template_uuid}`, {
        method: 'DELETE',
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to delete template');
    }
}

export function useTemplates() {
    return useQuery({
        queryKey: ['templates'],
        queryFn: fetchTemplates,
    });
}

export function useCreateTemplate() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: createTemplate,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['templates'] });
        },
    });
}

export function useDeleteTemplate() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: deleteTemplate,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['templates'] });
        },
    });
}
