import { type FC, useState, useRef, useLayoutEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import remarkMath from 'remark-math';
import rehypeHighlight from 'rehype-highlight';
import rehypeKatex from 'rehype-katex';
import { FileText, ChevronDown, ChevronUp } from 'lucide-react';
import { useLayoutStore } from '../../stores/useLayoutStore';

interface MessageBubbleProps {
    content: string;
    isStreaming: boolean;
    role: 'user' | 'agent' | 'system';
}

export const MessageBubble: FC<MessageBubbleProps> = ({ content, isStreaming, role }) => {
    const { maximizePanel } = useLayoutStore();
    const [isExpanded, setIsExpanded] = useState(false);
    const [isOverflowing, setIsOverflowing] = useState(false);
    const contentRef = useRef<HTMLDivElement>(null);

    // Transform content [Ref: X] to [📑 Ref: X](#doc:X) for clickable handling
    const processedContent = content.replace(/\[Ref: (.*?)\]/g, '[📑 Ref: $1](#doc:$1)');

    const components = {
        a: ({ href, children, ...props }: React.AnchorHTMLAttributes<HTMLAnchorElement>) => {
            if (href && href.startsWith('#doc:')) {
                const docId = href.replace('#doc:', '');
                return (
                    <button
                        onClick={(e) => {
                            e.preventDefault();
                            maximizePanel('right');
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

    // Check for overflow logic
    useLayoutEffect(() => {
        if (contentRef.current) {
            // Check if scrollHeight is significantly larger than max height (384px = 96 * 4)
            // We use a slightly smaller threshold (390) to avoid showing button for borderline cases
            setIsOverflowing(contentRef.current.scrollHeight > 390);
        }
    }, [content, isStreaming]);

    // Force expand while streaming to show latest content
    const effectiveExpanded = isExpanded || isStreaming;
    // Only apply collapse logic to AGENT messages
    const isAgent = role !== 'user';
    const shouldCollapse = isAgent && isOverflowing && !effectiveExpanded;

    return (
        <div className={`
            relative
            ${role === 'user' ? "ml-auto max-w-[80%]" : "w-full"}
            transition-all duration-300 ease-in-out
        `}>
            <div
                ref={contentRef}
                className={`
                    p-4 rounded-2xl text-sm leading-relaxed tracking-wide
                    overflow-x-auto
                    ${shouldCollapse ? 'max-h-96 overflow-y-hidden' : ''} 
                    ${role === 'user'
                        ? "bg-blue-600 text-white rounded-br-none shadow-md"
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

            {/* Expand/Collapse Controls for Agent */}
            {isAgent && isOverflowing && !isStreaming && (
                <div className={`
                    absolute bottom-0 left-0 right-0 
                    flex justify-center items-end pb-2 rounded-2xl
                    ${!isExpanded ? 'h-24 bg-gradient-to-t from-white via-white/90 to-transparent dark:from-slate-800 dark:via-slate-800/90' : 'relative mt-2 bg-transparent'}
                `}>
                    <button
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="flex items-center gap-1.5 px-3 py-1.5 bg-gray-100 dark:bg-slate-700 hover:bg-gray-200 dark:hover:bg-slate-600 text-xs font-semibold text-gray-700 dark:text-gray-200 rounded-full shadow-sm transition-colors backdrop-blur-sm border border-gray-200 dark:border-slate-600"
                    >
                        {isExpanded ? (
                            <>
                                <ChevronUp size={14} />
                                收起 (Show Less)
                            </>
                        ) : (
                            <>
                                <ChevronDown size={14} />
                                展开阅读全文 (Show More)
                            </>
                        )}
                    </button>
                </div>
            )}
        </div>
    );
};
