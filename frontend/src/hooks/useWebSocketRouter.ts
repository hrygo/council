import { useEffect, useCallback } from 'react';
import { useConnectStore } from '../stores/useConnectStore';
import { useSessionStore } from '../stores/useSessionStore';
import { useWorkflowRunStore } from '../stores/useWorkflowRunStore';
import type {
    WSMessage,
    TokenStreamData,
    NodeStateChangeData,
    ParallelStartData,
    TokenUsageData
} from '../types/websocket';

export const useWebSocketRouter = () => {
    const sessionStore = useSessionStore();
    const workflowStore = useWorkflowRunStore();

    const routeMessage = useCallback((msg: WSMessage) => {
        switch (msg.event) {
            case 'token_stream': {
                const data = msg.data as TokenStreamData;
                sessionStore.appendMessage({
                    nodeId: data.node_id,
                    agentId: data.agent_id,
                    role: 'agent',
                    content: data.chunk,
                    isStreaming: true,
                    isChunk: true,
                });
                break;
            }

            case 'node_state_change': {
                const data = msg.data as NodeStateChangeData;
                workflowStore.updateNodeStatus(data.node_id, data.status);
                sessionStore.updateNodeStatus(data.node_id, data.status);

                if (data.status === 'running') {
                    workflowStore.addActiveNode(data.node_id);
                } else if (data.status === 'completed' || data.status === 'failed') {
                    workflowStore.removeActiveNode(data.node_id);
                    sessionStore.finalizeMessage(data.node_id);
                }
                break;
            }

            case 'node:parallel_start': {
                const data = msg.data as ParallelStartData;
                workflowStore.setActiveNodes(data.branches);
                sessionStore.handleParallelStart(data.node_id, data.branches);
                break;
            }

            case 'token_usage': {
                const data = msg.data as TokenUsageData;
                sessionStore.updateTokenUsage(data.node_id, data.agent_id, {
                    inputTokens: data.input_tokens,
                    outputTokens: data.output_tokens,
                    estimatedCostUsd: data.estimated_cost_usd,
                });
                workflowStore.updateNodeTokenUsage(data.node_id, {
                    input: data.input_tokens,
                    output: data.output_tokens,
                    cost: data.estimated_cost_usd,
                });
                break;
            }

            case 'execution:paused':
                workflowStore.setExecutionStatus('paused');
                sessionStore.updateSessionStatus('paused');
                break;

            case 'execution:completed':
                workflowStore.setExecutionStatus('completed');
                sessionStore.updateSessionStatus('completed');
                break;

            case 'human_interaction_required': {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                const data = msg.data as any;
                workflowStore.setHumanReview({
                    sessionId: 'current', // Logic to get current session ID needed, or passed in broadcast
                    nodeId: msg.node_id || data.node_id, // Ensure protocol consistency
                    reason: data.reason,
                    timeout: data.timeout,
                });
                break;
            }

            case 'node_resumed': {
                workflowStore.setHumanReview(null);
                break;
            }

            case 'error': {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                const data = msg.data as any;
                const nodeId = data.node_id;
                const error = data.error || 'Unknown error';

                if (nodeId) {
                    workflowStore.updateNodeStatus(nodeId, 'failed', error);
                    sessionStore.updateNodeStatus(nodeId, 'failed');
                }
                console.error('WS Error:', error);
                break;
            }
        }
    }, [sessionStore, workflowStore]);

    useEffect(() => {
        // 订阅 WebSocket 消息
        const unsubscribe = useConnectStore.subscribe(
            (state) => state._lastMessage,
            (message) => {
                if (!message) return;
                routeMessage(message);
            }
        );

        return unsubscribe;
    }, [routeMessage]);
};
