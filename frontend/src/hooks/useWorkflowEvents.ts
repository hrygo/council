import { useEffect } from 'react';
import { useConnectStore } from '../stores/useConnectStore';
import { useWorkflowRunStore } from '../stores/useWorkflowRunStore';
import type { NodeStatus } from '../types/session';

interface WebSocketMessage {
    event: string;
    data: Record<string, unknown>;
}

export const useWorkflowEvents = () => {
    const lastMessage = useConnectStore((state) => state.lastMessage);
    const workflowStore = useWorkflowRunStore();

    useEffect(() => {
        if (!lastMessage) return;

        // Type assertion or safe parsing recommended here instead of rough casting
        // Assuming lastMessage is unknown conform to WebSocketMessage structure
        const message = lastMessage as unknown as WebSocketMessage;
        const { event, data } = message;

        if (!event || !data) return;

        switch (event) {
            case 'node_state_change':
                if (typeof data.node_id === 'string' && typeof data.status === 'string') {
                    workflowStore.updateNodeStatus(data.node_id, data.status as NodeStatus);
                    if (data.status === 'running') {
                        workflowStore.addActiveNode(data.node_id);
                    } else if (data.status === 'completed' || data.status === 'failed') {
                        workflowStore.removeActiveNode(data.node_id);
                    }
                }
                break;

            case 'node:parallel_start':
                if (Array.isArray(data.branches)) {
                    // Ensure branches is string[]
                    const branches = data.branches.filter((b): b is string => typeof b === 'string');
                    workflowStore.setActiveNodes(branches);
                }
                break;

            case 'token_usage':
                if (
                    typeof data.node_id === 'string' &&
                    typeof data.input_tokens === 'number' &&
                    typeof data.output_tokens === 'number' &&
                    typeof data.estimated_cost_usd === 'number'
                ) {
                    workflowStore.updateNodeTokenUsage(data.node_id, {
                        input: data.input_tokens,
                        output: data.output_tokens,
                        cost: data.estimated_cost_usd,
                    });
                }
                break;

            case 'execution:paused':
                workflowStore.setExecutionStatus('paused');
                workflowStore.stopTimer();
                break;

            case 'execution:completed':
                workflowStore.setExecutionStatus('completed');
                workflowStore.stopTimer();
                break;

            case 'execution:failed':
                workflowStore.setExecutionStatus('failed');
                workflowStore.stopTimer();
                break;

            case 'error':
                if (typeof data.node_id === 'string' && typeof data.error === 'string') {
                    workflowStore.updateNodeStatus(data.node_id, 'failed', data.error);
                }
                break;
        }
    }, [lastMessage, workflowStore]);
};
