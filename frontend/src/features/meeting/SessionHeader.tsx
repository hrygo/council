import { type FC } from 'react';
import { useTranslation } from 'react-i18next';
import { useSessionStore } from '../../stores/useSessionStore';
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore';
import { Activity, Clock, Coins, Hash, Minimize2 } from 'lucide-react';

interface SessionHeaderProps {
    onExitFullscreen?: () => void;
}

const StatusBadge: FC<{ status: string }> = ({ status }) => {
    let color = 'bg-gray-100 text-gray-600';
    if (status === 'running') color = 'bg-blue-100 text-blue-700 animate-pulse';
    if (status === 'completed') color = 'bg-green-100 text-green-700';
    if (status === 'failed') color = 'bg-red-100 text-red-700';
    if (status === 'paused') color = 'bg-amber-100 text-amber-700';

    return (
        <span className={`px-2 py-0.5 rounded-full text-xs font-bold uppercase tracking-wide flex items-center gap-1.5 ${color}`}>
            <div className={`w-1.5 h-1.5 rounded-full ${color.replace('bg-', 'bg-current ')}`} />
            {status}
        </span>
    );
};

export const SessionHeader: FC<SessionHeaderProps> = ({ onExitFullscreen }) => {
    const { t } = useTranslation();
    const currentSession = useSessionStore(state => state.currentSession);
    const stats = useWorkflowRunStore(state => state.stats);
    const executionStatus = useWorkflowRunStore(state => state.executionStatus);

    if (!currentSession) return null;

    const formatTime = (ms: number) => {
        const s = Math.floor(ms / 1000);
        const m = Math.floor(s / 60);
        const sec = s % 60;
        return `${m}:${sec.toString().padStart(2, '0')}`;
    };

    return (
        <div className="flex items-center justify-between px-6 py-3 border-b border-gray-100 dark:border-gray-800 bg-white/50 dark:bg-gray-900/50 backdrop-blur-sm z-10 sticky top-0">
            <div className="flex items-center gap-4">
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
                        <Activity size={20} className="text-blue-600 dark:text-blue-400" />
                    </div>
                    <div>
                        <h1 className="text-sm font-bold text-gray-900 dark:text-white leading-tight">{t('meeting.councilSession')}</h1>
                        <span className="text-xs text-gray-500 font-mono hidden sm:inline-block">{currentSession.session_uuid.slice(0, 8)}</span>
                    </div>
                </div>
                <StatusBadge status={executionStatus} />
            </div>

            <div className="flex items-center gap-4 sm:gap-6 text-xs font-medium text-gray-600 dark:text-gray-400">
                <div className="flex items-center gap-1.5 hidden md:flex" title="Elapsed Time">
                    <Clock size={14} className="text-gray-400" />
                    <span className="font-mono">{formatTime(stats.elapsedTimeMs)}</span>
                </div>

                <div className="flex items-center gap-1.5 hidden md:flex" title="Total Tokens">
                    <Hash size={14} className="text-gray-400" />
                    <span>{stats.totalTokens.toLocaleString()} T</span>
                </div>

                <div className="flex items-center gap-1.5 hidden sm:flex" title="Total Cost">
                    <Coins size={14} className="text-amber-500" />
                    <span>${stats.totalCostUsd.toFixed(4)}</span>
                </div>

                {onExitFullscreen && (
                    <button
                        onClick={onExitFullscreen}
                        className="p-1.5 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg ml-2 transition-colors"
                        title="Exit Fullscreen"
                    >
                        <Minimize2 size={16} />
                    </button>
                )}
            </div>
        </div>
    );
};
