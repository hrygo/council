import type { FC } from 'react';
import { Minimize2 } from 'lucide-react';
import type { SessionStatus } from '../../types/session';

interface ChatHeaderProps {
    sessionStatus?: SessionStatus;
    onExitFullscreen?: () => void;
}

export const ChatHeader: FC<ChatHeaderProps> = ({ sessionStatus = 'idle', onExitFullscreen }) => {
    return (
        <div className="p-4 border-b border-gray-100 dark:border-gray-800 flex items-center justify-between bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm sticky top-0 z-10">
            <div className="flex items-center gap-3">
                <h2 className="font-semibold text-gray-800 dark:text-gray-200">Council Chat</h2>
                <div className={`px-2 py-0.5 rounded-full text-xs font-medium border
                ${sessionStatus === 'running' ? 'bg-green-50 dark:bg-green-900/30 text-green-600 dark:text-green-400 border-green-100 dark:border-green-800' : ''}
                ${sessionStatus === 'idle' ? 'bg-gray-50 dark:bg-gray-800 text-gray-600 dark:text-gray-400 border-gray-100 dark:border-gray-700' : ''}
                ${sessionStatus === 'failed' ? 'bg-red-50 dark:bg-red-900/30 text-red-600 dark:text-red-400 border-red-100 dark:border-red-800' : ''}
            `}>
                    {sessionStatus.toUpperCase()}
                </div>
            </div>

            {onExitFullscreen && (
                <button
                    onClick={onExitFullscreen}
                    className="p-1.5 text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
                    title="Exit Fullscreen"
                >
                    <Minimize2 size={18} />
                </button>
            )}
        </div>
    );
};
