import type { FC } from 'react';
import type { NodeStatus } from '../../types/session';

interface GroupHeaderProps {
    nodeName: string;
    nodeType: string;
    status: NodeStatus;
}

const nodeTypeIcons: Record<string, string> = {
    start: '🚀',
    agent: '🤖',
    parallel: '⚡',
    sequence: '📝',
    vote: '🗳️',
    loop: '🔄',
    fact_check: '🔍',
    human_review: '👤',
    end: '🏁',
};

const statusColors: Record<NodeStatus, string> = {
    pending: 'text-gray-400',
    running: 'text-blue-500',
    completed: 'text-green-500',
    failed: 'text-red-500',
};

export const GroupHeader: FC<GroupHeaderProps> = ({ nodeName, nodeType, status }) => {
    const icon = nodeTypeIcons[nodeType] || '📍';

    return (
        <div className="flex items-center gap-2 text-sm font-medium text-gray-600">
            <span>{icon}</span>
            <span>{nodeName}</span>

            {/* 状态指示器 */}
            <span className={`ml-auto ${statusColors[status]}`}>
                {status === 'running' && (
                    <span className="inline-flex items-center gap-1">
                        <span className="animate-spin h-3 w-3 rounded-full border-2 border-current border-t-transparent" />
                        进行中
                    </span>
                )}
                {status === 'completed' && '✓ 已完成'}
                {status === 'failed' && '✕ 失败'}
            </span>
        </div>
    );
};
