import { type FC, useMemo } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ReferenceLine, ResponsiveContainer, Legend } from 'recharts';
import { TrendingUp, Target, Activity } from 'lucide-react';

interface ScoreHistoryEntry {
    round: number;
    score: number;
    timestamp?: string;
    event?: string;
}

interface LoopAnalyticsProps {
    scoreHistory: ScoreHistoryEntry[];
    threshold?: number;
    currentRound?: number;
}

export const LoopAnalytics: FC<LoopAnalyticsProps> = ({
    scoreHistory,
    threshold = 90,
    currentRound,
}) => {
    // Compute stats
    const stats = useMemo(() => {
        if (scoreHistory.length === 0) {
            return { baseline: 0, current: 0, delta: 0, trend: 'neutral' as const };
        }
        const baseline = scoreHistory[0].score;
        const current = scoreHistory[scoreHistory.length - 1].score;
        const delta = current - baseline;
        const trend = delta > 0 ? 'up' : delta < 0 ? 'down' : 'neutral';
        return { baseline, current, delta, trend };
    }, [scoreHistory]);

    // Prepare chart data
    const chartData = useMemo(() => {
        return scoreHistory.map((entry, idx) => ({
            round: entry.round || idx + 1,
            score: entry.score,
            event: entry.event,
        }));
    }, [scoreHistory]);

    if (scoreHistory.length === 0) {
        return (
            <div className="h-full flex flex-col items-center justify-center text-gray-500 dark:text-gray-400 gap-2 p-4">
                <Activity size={24} className="text-gray-300" />
                <p className="text-sm">No optimization data yet</p>
                <p className="text-xs">Score history will appear here during Loop execution.</p>
            </div>
        );
    }

    return (
        <div className="h-full flex flex-col bg-white dark:bg-gray-900">
            {/* Stats Header */}
            <div className="flex items-center justify-between px-4 py-3 border-b border-gray-200 dark:border-gray-700">
                <div className="flex items-center gap-2">
                    <TrendingUp size={16} className="text-blue-500" />
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                        Optimization Progress
                    </span>
                </div>
                <div className="flex items-center gap-4 text-xs">
                    <div className="flex items-center gap-1">
                        <span className="text-gray-500">Baseline:</span>
                        <span className="font-medium text-gray-700 dark:text-gray-300">{stats.baseline}</span>
                    </div>
                    <div className="flex items-center gap-1">
                        <span className="text-gray-500">Current:</span>
                        <span className={`font-medium ${stats.current >= threshold ? 'text-green-600' : 'text-blue-600'}`}>
                            {stats.current}
                        </span>
                    </div>
                    <div className="flex items-center gap-1">
                        <Target size={12} className="text-orange-500" />
                        <span className="text-gray-500">Target:</span>
                        <span className="font-medium text-orange-600">{threshold}</span>
                    </div>
                    <div className={`px-2 py-0.5 rounded-full text-xs font-medium ${stats.delta > 0
                        ? 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400'
                        : stats.delta < 0
                            ? 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400'
                            : 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                        }`}>
                        {stats.delta > 0 ? '+' : ''}{stats.delta} pts
                    </div>
                </div>
            </div>

            {/* Chart */}
            <div className="flex-1 p-4">
                <ResponsiveContainer width="100%" height="100%">
                    <LineChart data={chartData} margin={{ top: 10, right: 30, left: 0, bottom: 0 }}>
                        <CartesianGrid strokeDasharray="3 3" stroke="#374151" opacity={0.2} />
                        <XAxis
                            dataKey="round"
                            tick={{ fontSize: 11 }}
                            label={{ value: 'Round', position: 'insideBottom', offset: -5, fontSize: 11 }}
                        />
                        <YAxis
                            domain={[0, 100]}
                            tick={{ fontSize: 11 }}
                            label={{ value: 'Score', angle: -90, position: 'insideLeft', fontSize: 11 }}
                        />
                        <Tooltip
                            contentStyle={{
                                backgroundColor: '#1f2937',
                                border: 'none',
                                borderRadius: '8px',
                                fontSize: '12px',
                            }}
                            itemStyle={{ color: '#fff' }}
                            labelFormatter={(label) => `Round ${label}`}
                        />
                        <Legend />

                        {/* Threshold Line */}
                        <ReferenceLine
                            y={threshold}
                            stroke="#f97316"
                            strokeDasharray="5 5"
                            label={{ value: `Target: ${threshold}`, position: 'right', fontSize: 10, fill: '#f97316' }}
                        />

                        {/* Current Round Marker */}
                        {currentRound && (
                            <ReferenceLine
                                x={currentRound}
                                stroke="#3b82f6"
                                strokeDasharray="3 3"
                            />
                        )}

                        {/* Score Line */}
                        <Line
                            type="monotone"
                            dataKey="score"
                            stroke="#3b82f6"
                            strokeWidth={2}
                            dot={{ fill: '#3b82f6', strokeWidth: 2, r: 4 }}
                            activeDot={{ r: 6, fill: '#2563eb' }}
                        />
                    </LineChart>
                </ResponsiveContainer>
            </div>

            {/* Round Indicator */}
            <div className="px-4 py-2 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
                <div className="flex items-center gap-2 text-xs text-gray-500">
                    <span>Rounds completed: <strong>{scoreHistory.length}</strong></span>
                    {stats.current >= threshold && (
                        <span className="ml-auto px-2 py-0.5 bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400 rounded-full">
                            ✓ Target Reached
                        </span>
                    )}
                </div>
            </div>
        </div>
    );
};
