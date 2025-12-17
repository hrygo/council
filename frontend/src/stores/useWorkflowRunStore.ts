import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import { enableMapSet } from 'immer';
import type { Node, Edge } from '@xyflow/react';
import type { RuntimeNode, RunControlState, ControlAction } from '../types/workflow-run';
import type { NodeStatus } from '../types/session';

enableMapSet();

interface WorkflowRunState {
    // === State ===

    /**
     * React Flow 节点 (运行时增强)
     */
    nodes: Node<RuntimeNode>[];

    /**
     * React Flow 边
     */
    edges: Edge[];

    /**
     * 当前高亮的节点 ID 列表
     */
    activeNodeIds: Set<string>;

    /**
     * 执行状态
     */
    executionStatus: 'idle' | 'running' | 'paused' | 'completed' | 'failed';

    /**
     * 累计统计
     */
    stats: {
        totalNodes: number;
        completedNodes: number;
        failedNodes: number;
        totalTokens: number;
        totalCostUsd: number;
        elapsedTimeMs: number;
    };

    /**
     * 计时器 Ref (不作为 state 存储，但保留在 store 闭包中或作为 hidden property? Zustand 不建议存 Ref 在 state。这里用 module level variable 或者 custom store creator 比较好。为了简单，我们将 timerRef 放在 actions 闭包外或者使用 useRef 在组件侧。
     * 但 spec 建议在 store 中。我们可以放在 module scope 或者作为 non-reactive state?)
     * Spec: private timerRef ... 
     * 我们可以在 store 外部定义 let timerRef.
     */

    // === Actions ===

    loadWorkflow: (nodes: Node[], edges: Edge[]) => void;
    clearWorkflow: () => void;
    updateNodeStatus: (nodeId: string, status: NodeStatus, error?: string) => void;
    setActiveNodes: (nodeIds: string[]) => void;
    addActiveNode: (nodeId: string) => void;
    removeActiveNode: (nodeId: string) => void;
    updateNodeTokenUsage: (nodeId: string, usage: NonNullable<RuntimeNode['tokenUsage']>) => void;
    setExecutionStatus: (status: WorkflowRunState['executionStatus']) => void;
    sendControl: (sessionId: string, action: ControlAction) => Promise<void>;
    startTimer: () => void;
    stopTimer: () => void;
}

// Module-level timer reference (not persisted in state)
let timerRef: ReturnType<typeof setInterval> | null = null;

export const useWorkflowRunStore = create<WorkflowRunState>()(
    subscribeWithSelector(
        immer((set, get) => ({
            nodes: [],
            edges: [],
            activeNodeIds: new Set(),
            executionStatus: 'idle',
            stats: {
                totalNodes: 0,
                completedNodes: 0,
                failedNodes: 0,
                totalTokens: 0,
                totalCostUsd: 0,
                elapsedTimeMs: 0,
            },

            loadWorkflow: (nodes, edges) => {
                set((state) => {

                    state.nodes = nodes.map(n => ({
                        ...n,
                        data: {
                            ...n.data,
                            status: 'pending',
                            id: n.id,
                            type: n.type || 'default',
                            label: n.data?.label || n.id,
                        }
                    }));
                    state.edges = edges;
                    state.stats.totalNodes = nodes.length;
                });
            },

            clearWorkflow: () => {
                set((state) => {
                    state.nodes = [];
                    state.edges = [];
                    state.activeNodeIds = new Set();
                    state.executionStatus = 'idle';
                    state.stats = {
                        totalNodes: 0,
                        completedNodes: 0,
                        failedNodes: 0,
                        totalTokens: 0,
                        totalCostUsd: 0,
                        elapsedTimeMs: 0,
                    };
                });
                if (timerRef) {
                    clearInterval(timerRef);
                    timerRef = null;
                }
            },

            updateNodeStatus: (nodeId, status, error) => {
                set((state) => {
                    const node = state.nodes.find(n => n.id === nodeId);
                    if (node) {
                        node.data.status = status;
                        if (error) node.data.error = error;

                        if (status === 'completed') state.stats.completedNodes++;
                        if (status === 'failed') state.stats.failedNodes++;
                    }
                });
            },

            setActiveNodes: (nodeIds) => {
                set((state) => {
                    state.activeNodeIds = new Set(nodeIds);
                });
            },

            addActiveNode: (nodeId) => {
                set((state) => {
                    state.activeNodeIds.add(nodeId);
                });
            },

            removeActiveNode: (nodeId) => {
                set((state) => {
                    state.activeNodeIds.delete(nodeId);
                });
            },

            updateNodeTokenUsage: (nodeId, usage) => {
                set((state) => {
                    const node = state.nodes.find(n => n.id === nodeId);
                    if (node) {
                        node.data.tokenUsage = usage;
                    }
                    state.stats.totalTokens += usage.input + usage.output;
                    state.stats.totalCostUsd += usage.cost;
                });
            },

            setExecutionStatus: (status) => {
                set((state) => {
                    state.executionStatus = status;
                });
            },

            sendControl: async (sessionId, action) => {
                const { setExecutionStatus } = get();
                try {
                    const response = await fetch(`/api/v1/sessions/${sessionId}/control`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({ action }),
                    });

                    if (!response.ok) {
                        throw new Error(`Control action failed: ${response.statusText}`);
                    }

                    switch (action) {
                        case 'pause':
                            setExecutionStatus('paused');
                            break;
                        case 'resume':
                            setExecutionStatus('running');
                            break;
                        case 'stop':
                            setExecutionStatus('failed');
                            break;
                    }
                } catch (error) {
                    console.error('Failed to send control command:', error);
                    throw error;
                }
            },

            startTimer: () => {
                const { stopTimer } = get();
                stopTimer();

                const startTime = Date.now();
                // Reset or continue? Usually continue if resuming?
                // Spec implies resetting stats.elapsedTimeMs: 0, which implies a fresh start.
                // But if resuming, we might want to add to previous.
                // For now, follow spec: reset to 0 (maybe startTimer is only for fresh run).
                // Wait, if resuming, we shouldn't reset.
                // Let's assume startTimer is called on fresh start.
                // If resuming, we need another logic or just rely on elapsedTime adding up.
                // Spec 3.4 says: set({ stats: { ...get().stats, elapsedTimeMs: 0 } });
                // This implies reset. I will follow spec for now.

                set((state) => { state.stats.elapsedTimeMs = 0; });

                // Update elapsed time periodically
                timerRef = setInterval(() => {
                    set((state) => {
                        state.stats.elapsedTimeMs = Date.now() - startTime;
                    });
                }, 100);
            },

            stopTimer: () => {
                if (timerRef) {
                    clearInterval(timerRef);
                    timerRef = null;
                }
            },

        }))
    )
);

// Derived state getter
export const getControlState = (status: WorkflowRunState['executionStatus']): RunControlState => ({
    canPause: status === 'running',
    canResume: status === 'paused',
    canStop: status === 'running' || status === 'paused',
});
