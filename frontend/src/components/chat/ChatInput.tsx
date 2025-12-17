import { useState, type FC, type KeyboardEvent } from 'react';
import { Send } from 'lucide-react';
import { useConnectStore } from '../../stores/useConnectStore'; // Assuming connect store handles direct messages or session store handles append?
// Spec says ChatPanel uses useSessionStore for display but ChatInput usually sends via WebSocket.
// In ChatPanel.tsx (legacy), it used useConnectStore().sendMessage.
// However, SessionStore might need to update optimistically or wait for server echo?
// Spec says: ChatPanel -> ChatInput.
// Let's use useConnectStore for sending for now, or useSessionStore if it has send action?
// SessionStore has 'appendMessage' but that's for local updates.
// Usually we send via WS.

interface ChatInputProps {
    sessionId: string;
}

export const ChatInput: FC<ChatInputProps> = ({ sessionId }) => {
    const [input, setInput] = useState('');
    const sendMessage = useConnectStore((state) => state.sendMessage);

    const handleSend = () => {
        if (!input.trim()) return;

        // We send to backend, backend broadcasts back, store picks it up via event listener
        sendMessage({
            type: 'user_input',
            content: input,
            session_id: sessionId
        });

        setInput('');
    };

    const onKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            handleSend();
        }
    };

    return (
        <div className="p-4 border-t border-gray-100 bg-white">
            <div className="relative">
                <textarea
                    className="w-full p-3 pr-12 bg-gray-50 border border-gray-200 rounded-xl resize-none focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-500 min-h-[50px] max-h-[120px] shadow-sm transition-all"
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    onKeyDown={onKeyDown}
                    placeholder="输入消息..." // i18n later
                    rows={1}
                />
                <button
                    onClick={handleSend}
                    disabled={!input.trim()}
                    className="absolute right-2 bottom-2 p-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                    <Send size={16} />
                </button>
            </div>
        </div>
    );
};
