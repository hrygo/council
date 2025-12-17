import type { FC } from 'react';
import type { MessageGroup } from '../../types/session';
import { GroupHeader } from './GroupHeader';
import { SequentialMessage } from './SequentialMessage';

interface MessageGroupCardProps {
    group: MessageGroup;
    isActive: boolean;
}

export const MessageGroupCard: FC<MessageGroupCardProps> = ({ group, isActive }) => {
    return (
        <div
            className={`
        mb-6 transition-all duration-300
        ${isActive ? "ring-2 ring-blue-500/20 bg-blue-50/30 rounded-lg p-3" : ""}
      `}
        >
            {/* 阶段标题 */}
            <GroupHeader
                nodeName={group.nodeName}
                nodeType={group.nodeType}
                status={group.status}
            />

            {/* 消息内容 */}
            <div className="mt-3 pl-4 border-l-2 border-gray-200">
                {group.isParallel ? (
                    <div className="text-gray-500 italic text-sm p-2">并行消息渲染暂未实现 (SPEC-004)</div>
                ) : (
                    group.messages.map(msg => (
                        <SequentialMessage key={msg.id} message={msg} />
                    ))
                )}
            </div>
        </div>
    );
};
