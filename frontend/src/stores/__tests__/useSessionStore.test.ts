import { describe, it, expect, beforeEach } from 'vitest';
import { useSessionStore } from '../useSessionStore';

describe('useSessionStore', () => {
    beforeEach(() => {
        useSessionStore.getState().clearSession();
    });

    describe('initSession', () => {
        it('should initialize session with correct structure', () => {
            const { initSession } = useSessionStore.getState();

            initSession({
                session_uuid: 'sess-123',
                workflow_uuid: 'wf-456',
                group_uuid: 'group-789',
                nodes: [
                    { node_id: 'node-1', name: 'Start', type: 'start' },
                    { node_id: 'node-2', name: 'Analyst', type: 'agent' },
                ],
            });

            const session = useSessionStore.getState().currentSession;
            expect(session).not.toBeNull();
            expect(session?.session_uuid).toBe('sess-123');
            expect(session?.status).toBe('idle');
            // Map comparison in Vitest might need checking size or getting keys
            expect(session?.nodes.size).toBe(2);
            expect(session?.nodes.get('node-1')).toBeDefined();
        });
    });

    describe('appendMessage (streaming)', () => {
        it('should append chunks to existing streaming message', () => {
            const store = useSessionStore.getState();
            store.initSession({
                session_uuid: 'sess-1',
                workflow_uuid: 'wf-1',
                group_uuid: 'g-1',
                nodes: [{ node_id: 'node-2', name: 'Analyst', type: 'agent' }]
            });

            // First chunk
            store.appendMessage({
                node_id: 'node-2',
                agent_uuid: 'agent-1',
                role: 'agent',
                content: 'Hello ',
                isStreaming: true,
                isChunk: true,
            });

            // Second chunk
            store.appendMessage({
                node_id: 'node-2',
                agent_uuid: 'agent-1',
                role: 'agent',
                content: 'World!',
                isStreaming: true,
                isChunk: true,
            });

            const groups = useSessionStore.getState().messageGroups;
            expect(groups[0].messages.length).toBe(1);
            expect(groups[0].messages[0].content).toBe('Hello World!');
        });

        it('should create new message if not streaming or not chunk', () => {
            const store = useSessionStore.getState();
            store.initSession({
                session_uuid: 'sess-1',
                workflow_uuid: 'wf-1',
                group_uuid: 'g-1',
                nodes: [{ node_id: 'node-2', name: 'Analyst', type: 'agent' }]
            });

            store.appendMessage({
                node_id: 'node-2',
                agent_uuid: 'agent-1',
                role: 'agent',
                content: 'Msg 1',
                isStreaming: false,
            });

            store.appendMessage({
                node_id: 'node-2',
                agent_uuid: 'agent-1',
                role: 'agent',
                content: 'Msg 2',
                isStreaming: false,
            });

            const groups = useSessionStore.getState().messageGroups;
            expect(groups[0].messages.length).toBe(2);
        });
    });

    describe('updateTokenUsage', () => {
        it('should accumulate total cost', () => {
            const store = useSessionStore.getState();
            store.initSession({
                session_uuid: 'sess-1',
                workflow_uuid: 'wf-1',
                group_uuid: 'g-1',
                nodes: [
                    { node_id: 'node-2', name: 'Analyst', type: 'agent' },
                    { node_id: 'node-3', name: 'Reviewer', type: 'agent' }
                ]
            });

            store.updateTokenUsage('node-2', 'agent-1', {
                inputTokens: 100,
                outputTokens: 50,
                estimatedCostUsd: 0.01,
            });

            store.updateTokenUsage('node-3', 'agent-2', {
                inputTokens: 200,
                outputTokens: 100,
                estimatedCostUsd: 0.02,
            });

            const session = useSessionStore.getState().currentSession;
            expect(session?.totalCostUsd).toBeCloseTo(0.03);
            expect(session?.totalTokens).toBe(450);
        });
    });

    describe('updateSessionStatus', () => {
        it('should update status and logic timestamps', () => {
            const store = useSessionStore.getState();
            store.initSession({
                session_uuid: 'sess-1',
                workflow_uuid: 'wf-1',
                group_uuid: 'g-1',
                nodes: []
            });

            useSessionStore.getState().updateSessionStatus('running');
            expect(useSessionStore.getState().currentSession?.status).toBe('running');
            expect(useSessionStore.getState().currentSession?.startedAt).toBeDefined();

            useSessionStore.getState().updateSessionStatus('completed');
            expect(useSessionStore.getState().currentSession?.status).toBe('completed');
            expect(useSessionStore.getState().currentSession?.completedAt).toBeDefined();
        });
    });
});
