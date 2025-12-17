import { type NodeProps } from '@xyflow/react';
import { BaseNode } from './BaseNode';
import {
    Vote,
    RefreshCw,
    Search,
    UserCheck,
    Bot,
    Play,
    Square
} from 'lucide-react';
import type {
    AgentNodeData,
    VoteNodeData,
    LoopNodeData,
    FactCheckNodeData,
    HumanReviewNodeData
} from '../../../types/workflow';

// Helper to cast data safely
function getData<T>(data: unknown): T {
    return (data || {}) as T;
}

export const AgentNode = (props: NodeProps) => {
    const data = getData<AgentNodeData>(props.data);
    return (
        <BaseNode
            label={data.label || 'Agent'}
            icon={Bot}
            selected={props.selected}
            headerColor="bg-blue-50"
        >
            <div className="space-y-1">
                <div className="font-medium text-gray-900">Agent Task</div>
                <div className="text-gray-500 truncate max-w-[120px]">
                    Model Execution
                </div>
            </div>
        </BaseNode>
    );
};

export const VoteNode = (props: NodeProps) => {
    const data = getData<VoteNodeData>(props.data);
    return (
        <BaseNode
            label={data.label || 'Vote'}
            icon={Vote}
            selected={props.selected}
            headerColor="bg-orange-50"
        >
            <div className="text-center">
                Threshold: <span className="font-semibold text-orange-600">{data.threshold ?? 0.5}</span>
            </div>
        </BaseNode>
    );
};

export const LoopNode = (props: NodeProps) => {
    const data = getData<LoopNodeData>(props.data);
    return (
        <BaseNode
            label={data.label || 'Loop'}
            icon={RefreshCw}
            selected={props.selected}
            headerColor="bg-yellow-50"
        >
            <div className="flex justify-between items-center gap-2">
                <span>Max: {data.max_rounds || 3}</span>
                <span className="text-[10px] px-1 bg-yellow-100 rounded text-yellow-700">{data.exit_condition || 'Condition'}</span>
            </div>
        </BaseNode>
    );
};

export const FactCheckNode = (props: NodeProps) => {
    const data = getData<FactCheckNodeData>(props.data);
    return (
        <BaseNode
            label={data.label || 'Fact Check'}
            icon={Search}
            selected={props.selected}
            headerColor="bg-cyan-50"
        >
            <div>
                Source: <span className="font-medium">{data.search_sources?.length || 'Auto'}</span>
            </div>
            <div>
                Strictness: <span className="font-medium">{data.verify_threshold ?? 0.8}</span>
            </div>
        </BaseNode>
    );
};

export const HumanReviewNode = (props: NodeProps) => {
    const data = getData<HumanReviewNodeData>(props.data);
    return (
        <BaseNode
            label={data.label || 'Human Review'}
            icon={UserCheck}
            selected={props.selected}
            headerColor="bg-purple-50"
        >
            <div className="flex items-center gap-2 text-purple-700">
                <span>Wait {data.timeout_minutes || 60}m</span>
                {data.allow_skip && <span className="text-[10px] border border-purple-200 px-1 rounded">Skip</span>}
            </div>
        </BaseNode>
    );
};

export const StartNode = (props: NodeProps) => {
    return (
        <BaseNode
            label="Start"
            icon={Play}
            selected={props.selected}
            headerColor="bg-green-100"
            handles={['bottom']}
        >
            <div className="text-center text-green-700 font-medium">Entry Point</div>
        </BaseNode>
    );
};

export const EndNode = (props: NodeProps) => {
    return (
        <BaseNode
            label="End"
            icon={Square}
            selected={props.selected}
            headerColor="bg-red-50"
            handles={['top']}
        >
            <div className="text-center text-red-700 font-medium">Completion</div>
        </BaseNode>
    );
};
