import type { FC } from 'react';
import { Minimize2 } from 'lucide-react';
import type { SessionStatus } from '../../types/session';

interface ChatHeaderProps {
    sessionStatus?: SessionStatus;
    onExitFullscreen?: () => void;
}

export const ChatHeader: FC<ChatHeaderProps> = ({ sessionStatus = 'idle', onExitFullscreen }) => {
    return (
        <div className="p-4 border-b border-gray-100 flex items-center justify-between bg-white/80 backdrop-blur-sm sticky top-0 z-10">
            <div className="flex items-center gap-3">
                <h2 className="font-semibold text-gray-800">Council Chat</h2>
                <div className={`px-2 py-0.5 rounded-full text-xs font-medium border
                ${sessionStatus === 'running' ? 'bg-green-50 text-green-600 border-green-100' : ''}
                ${sessionStatus === 'idle' ? 'bg-gray-50 text-gray-600 border-gray-100' : ''}
                ${sessionStatus === 'failed' ? 'bg-red-50 text-red-600 border-red-100' : ''}
            `}>
                    {sessionStatus.toUpperCase()}
                </div>
            </div>

            {onExitFullscreen && (
                <button
                    onClick={onExitFullscreen}
                    className="p-1.5 text-gray-500 hover:bg-gray-100 rounded-lg transition-colors"
                    title="Exit Fullscreen"
                >
                    <Minimize2 size={18} />
                </button>
            )}
        </div>
    );
};
