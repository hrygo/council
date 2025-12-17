import type { FC } from 'react';
import type { Message } from '../../types/session';
import { MessageBubble } from './MessageBubble';
import { Bot, User } from 'lucide-react';

interface SequentialMessageProps {
    message: Message;
}

export const SequentialMessage: FC<SequentialMessageProps> = ({ message }) => {
    const isAgent = message.role === 'agent';
    const name = isAgent ? (message.agentName || 'Agent') : 'User';

    return (
        <div className={`flex gap-3 mb-4 ${!isAgent ? 'flex-row-reverse' : ''}`}>
            {/* Avatar */}
            <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gray-200 flex items-center justify-center overflow-hidden">
                {isAgent ? (
                    message.agentAvatar ? (
                        <img src={message.agentAvatar} alt={name} className="w-full h-full object-cover" />
                    ) : (
                        <Bot size={18} className="text-gray-600" />
                    )
                ) : (
                    <User size={18} className="text-gray-600" />
                )}
            </div>

            {/* 消息内容 */}
            <div className={`flex-1 min-w-0 flex flex-col ${!isAgent ? 'items-end' : 'items-start'}`}>
                {/* Name */}
                <div className="text-xs font-medium text-gray-500 mb-1">
                    {name}
                </div>

                {/* 消息气泡 */}
                <MessageBubble
                    content={message.content}
                    isStreaming={message.isStreaming ?? false}
                    role={message.role}
                />

                {/* Token 消耗 (如果有) */}
                {message.tokenUsage && (
                    <div className="mt-1 text-xs text-gray-400">
                        💰 ${message.tokenUsage.estimatedCostUsd.toFixed(4)}
                        ({message.tokenUsage.outputTokens} tokens)
                    </div>
                )}
            </div>
        </div>
    );
};
