import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useWebSocketRouter } from '../useWebSocketRouter';
import { useConnectStore } from '../../stores/useConnectStore';
import { useSessionStore } from '../../stores/useSessionStore';
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore'; // Adjust path
import type { WSMessage } from '../../types/websocket';

describe('useWebSocketRouter', () => {
    // ... setup code logic is likely fine in beforeEach but let's just make sure imports are clean.
    // The previous error was at line 24 for 'any'.
    // Let's re-read the file to be safe or just fix the standard header and the specific line if I can target it.
    // I'll replace the header first.

    beforeEach(() => {
        useSessionStore.getState().clearSession();
        useWorkflowRunStore.getState().clearWorkflow();
        useConnectStore.setState({ _lastMessage: null });

        // Mock init session to have a valid session to update
        useSessionStore.getState().initSession({
            sessionId: 'sess-1',
            workflowId: 'wf-1',
            groupId: 'g-1',
            nodes: [{ id: 'node-1', name: 'Node 1', type: 'agent' }]
        });
        // Mock workflow load
        useWorkflowRunStore.getState().loadWorkflow(
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            [{ id: 'node-1', data: {} } as any], []
        );
    });

    it('should route token_stream to sessionStore', async () => {
        const { unmount } = renderHook(() => useWebSocketRouter());

        const msg: WSMessage = {
            event: 'token_stream',
            data: {
                node_id: 'node-1',
                agent_id: 'agent-1',
                chunk: 'Hello',
                is_thinking: false
            }
        };

        act(() => {
            useConnectStore.setState({ _lastMessage: msg });
        });

        const groups = useSessionStore.getState().messageGroups;
        expect(groups).toHaveLength(1);
        expect(groups[0].messages[0].content).toBe('Hello');
        unmount();
    });

    it('should route node_state_change to both stores', () => {
        const { unmount } = renderHook(() => useWebSocketRouter());

        const msg: WSMessage = {
            event: 'node_state_change',
            data: {
                node_id: 'node-1',
                status: 'running'
            }
        };

        act(() => {
            useConnectStore.setState({ _lastMessage: msg });
        });

        // Check Workflow Store
        expect(useWorkflowRunStore.getState().nodes[0].data.status).toBe('running');
        expect(useWorkflowRunStore.getState().activeNodeIds.has('node-1')).toBe(true);
        expect(useSessionStore.getState().currentSession?.nodes.get('node-1')?.status).toBe('running');
        unmount();
    });

    it('should route node_state_change for completion', () => {
        renderHook(() => useWebSocketRouter());
        const msg: WSMessage = {
            event: 'node_state_change',
            data: { node_id: 'node-1', status: 'completed' }
        };

        act(() => { useConnectStore.setState({ _lastMessage: msg }); });

        expect(useWorkflowRunStore.getState().nodes[0].data.status).toBe('completed');
        expect(useWorkflowRunStore.getState().activeNodeIds.has('node-1')).toBe(false);
    });

    it('should route parallel_start', () => {
        renderHook(() => useWebSocketRouter());
        const msg: WSMessage = {
            event: 'node:parallel_start',
            data: { node_id: 'node-root', branches: ['br-1', 'br-2'] }
        };

        act(() => { useConnectStore.setState({ _lastMessage: msg }); });

        expect(useWorkflowRunStore.getState().activeNodeIds.size).toBe(2);
    });

    it('should route token_usage', () => {
        renderHook(() => useWebSocketRouter());
        const msg: WSMessage = {
            event: 'token_usage',
            data: {
                node_id: 'node-1',
                agent_id: 'a-1',
                input_tokens: 10,
                output_tokens: 20,
                estimated_cost_usd: 0.001
            }
        };

        act(() => { useConnectStore.setState({ _lastMessage: msg }); });

        expect(useWorkflowRunStore.getState().stats.totalTokens).toBe(30);
    });

    it('should handle execution:completed', () => {
        renderHook(() => useWebSocketRouter());
        act(() => {
            useConnectStore.setState({ _lastMessage: { event: 'execution:completed' } as any });
        });
        expect(useWorkflowRunStore.getState().executionStatus).toBe('completed');
    });

    it('should handle error events', () => {
        renderHook(() => useWebSocketRouter());
        act(() => {
            useConnectStore.setState({
                _lastMessage: {
                    event: 'error',
                    data: { node_id: 'node-1', error: 'Boom' }
                } as any
            });
        });
        expect(useWorkflowRunStore.getState().nodes[0].data.status).toBe('failed');
    });
});
