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
    parallelNodeMap: Map<string, string>; // Maps branch_node_id -> parent_parallel_node_id

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
        workflow_uuid: string;
        group_uuid: string;
        nodes: Array<{ node_id: string; name: string; type: string }>;
        status?: SessionStatus;
        node_statuses?: Record<string, NodeStatus>;
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
        parallelNodeMap: new Map(),
        connectionStatus: 'disconnected',

        // Actions
        initSession: ({ session_uuid, workflow_uuid, group_uuid, nodes, status, node_statuses }) => {
            const initialNodes = new Map();
            const initialGroups: MessageGroup[] = [];

            nodes.forEach(node => {
                const nodeStatus = node_statuses?.[node.node_id] || 'pending';
                initialNodes.set(node.node_id, {
                    node_id: node.node_id,
                    name: node.name,
                    type: node.type,
                    status: nodeStatus
                });

                // Proactively create MessageGroups for already running nodes to avoid blank UI
                if (nodeStatus === 'running') {
                    initialGroups.push({
                        node_id: node.node_id,
                        nodeName: node.name,
                        nodeType: node.type as MessageGroup['nodeType'],
                        isParallel: false, // Default to false, parallelNodeMap will handle specialized cases
                        messages: [],
                        status: 'running'
                    });
                }
            });

            set({
                currentSession: {
                    session_uuid,
                    workflow_uuid,
                    group_uuid,
                    status: status || 'idle',
                    nodes: initialNodes,
                    active_node_ids: initialGroups.map(g => g.node_id),
                    totalTokens: 0,
                    totalCostUsd: 0,
                    startedAt: status === 'running' ? new Date() : undefined,
                },
                messageGroups: initialGroups,
                parallelNodeMap: new Map(),
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
                let group = state.messageGroups.find(g => g.node_id === node_id);

                // Robustness: If running and no group, create one (e.g. slow start or first event)
                if (!group && status === 'running') {
                    const node = state.currentSession?.nodes.get(node_id);
                    // Skip if parallel parent mapping handles it (but here we don't have mapping check easily without more logic)
                    // Actually, if it's a parallel child, we might want to wait? 
                    // But usually parallel start handles it. 
                    // If this is a regular node, we definitely want it.
                    // If it is parallel child, adding it as separate group might be duplicate if parallel parent group exists?
                    // Strategy: Only add if NOT in parallel map?
                    // But parallelNodeMap is keyed by child ID.
                    const isParallelChild = state.parallelNodeMap.has(node_id);

                    if (!isParallelChild) {
                        group = {
                            node_id,
                            nodeName: node?.name || node_id,
                            nodeType: (node?.type as MessageGroup['nodeType']) || 'agent',
                            isParallel: false,
                            messages: [],
                            status: 'running'
                        };
                        state.messageGroups.push(group);
                    }
                }

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
                // 1.5 Check if this node is part of a parallel branch
                const parallelParentId = state.parallelNodeMap.get(msg.node_id);

                // 1. 查找对应的消息组 (Priority: Parallel Parent Group -> Direct Node Group)
                let group: MessageGroup | undefined;

                if (parallelParentId) {
                    group = state.messageGroups.find(g => g.node_id === parallelParentId);
                } else {
                    group = state.messageGroups.find(g => g.node_id === msg.node_id);
                }

                // 2. 如果不存在，创建新组 (Only for non-parallel or if parent group missing which shouldn't happen for parallel)
                if (!group) {
                    // Fail-safe: if parallel parent missing, treat as normal sequential
                    // 尝试从 nodes 中获取节点信息
                    const node = state.currentSession?.nodes.get(msg.node_id);

                    group = {
                        node_id: msg.node_id,
                        nodeName: node?.name || msg.node_id || 'Unknown Node',
                        nodeType: (node?.type as MessageGroup['nodeType']) || 'agent',
                        isParallel: false,
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
                // Check parallel parent first
                const parallelParentId = state.parallelNodeMap.get(node_id);
                const targetGroupId = parallelParentId || node_id;

                const group = state.messageGroups.find(g => g.node_id === targetGroupId);
                if (group) {
                    if (agent_uuid) {
                        const msgs = group.messages.filter((m: Message) => m.agent_uuid === agent_uuid && m.isStreaming);
                        msgs.forEach((m: Message) => { m.isStreaming = false; });
                    } else {
                        // If no agent_uuid provided, finalize ALL streaming messages in this group
                        group.messages.forEach((m: Message) => {
                            if (m.isStreaming) m.isStreaming = false;
                        });
                    }
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

                // 更新单条消息或最近一条消息的 token usage
                // Check parallel parent first
                const parallelParentId = state.parallelNodeMap.get(node_id);
                const targetGroupId = parallelParentId || node_id;

                const group = state.messageGroups.find(g => g.node_id === targetGroupId);
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
                parallelNodeMap: new Map(),
                connectionStatus: 'disconnected',
            });
        },

        setConnectionStatus: (status) => {
            set({ connectionStatus: status });
        },

        handleParallelStart: (node_id, branchIds) => {
            set(state => {
                // Register branch mappings
                branchIds.forEach(bid => {
                    state.parallelNodeMap.set(bid, node_id);
                });

                // check if group already exists to avoid duplicates
                if (state.messageGroups.find(g => g.node_id === node_id)) return;

                state.messageGroups.push({
                    node_id,
                    nodeName: 'Parallel Execution', // Could fetch proper name
                    nodeType: 'parallel',
                    isParallel: true,
                    messages: [],
                    status: 'running',
                });
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
