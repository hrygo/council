import { useState, useRef, useEffect } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import rehypeHighlight from 'rehype-highlight';
import { Send, Bot, User } from 'lucide-react';
import 'highlight.js/styles/github.css'; // Or any other style you prefer
import { useConnectStore } from '../../stores/useConnectStore';

interface Message {
    role: 'user' | 'assistant';
    content: string;
}

interface ChatPanelProps {
    fullscreen?: boolean;
    onExitFullscreen?: () => void;
}

export default function ChatPanel({ fullscreen, onExitFullscreen }: ChatPanelProps) {
    const [input, setInput] = useState('');
    const [messages, setMessages] = useState<Message[]>([
        { role: 'assistant', content: '### Welcome to The Council\nI am ready to assist you with your workflows.' }
    ]);
    const messagesEndRef = useRef<HTMLDivElement>(null);
    const { sendMessage, lastMessage, connect } = useConnectStore();

    useEffect(() => {
        // Auto-connect for demo purposes (in real app, might happen on session start)
        connect('ws://localhost:8080/ws');
    }, [connect]);

    const processedMsgRef = useRef<unknown>(null);

    useEffect(() => {
        if (lastMessage && lastMessage !== processedMsgRef.current) {
            processedMsgRef.current = lastMessage;
            const msg = lastMessage as Record<string, unknown>;

            // Example handling: if lastMessage has 'content'
            if (msg.type === 'agent:speaking' || msg.content) {
                const data = msg.data as Record<string, unknown> | undefined;
                const content = (msg.content as string) || (data?.content as string);

                if (content) {
                    // eslint-disable-next-line react-hooks/set-state-in-effect
                    setMessages(prev => [...prev, { role: 'assistant', content }]);
                }
            }
        }
    }, [lastMessage]);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    };

    useEffect(scrollToBottom, [messages]);

    const handleSend = () => {
        if (!input.trim()) return;
        const newMsg: Message = { role: 'user', content: input };
        setMessages(prev => [...prev, newMsg]);
        setInput('');

        // Send to backend
        sendMessage({ type: 'user_input', content: input });
    };

    return (
        <div className={`flex flex-col h-full bg-white border-l border-gray-200 shadow-xl z-10 w-full ${fullscreen ? 'fixed inset-0 z-50 p-8' : ''}`}>
            <div className="p-4 border-b border-gray-100 flex items-center justify-between">
                <h2 className="font-semibold text-gray-800">Council Chat</h2>
                <div className="flex items-center gap-2">
                    <span className="text-xs text-green-600 bg-green-50 px-2 py-0.5 rounded-full">Online</span>
                    {fullscreen && <button onClick={onExitFullscreen} className="text-sm bg-gray-100 px-2 py-1 rounded">Min</button>}
                </div>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {messages.map((msg, idx) => (
                    <div key={idx} className={`flex gap-3 ${msg.role === 'user' ? 'flex-row-reverse' : ''}`}>
                        <div className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 ${msg.role === 'assistant' ? 'bg-indigo-100 text-indigo-600' : 'bg-gray-100 text-gray-600'}`}>
                            {msg.role === 'assistant' ? <Bot size={18} /> : <User size={18} />}
                        </div>

                        <div className={`max-w-[85%] p-3 rounded-2xl text-sm ${msg.role === 'user'
                            ? 'bg-blue-600 text-white rounded-br-none'
                            : 'bg-gray-50 border border-gray-100 text-gray-800 rounded-bl-none'
                            }`}>
                            <div className={`prose prose-sm max-w-none ${msg.role === 'user' ? 'prose-invert' : ''}`}>
                                <ReactMarkdown
                                    remarkPlugins={[remarkGfm]}
                                    rehypePlugins={[rehypeHighlight]}
                                >
                                    {msg.content}
                                </ReactMarkdown>
                            </div>
                        </div>
                    </div>
                ))}
                <div ref={messagesEndRef} />
            </div>

            <div className="p-4 border-t border-gray-100">
                <div className="relative">
                    <textarea
                        className="w-full p-3 pr-10 bg-gray-50 border border-gray-200 rounded-xl resize-none focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 min-h-[50px] max-h-[120px]"
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        onKeyDown={(e) => {
                            if (e.key === 'Enter' && !e.shiftKey) {
                                e.preventDefault();
                                handleSend();
                            }
                        }}
                        placeholder="Type a message..."
                        rows={1}
                    />
                    <button
                        onClick={handleSend}
                        disabled={!input.trim()}
                        className="absolute right-2 bottom-2.5 p-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                    >
                        <Send size={16} />
                    </button>
                </div>
            </div>
        </div>
    );
}
