import type { FC } from 'react';
import { useWorkflowRunStore, getControlState } from '../../stores/useWorkflowRunStore';
import { Pause, Play, Square } from 'lucide-react';

interface ExecutionControlBarProps {
    sessionId: string;
}

// TODO: Move StatusDot to shared components if used elsewhere
const StatusDot: FC<{ status: string }> = ({ status }) => {
    let color = 'bg-gray-400';
    if (status === 'running') color = 'bg-blue-500 animate-pulse';
    if (status === 'completed') color = 'bg-green-500';
    if (status === 'failed') color = 'bg-red-500';
    if (status === 'paused') color = 'bg-yellow-500';

    return <div className={`w-2.5 h-2.5 rounded-full ${color}`} />;
};

export const ExecutionControlBar: FC<ExecutionControlBarProps> = ({ sessionId }) => {
    const executionStatus = useWorkflowRunStore((state) => state.executionStatus);
    const stats = useWorkflowRunStore((state) => state.stats);
    const sendControl = useWorkflowRunStore((state) => state.sendControl);

    // Check if using derived state getter directly in component or selector is better. 
    // Using simple derivation here since we have executionStatus
    const controlState = getControlState(executionStatus);

    const formatTime = (ms: number) => {
        const seconds = Math.floor(ms / 1000);
        const minutes = Math.floor(seconds / 60);
        return `${minutes}:${(seconds % 60).toString().padStart(2, '0')}`;
    };

    return (
        <div className="flex items-center gap-4 p-2 bg-gray-50 rounded-lg border border-gray-100 shadow-sm">
            {/* 状态指示器 */}
            <div className="flex items-center gap-2 px-2">
                <StatusDot status={executionStatus} />
                <span className="text-sm font-medium capitalize text-gray-700">{executionStatus}</span>
            </div>

            {/* 控制按钮 */}
            <div className="flex gap-2">
                {controlState.canPause && (
                    <button
                        className="flex items-center px-3 py-1.5 text-xs font-medium text-gray-700 bg-white border border-gray-300 rounded hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                        onClick={() => sendControl(sessionId, 'pause')}
                    >
                        <Pause size={14} className="mr-1.5" /> 暂停
                    </button>
                )}

                {controlState.canResume && (
                    <button
                        className="flex items-center px-3 py-1.5 text-xs font-medium text-white bg-blue-600 border border-transparent rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                        onClick={() => sendControl(sessionId, 'resume')}
                    >
                        <Play size={14} className="mr-1.5" /> 继续
                    </button>
                )}

                {controlState.canStop && (
                    <button
                        className="flex items-center px-3 py-1.5 text-xs font-medium text-white bg-red-600 border border-transparent rounded hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                        onClick={() => sendControl(sessionId, 'stop')}
                    >
                        <Square size={14} className="mr-1.5" /> 停止
                    </button>
                )}
            </div>

            {/* 统计信息 */}
            <div className="ml-auto flex items-center gap-4 text-xs text-gray-500 font-mono">
                <span title="Elapsed Time">⏱️ {formatTime(stats.elapsedTimeMs)}</span>
                <span title="Progress">📊 {stats.completedNodes}/{stats.totalNodes}</span>
                <span title="Total Cost">💰 ${stats.totalCostUsd.toFixed(4)}</span>
            </div>
        </div>
    );
};
