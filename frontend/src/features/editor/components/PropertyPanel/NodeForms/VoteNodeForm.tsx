import { type FC } from 'react';
import type { VoteNodeData } from '../../../../../types/workflow';

interface VoteNodeFormProps {
    data: VoteNodeData;
    onChange: (data: Partial<VoteNodeData>) => void;
}

export const VoteNodeForm: FC<VoteNodeFormProps> = ({ data, onChange }) => {
    return (
        <div className="space-y-4">
            <div className="space-y-1.5">
                <div className="flex justify-between text-xs">
                    <label className="font-semibold text-gray-500 uppercase">Approval Threshold</label>
                    <span className="font-mono text-gray-700 dark:text-gray-300">{Math.round((data.threshold || 0.67) * 100)}%</span>
                </div>
                <input
                    type="range"
                    min="50"
                    max="100"
                    step="5"
                    value={(data.threshold || 0.67) * 100}
                    onChange={(e) => onChange({ threshold: parseInt(e.target.value) / 100 })}
                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-amber-500"
                />
            </div>

            <div className="space-y-2">
                <label className="text-xs font-semibold text-gray-500 uppercase">Vote Type</label>
                <div className="flex gap-2">
                    <button
                        className={`flex-1 py-1.5 px-3 rounded-lg text-xs font-medium border transition-colors ${data.vote_type === 'yes_no'
                                ? 'bg-amber-50 border-amber-200 text-amber-700 dark:bg-amber-900/20 dark:border-amber-800 dark:text-amber-400'
                                : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:bg-gray-50'
                            }`}
                        onClick={() => onChange({ vote_type: 'yes_no' })}
                    >
                        Yes/No
                    </button>
                    <button
                        className={`flex-1 py-1.5 px-3 rounded-lg text-xs font-medium border transition-colors ${data.vote_type === 'score_1_10'
                                ? 'bg-amber-50 border-amber-200 text-amber-700 dark:bg-amber-900/20 dark:border-amber-800 dark:text-amber-400'
                                : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 text-gray-600 dark:text-gray-400 hover:bg-gray-50'
                            }`}
                        onClick={() => onChange({ vote_type: 'score_1_10' })}
                    >
                        1-10 Score
                    </button>
                </div>
            </div>

            {/* Agent Select Placeholder - Waiting for Agent MultiSelect Component */}
            <div className="space-y-1.5">
                <label className="text-xs font-semibold text-gray-500 uppercase">Participants</label>
                <div className="p-3 bg-gray-50 dark:bg-gray-900/50 border border-dashed border-gray-200 dark:border-gray-700 rounded-lg text-xs text-gray-400 text-center">
                    Agent selection coming soon
                </div>
            </div>
        </div>
    );
};
