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
                sessionId: 'sess-123',
                workflowId: 'wf-456',
                groupId: 'group-789',
                nodes: [
                    { id: 'node-1', name: 'Start', type: 'start' },
                    { id: 'node-2', name: 'Analyst', type: 'agent' },
                ],
            });

            const session = useSessionStore.getState().currentSession;
            expect(session).not.toBeNull();
            expect(session?.id).toBe('sess-123');
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
                sessionId: 'sess-1',
                workflowId: 'wf-1',
                groupId: 'g-1',
                nodes: [{ id: 'node-2', name: 'Analyst', type: 'agent' }]
            });

            // First chunk
            store.appendMessage({
                nodeId: 'node-2',
                agentId: 'agent-1',
                role: 'agent',
                content: 'Hello ',
                isStreaming: true,
                isChunk: true,
            });

            // Second chunk
            store.appendMessage({
                nodeId: 'node-2',
                agentId: 'agent-1',
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
                sessionId: 'sess-1',
                workflowId: 'wf-1',
                groupId: 'g-1',
                nodes: [{ id: 'node-2', name: 'Analyst', type: 'agent' }]
            });

            store.appendMessage({
                nodeId: 'node-2',
                agentId: 'agent-1',
                role: 'agent',
                content: 'Msg 1',
                isStreaming: false,
            });

            store.appendMessage({
                nodeId: 'node-2',
                agentId: 'agent-1',
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
                sessionId: 'sess-1',
                workflowId: 'wf-1',
                groupId: 'g-1',
                nodes: [
                    { id: 'node-2', name: 'Analyst', type: 'agent' },
                    { id: 'node-3', name: 'Reviewer', type: 'agent' }
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
                sessionId: 'sess-1',
                workflowId: 'wf-1',
                groupId: 'g-1',
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
