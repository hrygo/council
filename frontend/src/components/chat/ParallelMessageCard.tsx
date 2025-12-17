import type { FC } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import { Bot, User } from 'lucide-react';
import type { Message } from '../../types/session';

interface ParallelMessageCardProps {
    message: Message;
    index: number;
    accentColor: string;
}

export const ParallelMessageCard: FC<ParallelMessageCardProps> = ({
    message, accentColor,
}) => (
    <div className={`border rounded-lg border-t-4 bg-white shadow-sm ${accentColor}`}>
        <div className="p-3 border-b flex items-center gap-2">
            <div className="flex-shrink-0 w-6 h-6 rounded-full bg-gray-100 flex items-center justify-center overflow-hidden">
                {message.role === 'agent' ? (
                    message.agentAvatar ? (
                        <img src={message.agentAvatar} alt={message.agentName} className="w-full h-full object-cover" />
                    ) : (
                        <Bot size={14} className="text-gray-600" />
                    )
                ) : (
                    <User size={14} className="text-gray-600" />
                )}
            </div>
            <span className="font-medium text-sm truncate max-w-[120px]" title={message.agentName}>
                {message.agentName || 'Agent'}
            </span>
            {message.isStreaming && (
                <span className="animate-spin h-3 w-3 rounded-full border-2 border-blue-500 border-t-transparent ml-auto" />
            )}
        </div>
        <div className="p-3 prose prose-xs max-h-[400px] overflow-y-auto dark:prose-invert">
            <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeHighlight]}
            >
                {message.content}
            </ReactMarkdown>
            {message.isStreaming && (
                <span className="inline-block w-1.5 h-3 bg-gray-400 animate-blink ml-1 align-middle" />
            )}
        </div>
        <div className="px-3 py-2 bg-gray-50 text-xs text-gray-500 rounded-b-lg border-t">
            {message.tokenUsage ? (
                <span>💰 ${message.tokenUsage.estimatedCostUsd.toFixed(4)}</span>
            ) : (
                <span>-</span>
            )}
        </div>
    </div>
);
