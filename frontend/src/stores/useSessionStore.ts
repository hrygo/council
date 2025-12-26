import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { enableMapSet } from 'immer';

enableMapSet();

import type {
    WorkflowSession,
    MessageGroup,
    SessionStatus,
    NodeStatus,
    Message
} from '../types/session';

interface SessionState {
    // === State ===

    /**
     * 当前活跃会话
     */
    currentSession: WorkflowSession | null;

    /**
     * 消息列表 (按节点分组)
     */
    messageGroups: MessageGroup[];

    /**
     * 连接状态
     */
    connectionStatus: 'disconnected' | 'connecting' | 'connected' | 'error';

    // === Actions ===

    /**
     * 初始化新会话
     */
    initSession: (params: {
        session_uuid: string;
        workflow_id: string;
        group_uuid: string;
        nodes: Array<{ node_id: string; name: string; type: string }>;
    }) => void;

    /**
     * 更新会话状态
     */
    updateSessionStatus: (status: SessionStatus) => void;

    /**
     * 更新节点状态
     */
    updateNodeStatus: (node_id: string, status: NodeStatus) => void;

    /**
     * 设置当前活跃节点
     */
    setActiveNodes: (node_ids: string[]) => void;

    /**
     * 追加/更新消息 (处理流式输出)
     * @param isStreaming - 如果为 true且isChunk为true，则追加到最后一条消息；否则创建新消息
     */
    appendMessage: (message: Omit<Message, 'message_uuid' | 'timestamp'> & { isChunk?: boolean }) => void;

    /**
     * 标记消息流式完成
     */
    finalizeMessage: (node_id: string, agent_uuid?: string) => void;

    /**
     * 更新 Token 使用量
     */
    updateTokenUsage: (node_id: string, agent_uuid: string, usage: {
        inputTokens: number;
        outputTokens: number;
        estimatedCostUsd: number;
    }) => void;

    /**
     * 清理当前会话
     */
    clearSession: () => void;

    /**
     * 设置连接状态
     */
    setConnectionStatus: (status: SessionState['connectionStatus']) => void;

    /**
     * 处理并行开始 (helper action)
     */
    handleParallelStart: (node_id: string, branchIds: string[]) => void;
}

