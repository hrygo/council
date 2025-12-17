import type { FC } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';

interface MessageBubbleProps {
    content: string;
    isStreaming: boolean;
    role: 'user' | 'agent' | 'system';
}

export const MessageBubble: FC<MessageBubbleProps> = ({ content, isStreaming, role }) => {
    return (
        <div
            className={`
        p-3 rounded-2xl text-sm
        ${role === 'user'
                    ? "bg-blue-600 text-white rounded-br-none ml-auto max-w-[80%]"
                    : "bg-gray-50 border border-gray-100 text-gray-800 rounded-bl-none"}
        ${isStreaming ? "animate-pulse" : ""}
      `}
        >
            <div className="prose prose-sm max-w-none dark:prose-invert">
                <ReactMarkdown
                    remarkPlugins={[remarkGfm]}
                    rehypePlugins={[rehypeHighlight]}
                >
                    {content}
                </ReactMarkdown>

                {/* 流式输入光标 */}
                {isStreaming && (
                    <span className="inline-block w-2 h-4 bg-gray-400 animate-blink ml-1 align-middle" />
                )}
            </div>
        </div>
    );
};
