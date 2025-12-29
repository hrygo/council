/**
 * 节点执行状态枚举
 * 对应后端 NodeStatus
 */
export type NodeStatus = 'pending' | 'running' | 'completed' | 'failed';

/**
 * 会话整体状态
 * 对应后端 SessionStatus
 */
export type SessionStatus =
  | 'idle'       // 未开始
  | 'running'    // 执行中
  | 'paused'     // 已暂停
  | 'completed'  // 已完成
  | 'failed'     // 执行失败
  | 'cancelled'; // 用户取消

/**
 * 消息角色
 */
export type MessageRole = 'user' | 'agent' | 'system';

/**
 * 单条消息
 */
export interface Message {
  message_uuid: string;
  node_id: string;          // 所属节点 ID
  agent_uuid?: string;      // Agent UUID (如有)
  agentName?: string;      // Agent 显示名
  agentAvatar?: string;    // Agent 头像
  role: MessageRole;
  content: string;
  isStreaming: boolean;    // 是否正在流式输出
  timestamp: Date;
  tokenUsage?: {
    inputTokens: number;
    outputTokens: number;
    estimatedCostUsd: number;
  };
}

/**
 * 消息组 (按节点分组)
 */
export interface MessageGroup {
  node_id: string;
  nodeName: string;
  nodeType: 'start' | 'agent' | 'parallel' | 'sequence' | 'vote' | 'loop' | 'fact_check' | 'human_review' | 'end';
  isParallel: boolean;     // 是否为并行组
  messages: Message[];
  status: NodeStatus;
}

/**
 * 节点状态快照
 */
export interface NodeStateSnapshot {
  node_id: string;
  name?: string;         // 节点显示名称
  type?: string;         // 节点类型
  status: NodeStatus;
  startedAt?: Date;
  completedAt?: Date;
  tokenUsage?: {
    inputTokens: number;
    outputTokens: number;
  };
}

/**
 * 会话状态
 */
export interface WorkflowSession {
  session_uuid: string;
  workflow_uuid: string;
  group_uuid: string;
  status: SessionStatus;
  startedAt?: Date;
  completedAt?: Date;

  // 工作流图
  nodes: Map<string, NodeStateSnapshot>;

  // 当前高亮节点 (可能多个，如并行执行)
  active_node_ids: string[];

  // 累计统计
  totalTokens: number;
  totalCostUsd: number;
}
