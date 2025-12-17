import { useEffect, useRef } from 'react';
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

    useEffect(() => {
        // 订阅 WebSocket 消息
        const unsubscribe = useConnectStore.subscribe(
            (state) => state._lastMessage,
            (message) => {
                if (!message) return;

                // 防重复处理
                // Assuming message has some ID or unique timestamp combination if simpler approach needed.
                // If message doesn't have unique ID, we might need to rely on strict sequential processing.
                // The spec suggested `${message.event}-${message.timestamp}`, assuming timestamp exists.
                // If data doesn't guarantee unique timestamp, simpler logic: just process it.
                // React effect re-runs might cause double subscription if not careful, but `subscribeWithSelector` returns unsubscribe.
                // We will skip strict deduping for now unless specific message IDs are added to protocol.

                // 路由到对应处理器
                routeMessage(message);
            }
        );

        return unsubscribe;
    }, [sessionStore, workflowStore]);

    const routeMessage = (msg: WSMessage) => {
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
                sessionStore.updateNodeStatus(data.node_id, data.status); // Session needs updates too

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
    };
};
