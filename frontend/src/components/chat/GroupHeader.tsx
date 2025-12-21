import { type FC } from 'react';
import { useTranslation } from 'react-i18next';
import type { NodeStatus } from '../../types/session';
import { AgentAvatar } from '../common/AgentAvatar';

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
    pending: 'text-gray-400 font-normal',
    running: 'text-blue-600 font-semibold',
    completed: 'text-green-600 font-medium',
    failed: 'text-red-500 font-bold',
};

export const GroupHeader: FC<GroupHeaderProps> = ({ nodeName, nodeType, status }) => {
    const { t } = useTranslation();
    const icon = nodeTypeIcons[nodeType] || '📍';
    const isAgent = nodeType === 'agent';

    return (
        <div className="flex items-center gap-3 text-sm text-gray-700 dark:text-gray-200">
            {isAgent ? (
                <AgentAvatar name={nodeName} size="sm" />
            ) : (
                <span className="text-lg w-6 h-6 flex items-center justify-center grayscale-[0.5] opacity-80">{icon}</span>
            )}

            <div className="flex flex-col">
                <span className={`font-medium ${status === 'running' ? 'text-blue-600 dark:text-blue-400' : ''}`}>
                    {nodeName}
                </span>
                {isAgent && status === 'running' && (
                    <span className="text-[10px] text-gray-400 leading-none">{t('meeting.status.thinking')}</span>
                )}
            </div>

            {/* 状态指示器 */}
            <span className={`ml-auto text-xs ${statusColors[status]} flex items-center gap-1.5`}>
                {status === 'running' && (
                    <>
                        <span className="animate-spin h-3 w-3 rounded-full border-2 border-current border-t-transparent opacity-70" />
                        <span className="hidden sm:inline">{t('meeting.status.processing')}</span>
                    </>
                )}
                {status === 'completed' && (
                    <span className="opacity-80">{t('meeting.status.completed')}</span>
                )}
                {status === 'failed' && t('meeting.status.failed')}
                {status === 'pending' && <span className="opacity-50">{t('meeting.status.pending')}</span>}
            </span>
        </div>
    );
};
