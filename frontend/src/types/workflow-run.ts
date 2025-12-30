import type { NodeStatus } from '../types/session';

/**
 * 运行时节点数据 (覆盖 React Flow Node)
 */
export interface RuntimeNode {
    node_id: string;
    type: string;
    label: string;
    status: NodeStatus;
    progress?: number;           // 0-100, 用于长时间节点
    error?: string;              // 错误信息
    tokenUsage?: {
        input: number;
        output: number;
        cost: number;
    };
    [key: string]: unknown;
}

/**
 * 控制命令
 */
export type ControlAction = 'pause' | 'resume' | 'stop';

/**
 * 运行控制状态
 */
export interface RunControlState {
    canPause: boolean;
    canResume: boolean;
    canStop: boolean;
}

export interface HumanReviewRequest {
    session_uuid: string;
    node_id: string;
    reason: string;
    timeout: number;
    payload?: {
        tool_calls?: Array<{
            name: string;
            arguments: Record<string, unknown>;
        }>;
        original_content?: string;
        [key: string]: unknown;
    };
}
