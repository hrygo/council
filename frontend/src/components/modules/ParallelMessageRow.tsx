import type { FC } from 'react';
import ReactMarkdown from 'react-markdown';

interface Message {
    agentName: string;
    content: string;
}

export const ParallelMessageRow: FC<{ messages: Message[] }> = ({ messages }) => {
    return (
        <div className="flex gap-4 w-full overflow-x-auto pb-2">
            {messages.map((msg, idx) => (
                <div key={idx} className="flex-1 min-w-[300px] border rounded-lg p-3 bg-white dark:bg-gray-800 shadow-sm">
                    <div className="flex items-center gap-2 mb-2 border-b pb-1">
                        <div className="w-6 h-6 rounded-full bg-blue-500 flex items-center justify-center text-white text-xs font-bold">
                            {msg.agentName[0]}
                        </div>
                        <span className="font-semibold text-sm">{msg.agentName}</span>
                    </div>
                    <div className="text-sm prose dark:prose-invert max-w-none">
                        <ReactMarkdown>{msg.content}</ReactMarkdown>
                    </div>
                </div>
            ))}
        </div>
    );
};
