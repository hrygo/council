import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { Agent, CreateAgentInput } from '../types/agent';

const API_BASE = '/api/v1/agents';

async function fetchAgents(): Promise<Agent[]> {
    const res = await fetch(API_BASE);
    if (!res.ok) {
        throw new Error('Failed to fetch agents');
    }
    return res.json();
}

async function fetchAgent(id: string): Promise<Agent> {
    const res = await fetch(`${API_BASE}/${id}`);
    if (!res.ok) {
        throw new Error('Failed to fetch agent');
    }
    return res.json();
}

async function createAgent(data: CreateAgentInput): Promise<Agent> {
    const res = await fetch(API_BASE, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to create agent');
    }
    return res.json();
}

async function updateAgent(agent: Agent): Promise<Agent> {
    const res = await fetch(`${API_BASE}/${agent.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(agent),
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to update agent');
    }
    return res.json();
}

async function deleteAgent(id: string): Promise<void> {
    const res = await fetch(`${API_BASE}/${id}`, {
        method: 'DELETE',
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to delete agent');
    }
}

export function useAgents() {
    return useQuery({
        queryKey: ['agents'],
        queryFn: fetchAgents,
    });
}

export function useAgent(id: string) {
    return useQuery({
        queryKey: ['agents', id],
        queryFn: () => fetchAgent(id),
        enabled: !!id,
    });
}

export function useCreateAgent() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: createAgent,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['agents'] });
        },
    });
}

export function useUpdateAgent() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: updateAgent,
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ['agents'] });
            queryClient.invalidateQueries({ queryKey: ['agents', data.id] });
        },
    });
}

export function useDeleteAgent() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: deleteAgent,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['agents'] });
        },
    });
}
