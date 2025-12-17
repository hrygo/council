import type { Node as ReactFlowNode } from '@xyflow/react';

export type NodeType = 'start' | 'end' | 'agent' | 'vote' | 'loop' | 'fact_check' | 'human_review' | 'parallel' | 'sequence';

export interface BaseNodeData {
    label: string;
    [key: string]: unknown;
}

export interface AgentNodeData extends BaseNodeData {
    agent_id: string;
    input_source?: string;
}

export interface VoteNodeData extends BaseNodeData {
    threshold: number;      // 0.5-1.0
    vote_type: 'yes_no' | 'score_1_10';
    agent_ids: string[];    // Participants
}

export interface LoopNodeData extends BaseNodeData {
    max_rounds: number;
    exit_condition: 'max_rounds' | 'consensus';
    agent_pairs: [string, string][];
}

export interface FactCheckNodeData extends BaseNodeData {
    search_sources: ('tavily' | 'serper' | 'local_kb')[];
    max_queries: number;
    verify_threshold: number;
}

export interface HumanReviewNodeData extends BaseNodeData {
    review_type: 'approve_reject' | 'edit_content';
    timeout_minutes: number;
    allow_skip: boolean;
}

export type WorkflowNodeData = AgentNodeData | VoteNodeData | LoopNodeData | FactCheckNodeData | HumanReviewNodeData | BaseNodeData;

export interface WorkflowNode extends ReactFlowNode {
    type: NodeType;
    data: WorkflowNodeData;
}
