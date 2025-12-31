export type WSEventType =
    | 'token_stream'        // Token 流
    | 'node_state_change'   // 节点状态变化
    | 'node:parallel_start' // 并行节点开始
    | 'token_usage'         // Token 使用统计
    | 'execution:paused'    // 执行已暂停
    | 'execution:completed' // 执行完成
    | 'error'               // 错误
    | 'human_interaction_required' // 人工介入请求
    | 'node_resumed'        // 节点恢复执行
    | 'tool_execution';     // 工具执行

export interface WSMessage<T = unknown> {
    event: WSEventType;
    data: T;
    timestamp?: string;
    node_id?: string;
}

// 具体事件数据类型
export interface TokenStreamData {
    node_id: string;
    agent_id: string;
    chunk: string;
    is_thinking?: boolean;
}

export interface NodeStateChangeData {
    node_id: string;
    status: 'pending' | 'running' | 'completed' | 'failed';
}

export interface TokenUsageData {
    node_id: string;
    agent_id: string;
    input_tokens: number;
    output_tokens: number;
    estimated_cost_usd: number;
}

export interface ParallelStartData {
    node_id: string;
    branches: string[];
}

export interface ToolExecutionData {
    node_id: string;
    agent_id?: string;
    tool: string;
    input: string;
    output: string;
}

// 上行命令 (Client -> Server)
export type WSCommandType = 'start_session' | 'pause_session' | 'resume_session' | 'user_input';

export interface WSCommand<T = unknown> {
    cmd: WSCommandType;
    data?: T;
}
