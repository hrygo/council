import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useWorkflowRunStore, getControlState } from '../useWorkflowRunStore';


describe('useWorkflowRunStore', () => {
    beforeEach(() => {
        useWorkflowRunStore.getState().clearWorkflow();
    });

    describe('updateNodeStatus', () => {
        it('should update node status', () => {
            const { loadWorkflow } = useWorkflowRunStore.getState();

            const mockNodes = [
                { id: 'node-1', type: 'start', position: { x: 0, y: 0 }, data: {} },
            ];

            loadWorkflow(mockNodes, []);

            useWorkflowRunStore.getState().updateNodeStatus('node-1', 'running');

            const updatedNodes = useWorkflowRunStore.getState().nodes;
            expect(updatedNodes[0].data.status).toBe('running');
        });

        it('should track stats correctly', () => {
            const { loadWorkflow } = useWorkflowRunStore.getState();
            const mockNodes = [
                { id: 'node-1', type: 'start', position: { x: 0, y: 0 }, data: {} },
                { id: 'node-2', type: 'agent', position: { x: 0, y: 0 }, data: {} },
            ];

            loadWorkflow(mockNodes, []);

            useWorkflowRunStore.getState().updateNodeStatus('node-1', 'completed');
            expect(useWorkflowRunStore.getState().stats.completedNodes).toBe(1);

            useWorkflowRunStore.getState().updateNodeStatus('node-2', 'failed');
            expect(useWorkflowRunStore.getState().stats.failedNodes).toBe(1);
        });
    });

    describe('controlState', () => {
        // Since controlState is a derived getter on the exported hook or util, 
        // but here we might want to test the utility function if exported or check store behavior if mapped.
        // In our implementation we exported `getControlState` util.



        it('should derive correct control states', () => {
            expect(getControlState('running').canPause).toBe(true);
            expect(getControlState('running').canResume).toBe(false);

            expect(getControlState('paused').canPause).toBe(false);
            expect(getControlState('paused').canResume).toBe(true);

            expect(getControlState('idle').canPause).toBe(false);
        });
    });

    describe('activeNodeIds', () => {
        it('should manage active nodes for parallel execution', () => {
            useWorkflowRunStore.getState().setActiveNodes(['node-a', 'node-b']);
            expect(useWorkflowRunStore.getState().activeNodeIds.size).toBe(2);
            expect(useWorkflowRunStore.getState().activeNodeIds.has('node-a')).toBe(true);

            useWorkflowRunStore.getState().removeActiveNode('node-a');
            expect(useWorkflowRunStore.getState().activeNodeIds.has('node-a')).toBe(false);
            expect(useWorkflowRunStore.getState().activeNodeIds.has('node-b')).toBe(true);

            useWorkflowRunStore.getState().addActiveNode('node-c');
            expect(useWorkflowRunStore.getState().activeNodeIds.has('node-c')).toBe(true);
        });
    });

    describe('timers', () => {
        it('startTimer should set elapsedTime to 0 and start updating', () => {
            vi.useFakeTimers();
            const store = useWorkflowRunStore.getState();

            store.startTimer();
            expect(store.stats.elapsedTimeMs).toBe(0);

            vi.advanceTimersByTime(200);
            // Since we use setInterval 100ms, it should update.
            // However, store update happens inside interval. 
            // We need to check store state again.

            expect(useWorkflowRunStore.getState().stats.elapsedTimeMs).toBeGreaterThan(0);

            store.stopTimer();
            vi.useRealTimers();
        });
        describe('updateNodeTokenUsage', () => {
            it('should update token usage for a node and global stats', () => {
                const { loadWorkflow } = useWorkflowRunStore.getState();
                loadWorkflow([{ id: 'node-1', position: { x: 0, y: 0 }, data: {} as any }], []);

                useWorkflowRunStore.getState().updateNodeTokenUsage('node-1', {
                    input: 100,
                    output: 50,
                    cost: 0.002
                });

                const node = useWorkflowRunStore.getState().nodes[0];
                expect(node.data.tokenUsage?.input).toBe(100);
                expect(useWorkflowRunStore.getState().stats.totalTokens).toBe(150);
                expect(useWorkflowRunStore.getState().stats.totalCostUsd).toBe(0.002);
            });
        });

        describe('async actions', () => {
            beforeEach(() => {
                vi.stubGlobal('fetch', vi.fn());
            });

            it('sendControl should call API and update state', async () => {
                (fetch as any).mockResolvedValue({ ok: true });

                await useWorkflowRunStore.getState().sendControl('session-1', 'pause');

                expect(fetch).toHaveBeenCalledWith('/api/v1/sessions/session-1/control', expect.any(Object));
                expect(useWorkflowRunStore.getState().executionStatus).toBe('paused');
            });

            it('submitHumanReview should call API and clear review request', async () => {
                (fetch as any).mockResolvedValue({ ok: true });

                const request = { sessionId: 's1', nodeId: 'n1', reason: 'test', timeout: 30 };
                useWorkflowRunStore.getState().setHumanReview(request);

                await useWorkflowRunStore.getState().submitHumanReview(request, 'approve');

                expect(fetch).toHaveBeenCalledWith('/api/v1/sessions/s1/review', expect.any(Object));
                expect(useWorkflowRunStore.getState().humanReview).toBeNull();
            });
        });
    });
});
