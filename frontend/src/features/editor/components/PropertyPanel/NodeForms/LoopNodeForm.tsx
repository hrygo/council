import { type FC } from 'react';
import type { LoopNodeData } from '../../../../../types/workflow';

interface LoopNodeFormProps {
    data: LoopNodeData;
    onChange: (data: Partial<LoopNodeData>) => void;
}

export const LoopNodeForm: FC<LoopNodeFormProps> = ({ data, onChange }) => {
    return (
        <div className="space-y-4">
            <div className="space-y-1.5">
                <div className="flex justify-between text-xs">
                    <label className="font-semibold text-gray-500 uppercase">Max Rounds</label>
                    <span className="font-mono text-gray-700 dark:text-gray-300">{data.max_rounds || 3}</span>
                </div>
                <input
                    type="range"
                    min="1"
                    max="10"
                    step="1"
                    value={data.max_rounds || 3}
                    onChange={(e) => onChange({ max_rounds: parseInt(e.target.value) })}
                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-purple-500"
                />
            </div>

            <div className="space-y-1.5">
                <label className="text-xs font-semibold text-gray-500 uppercase">Exit Condition</label>
                <select
                    value={data.exit_condition || 'max_rounds'}
                    onChange={(e) => onChange({ exit_condition: e.target.value as LoopNodeData['exit_condition'] })}
                    className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500/50"
                >
                    <option value="max_rounds">Reaching Max Rounds</option>
                    <option value="consensus">Reaching Consensus</option>
                </select>
            </div>

            {/* Agent Pairs Placeholder */}
            <div className="space-y-1.5">
                <label className="text-xs font-semibold text-gray-500 uppercase">Debate Pairs</label>
                <div className="p-3 bg-gray-50 dark:bg-gray-900/50 border border-dashed border-gray-200 dark:border-gray-700 rounded-lg text-xs text-gray-400 text-center">
                    Debate pairs configuration coming soon
                </div>
            </div>
        </div>
    );
};
