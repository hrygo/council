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
    const routeMessage = useCallback((msg: WSMessage) => {
        const sessionStore = useSessionStore.getState();
        const workflowStore = useWorkflowRunStore.getState();

        switch (msg.event) {
            case 'token_stream': {
                const data = msg.data as TokenStreamData;
                sessionStore.appendMessage({
                    node_id: data.node_id,
                    agent_uuid: data.agent_id,
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
                    // Auto-update session status to running when first node starts
                    sessionStore.updateSessionStatus('running');
                    // Start timer if not already running
                    if (workflowStore.executionStatus !== 'running') {
                        workflowStore.setExecutionStatus('running');
                        workflowStore.startTimer();
                    }
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
                workflowStore.stopTimer();
                sessionStore.updateSessionStatus('paused');
                break;

            case 'execution:completed':
                workflowStore.setExecutionStatus('completed');
                workflowStore.stopTimer();
                sessionStore.updateSessionStatus('completed');
                break;

            case 'human_interaction_required': {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                const data = msg.data as any;
                workflowStore.setHumanReview({
                    session_uuid: 'current', // Logic to get current session ID needed, or passed in broadcast
                    node_id: msg.node_id || data.node_id, // Ensure protocol consistency
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

            case 'tool_execution': {
                // eslint-disable-next-line @typescript-eslint/no-explicit-any
                const data = msg.data as any;
                // data: { node_id, tool, input, output }
                // We want to show "Executing tool: X"
                // Append as a system-like message or agent thought?
                // Let's format it as a thought block for now since we don't have a distinct 'tool' role UI.
                // Or better, just append to content if streaming.

                const toolInfo = `\n\n> 🛠 **Executing Tool**: \`${data.tool}\`\n\n`;

                sessionStore.appendMessage({
                    node_id: data.node_id,
                    // If we don't have agent_id in data, we might need it. 
                    // AgentProcessor.Process sends: map[string]interface{ "node_id": a.NodeID, "tool": ..., "input": ..., "output": ... }
                    // Wait, Step 201 code shows: "node_id": a.NodeID, "tool": toolName, ...
                    // It does NOT send agent_id explicitly in Data, but 'token_stream' DOES.
                    // Ideally we need agent_id to route to the correct bubble if parallel.

                    // Let's check AgentProcessor in Step 168.
                    // map[string]interface{}{ "node_id": a.NodeID, "tool": toolName, ... }
                    // It is missing "agent_id". 
                    // However, 'token_stream' uses "agent_id".
                    // If we don't have agent_id, we can't find the exact parallel column easily if multiple agents share a node (unlikely in this design).
                    // In this design, 1 node = 1 agent usually.
                    // But `appendMessage` uses `agent_uuid` to find `lastMsg`.
                    // If we assume node_id <-> agent_id 1:1, we might get away with it, OR we fix backend to send agent_id.

                    // Let's fix backend? No, I want frontend fix first.
                    // Use node_id as agent_id fallback?
                    agent_uuid: data.agent_id || data.node_id,
                    role: 'agent',
                    content: toolInfo,
                    isStreaming: true,
                    isChunk: true,
                });
                break;
            }
        }
    }, []);

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