export const useSessionStore = create<SessionState>()(
    immer((set) => ({
        // Initial State
        currentSession: null,
        messageGroups: [],
        connectionStatus: 'disconnected',

        // Actions
        initSession: ({ session_uuid, workflow_id, group_uuid, nodes }) => {
            const initialNodes = new Map();
            nodes.forEach(node => {
                initialNodes.set(node.node_id, {
                    node_id: node.node_id,
                    name: node.name,
                    type: node.type,
                    status: 'pending'
                });
            });

            set({
                currentSession: {
                    session_uuid,
                    workflow_id,
                    group_uuid,
                    status: 'idle',
                    nodes: initialNodes,
                    active_node_ids: [],
                    totalTokens: 0,
                    totalCostUsd: 0,
                },
                messageGroups: [],
                connectionStatus: 'connecting',
            });
        },

        updateSessionStatus: (status) => {
            set(state => {
                if (state.currentSession) {
                    state.currentSession.status = status;
                    if (status === 'running' && !state.currentSession.startedAt) {
                        state.currentSession.startedAt = new Date();
                    }
                    if (['completed', 'failed', 'cancelled'].includes(status)) {
                        state.currentSession.completedAt = new Date();
                    }
                }
            });
        },

        updateNodeStatus: (node_id, status) => {
            set(state => {
                if (state.currentSession) {
                    const node = state.currentSession.nodes.get(node_id);
                    if (node) {
                        node.status = status;
                        if (status === 'running') node.startedAt = new Date();
                        if (['completed', 'failed'].includes(status)) node.completedAt = new Date();
                    }
                }

                // 同时更新 MessageGroup 的状态
                const group = state.messageGroups.find(g => g.node_id === node_id);
                if (group) {
                    group.status = status;
                }
            });
        },

        setActiveNodes: (node_ids) => {
            set(state => {
                if (state.currentSession) {
                    state.currentSession.active_node_ids = node_ids;
                }
            });
        },

        appendMessage: (msg) => {
            set(state => {
                // 1. 查找对应的消息组
                let group = state.messageGroups.find(g => g.node_id === msg.node_id);

                // 2. 如果不存在，创建新组
                if (!group) {
                    // 尝试从 nodes 中获取节点信息
                    // 注意：Immer draft Map 需要特殊处理，这里简单假定如果 session 存在则 node 存在
                    // 这里的 Map 是 Immer 代理后的，可以直接 get
                    const node = state.currentSession?.nodes.get(msg.node_id);

                    group = {
                        node_id: msg.node_id,
                        nodeName: node?.name || msg.node_id || 'Unknown Node',
                        nodeType: (node?.type as MessageGroup['nodeType']) || 'agent',
                        isParallel: false, // 默认为非并行
                        messages: [],
                        status: 'running',
                    };
                    state.messageGroups.push(group);
                }

                // 3. 处理流式消息
                if (msg.isChunk && msg.isStreaming) {
                    // 查找同一 Agent 的最后一条流式消息

                    const existingMsg = group.messages.findLast(
                        (m: Message) => m.agent_uuid === msg.agent_uuid && m.isStreaming
                    );

                    if (existingMsg) {
                        // 追加内容
                        existingMsg.content += msg.content;
                        return;
                    }
                }

                // 4. 创建新消息
                group.messages.push({
                    message_uuid: crypto.randomUUID(),
                    node_id: msg.node_id,
                    agent_uuid: msg.agent_uuid,
                    agentName: msg.agentName,
                    agentAvatar: msg.agentAvatar,
                    role: msg.role,
                    content: msg.content,
                    isStreaming: msg.isStreaming,
                    timestamp: new Date(),
                });
            });
        },

        finalizeMessage: (node_id, agent_uuid) => {
            set(state => {
                const group = state.messageGroups.find(g => g.node_id === node_id);
                if (group) {

                    const msgs = group.messages.filter((m: Message) => m.agent_uuid === agent_uuid && m.isStreaming);
                    msgs.forEach((m: Message) => { m.isStreaming = false; });
                }
            });
        },

        updateTokenUsage: (node_id, agent_uuid, usage) => {
            set(state => {
                // 更新会话总成本
                if (state.currentSession) {
                    state.currentSession.totalTokens += usage.inputTokens + usage.outputTokens;
                    state.currentSession.totalCostUsd += usage.estimatedCostUsd;

                    // 更新节点成本
                    const node = state.currentSession.nodes.get(node_id);
                    if (node) {
                        if (!node.tokenUsage) {
                            node.tokenUsage = { inputTokens: 0, outputTokens: 0 };
                        }
                        node.tokenUsage.inputTokens += usage.inputTokens;
                        node.tokenUsage.outputTokens += usage.outputTokens;
                    }
                }

                // 更新单条消息或最近一条消息的 token usage (Optional, not specified in detail but good to have)
                const group = state.messageGroups.find(g => g.node_id === node_id);
                if (group) {
                    const lastMsg = group.messages.findLast((m: Message) => m.agent_uuid === agent_uuid);
                    if (lastMsg) {
                        if (!lastMsg.tokenUsage) lastMsg.tokenUsage = { inputTokens: 0, outputTokens: 0, estimatedCostUsd: 0 };
                        lastMsg.tokenUsage.inputTokens += usage.inputTokens;
                        lastMsg.tokenUsage.outputTokens += usage.outputTokens;
                        lastMsg.tokenUsage.estimatedCostUsd += usage.estimatedCostUsd;
                    }
                }
            });
        },

        clearSession: () => {
            set({
                currentSession: null,
                messageGroups: [],
                connectionStatus: 'disconnected',
            });
        },

        setConnectionStatus: (status) => {
            set({ connectionStatus: status });
        },

        handleParallelStart: (node_id, branchIds) => {
            set(state => {
                // 创建并行消息组
                state.messageGroups.push({
                    node_id,
                    nodeName: 'Parallel Execution',
                    nodeType: 'parallel', // 这里需要根据 node_id 对应类型来定，简化处理
                    isParallel: true,
                    messages: [], // 并行消息将被收集到这里 (或者是其子节点各自有 group? 这是一个设计点，Spec 似乎暗示并行组是一个容器)
                    // 根据 SPEC-001 4.2: state.messageGroups.push({ ... messages: [] })
                    // 实际上并行执行时，消息可能分散在各个并行的子节点 Group 中，或者聚合在这里。
                    // Spec 4.2 代码片段显示创建一个并行组。
                    status: 'running',
                });

                // 标记分支节点为活跃
                if (state.currentSession) {
                    state.currentSession.active_node_ids = branchIds;
                }
            });
        }

    }))
);

// === Selectors ===

/**
 * 获取当前活跃节点的消息组
 */
export const selectActiveMessageGroups = (state: SessionState): MessageGroup[] => {
    if (!state.currentSession) return [];
    // 如果当前是并行节点，可能直接显示并行组？
    // 或者显示 active_node_ids 对应的所有组
    return state.messageGroups.filter(
        g => state.currentSession!.active_node_ids.includes(g.node_id) || (g.isParallel && g.status === 'running')
    );
};

/**
 * 获取指定节点的状态
 */
export const selectNodeStatus = (node_id: string) => (state: SessionState): NodeStatus | null => {
    return state.currentSession?.nodes.get(node_id)?.status ?? null;
};

/**
 * 获取累计成本
 */
export const selectTotalCost = (state: SessionState): number => {
    return state.currentSession?.totalCostUsd ?? 0;
};
