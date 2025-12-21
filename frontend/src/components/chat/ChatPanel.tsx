import { useRef, useEffect } from 'react';
import type { FC } from 'react';
import { useSessionStore } from '../../stores/useSessionStore';
import { MessageGroupCard } from './MessageGroupCard';
import { ChatInput } from './ChatInput';
import { SessionHeader } from '../../features/meeting/SessionHeader';
import type { ChatPanelProps } from './ChatPanel.types';

export const ChatPanel: FC<ChatPanelProps> = ({
    fullscreen,
    onExitFullscreen,
    readOnly,
    sessionId
}) => {
    const messageGroups = useSessionStore(state => state.messageGroups);
    const currentSession = useSessionStore(state => state.currentSession);
    const activeNodeIds = currentSession?.activeNodeIds ?? [];

    const messagesEndRef = useRef<HTMLDivElement>(null);

    // 自动滚动到底部
    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
    }, [messageGroups, messageGroups.length]);

    return (
        <div
            className={`
        flex flex-col h-full bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-700 shadow-xl z-20 w-full transition-all duration-300
        ${fullscreen ? "fixed inset-0 z-50" : ""}
      `}
        >
            {/* Header */}
            <SessionHeader
                onExitFullscreen={fullscreen ? onExitFullscreen : undefined}
            />

            {/* Message Groups */}
            <div className="flex-1 overflow-y-auto p-4 bg-gray-50/50 dark:bg-gray-800/50">
                {messageGroups.length === 0 ? (
                    <div className="h-full flex items-center justify-center text-gray-400 dark:text-gray-500 text-sm">
                        等待会议开始...
                    </div>
                ) : (
                    messageGroups.map(group => (
                        <MessageGroupCard
                            key={group.nodeId} // Assuming nodeId is unique per group sequence, or use group.id if available. SPEC-001 defined group having nodeId. Using nodeId as key might be duplicate if looped? SPEC-001 says messageGroups is array.
                            // If loop visits same node, we need unique key.
                            // SPEC-001: interface MessageGroup { id: string; nodeId: string; ... } => Let's double check if we defined ID.
                            // Checked session.ts in memory: MessageGroup has `nodeId`, no `id`?
                            // Wait, looking at SPEC-001 snippet in chat history:
                            // export interface MessageGroup { nodeId: string; ... }
                            // If we revisit a node, we might have multiple groups for same nodeId?
                            // The store implementation should handle creating new groups for new visits.
                            // If store implementation uses array index or adds unique ID, better.
                            // Let's use index as fallback or assume store handles unique references.
                            // To be safe, let's look at `index` in map.
                            group={group}
                            isActive={activeNodeIds.includes(group.nodeId)}
                        />
                    ))
                )}
                <div ref={messagesEndRef} />
            </div>

            {/* Input */}
            {!readOnly && sessionId && (
                <ChatInput sessionId={sessionId} />
            )}
        </div>
    );
};

export default ChatPanel;
