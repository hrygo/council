import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { Group, CreateGroupInput } from '../types/group';

const API_BASE = '/api/v1/groups';

async function fetchGroups(): Promise<Group[]> {
    const res = await fetch(API_BASE);
    if (!res.ok) {
        throw new Error('Failed to fetch groups');
    }
    return res.json();
}

async function fetchGroup(group_uuid: string): Promise<Group> {
    const res = await fetch(`${API_BASE}/${group_uuid}`);
    if (!res.ok) {
        throw new Error('Failed to fetch group');
    }
    return res.json();
}

async function createGroup(data: CreateGroupInput): Promise<Group> {
    const res = await fetch(API_BASE, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to create group');
    }
    return res.json();
}

async function updateGroup(group: Group): Promise<Group> {
    const res = await fetch(`${API_BASE}/${group.group_uuid}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(group),
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to update group');
    }
    return res.json();
}

async function deleteGroup(group_uuid: string): Promise<void> {
    const res = await fetch(`${API_BASE}/${group_uuid}`, {
        method: 'DELETE',
    });
    if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || 'Failed to delete group');
    }
}

export function useGroups() {
    return useQuery({
        queryKey: ['groups'],
        queryFn: fetchGroups,
    });
}

export function useGroup(group_uuid: string) {
    return useQuery({
        queryKey: ['groups', group_uuid],
        queryFn: () => fetchGroup(group_uuid),
        enabled: !!group_uuid,
    });
}

export function useCreateGroup() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: createGroup,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['groups'] });
        },
    });
}

export function useUpdateGroup() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: updateGroup,
        onSuccess: (data) => {
            queryClient.invalidateQueries({ queryKey: ['groups'] });
            queryClient.invalidateQueries({ queryKey: ['groups', data.group_uuid] });
        },
    });
}

export function useDeleteGroup() {
    const queryClient = useQueryClient();
    return useMutation({
        mutationFn: deleteGroup,
        onSuccess: () => {
            queryClient.invalidateQueries({ queryKey: ['groups'] });
        },
    });
}
