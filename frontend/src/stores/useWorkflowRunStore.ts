import { create } from 'zustand';
import { subscribeWithSelector } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import { enableMapSet } from 'immer';
import type { Node, Edge } from '@xyflow/react';
import type { RuntimeNode, RunControlState, ControlAction, HumanReviewRequest } from '../types/workflow-run';
import type { NodeStatus } from '../types/session';
import type { BackendGraph } from '../utils/graphUtils';
import { transformToReactFlow } from '../utils/graphUtils';
import type { Template } from '../types/template';

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
    active_node_ids: Set<string>;

    /**
     * 执行状态
     */
    executionStatus: 'idle' | 'running' | 'paused' | 'completed' | 'failed';
    humanReview: HumanReviewRequest | null;

    /**
     * 原始图定义 (用于 Live Monitor)
     */
    graphDefinition: BackendGraph | null;

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
    setGraphFromTemplate: (template: Template) => void;
    clearWorkflow: () => void;
    updateNodeStatus: (node_id: string, status: NodeStatus, error?: string) => void;
    setActiveNodes: (node_ids: string[]) => void;
    addActiveNode: (node_id: string) => void;
    removeActiveNode: (node_id: string) => void;
    updateNodeTokenUsage: (node_id: string, usage: NonNullable<RuntimeNode['tokenUsage']>) => void;
    setExecutionStatus: (status: WorkflowRunState['executionStatus']) => void;
    sendControl: (session_uuid: string, action: ControlAction) => Promise<void>;
    setHumanReview: (request: HumanReviewRequest | null) => void;
    submitHumanReview: (req: HumanReviewRequest, action: 'approve' | 'reject' | 'modify', data?: Record<string, unknown>) => Promise<void>;
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
            active_node_ids: new Set(),
            executionStatus: 'idle',
            humanReview: null as HumanReviewRequest | null,
            graphDefinition: null,
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
                            node_id: n.id,
                            type: n.type || 'default',
                            label: (n.data?.label as string) || n.id,
                        }
                    }));
                    state.edges = edges;
                    state.stats.totalNodes = nodes.length;
                });
            },

            setGraphFromTemplate: (template) => {
                set((state) => {
                    state.graphDefinition = template.graph;
                });
                const { nodes, edges } = transformToReactFlow(template.graph);
                get().loadWorkflow(nodes, edges);
            },

            clearWorkflow: () => {
                set((state) => {
                    state.nodes = [];
                    state.edges = [];
                    state.active_node_ids = new Set();
                    state.executionStatus = 'idle';
                    state.graphDefinition = null;
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

            updateNodeStatus: (node_id, status, error) => {
                set((state) => {
                    const node = state.nodes.find(n => n.id === node_id);
                    if (node) {
                        node.data.status = status;
                        if (error) node.data.error = error;

                        if (status === 'completed') state.stats.completedNodes++;
                        if (status === 'failed') state.stats.failedNodes++;
                    }
                });
            },

            setActiveNodes: (node_ids) => {
                set((state) => {
                    state.active_node_ids = new Set(node_ids);
                });
            },

            addActiveNode: (node_id) => {
                set((state) => {
                    state.active_node_ids.add(node_id);
                });
            },

            removeActiveNode: (node_id) => {
                set((state) => {
                    state.active_node_ids.delete(node_id);
                });
            },

            updateNodeTokenUsage: (node_id, usage) => {
                set((state) => {
                    const node = state.nodes.find(n => n.id === node_id);
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

            sendControl: async (session_uuid, action) => {
                const { setExecutionStatus } = get();
                try {
                    const response = await fetch(`/api/v1/sessions/${session_uuid}/control`, {
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

            setHumanReview: (request) => {
                set((state) => {
                    state.humanReview = request;
                });
            },

            submitHumanReview: async (req, action, data) => {
                const { setHumanReview } = get();
                try {
                    const payload = {
                        node_id: req.node_id,
                        action,
                        data,
                    };
                    const response = await fetch(`/api/v1/sessions/${req.session_uuid}/review`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify(payload),
                    });

                    if (!response.ok) {
                        throw new Error(`Review submission failed: ${response.statusText}`);
                    }

                    // On success, clear the modal state
                    setHumanReview(null);
                } catch (error) {
                    console.error('Failed to submit human review:', error);
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
