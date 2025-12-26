import type { FC } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeHighlight from 'rehype-highlight';
import rehypeKatex from 'rehype-katex';
import { FileText } from 'lucide-react';
import { useLayoutStore } from '../../stores/useLayoutStore';

interface MessageBubbleProps {
    content: string;
    isStreaming: boolean;
    role: 'user' | 'agent' | 'system';
}

export const MessageBubble: FC<MessageBubbleProps> = ({ content, isStreaming, role }) => {
    const { maximizePanel } = useLayoutStore();

    // Transform [Ref: X] to [📑 Ref: X](#doc:X) for clickable handling
    // Regex: \[Ref: (.*?)\] -> [📑 Ref: $1](#doc:$1)
    const processedContent = content.replace(/\[Ref: (.*?)\]/g, '[📑 Ref: $1](#doc:$1)');

    const components = {
        a: ({ href, children, ...props }: React.AnchorHTMLAttributes<HTMLAnchorElement>) => {
            if (href && href.startsWith('#doc:')) {
                const docId = href.replace('#doc:', '');
                return (
                    <button
                        onClick={(e) => {
                            e.preventDefault();
                            // Logic to open document reader or scroll to ref
                            // For MVP, we maximize the Document Panel (Right)
                            maximizePanel('right');
                            // In a real app, we'd also trigger the DocStore to query `docId`
                            // e.g. useDocumentStore.getState().setSearchTerm(docId);
                            console.log('Jump to document:', docId);
                        }}
                        className="inline-flex items-center gap-1 text-blue-600 dark:text-blue-400 hover:underline cursor-pointer bg-blue-50 dark:bg-blue-900/30 px-1.5 py-0.5 rounded text-xs font-medium transition-colors"
                    >
                        <FileText size={12} />
                        {children}
                    </button>
                );
            }
            return <a href={href} {...props}>{children}</a>;
        }
    };

    return (
        <div
            className={`
        p-4 rounded-2xl text-sm leading-relaxed tracking-wide
        overflow-x-auto overflow-y-auto max-h-96
        [&::-webkit-scrollbar]:w-1.5
        [&::-webkit-scrollbar-track]:bg-transparent
        [&::-webkit-scrollbar-thumb]:bg-gray-300 
        dark:[&::-webkit-scrollbar-thumb]:bg-slate-600
        [&::-webkit-scrollbar-thumb]:rounded-full
        ${role === 'user'
                    ? "bg-blue-600 text-white rounded-br-none ml-auto max-w-[80%] shadow-md"
                    : "bg-white dark:bg-slate-800 border border-gray-100 dark:border-slate-700 text-gray-800 dark:text-gray-100 rounded-bl-none shadow-sm"}
      `}
        >
            <div className={`prose prose-sm max-w-none dark:prose-invert 
                ${role === 'user' ? 'text-white prose-headings:text-white prose-p:text-white prose-a:text-white prose-code:text-white' : 'dark:text-gray-100'}
                prose-p:leading-relaxed prose-pre:my-2 prose-pre:bg-gray-800 dark:prose-pre:bg-black/30 prose-pre:rounded-lg`}>
                <ReactMarkdown
                    remarkPlugins={[remarkGfm, remarkMath]}
                    rehypePlugins={[rehypeHighlight, rehypeKatex]}
                    components={components}
                >
                    {processedContent}
                </ReactMarkdown>

                {/* 流式输入光标 */}
                {isStreaming && (
                    <span className="inline-block w-2 h-4 bg-gray-400 animate-blink ml-1 align-middle" />
                )}
            </div>
        </div>
    );
};
