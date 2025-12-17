import { type FC } from 'react';
import type { HumanReviewNodeData } from '../../../../../types/workflow';
import { AlertCircle } from 'lucide-react';

interface HumanReviewNodeFormProps {
    data: HumanReviewNodeData;
    onChange: (data: Partial<HumanReviewNodeData>) => void;
}

export const HumanReviewNodeForm: FC<HumanReviewNodeFormProps> = ({ data, onChange }) => {
    return (
        <div className="space-y-4">
            <div className="p-3 bg-amber-50 dark:bg-amber-900/20 text-amber-700 dark:text-amber-400 text-xs rounded-lg flex items-start gap-2">
                <AlertCircle size={14} className="mt-0.5 shrink-0" />
                <span>This node will pause execution and wait for human intervention.</span>
            </div>

            <div className="space-y-1.5">
                <label className="text-xs font-semibold text-gray-500 uppercase">Review Type</label>
                <select
                    value={data.review_type || 'approve_reject'}
                    onChange={(e) => onChange({ review_type: e.target.value as HumanReviewNodeData['review_type'] })}
                    className="w-full px-3 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-rose-500/50"
                >
                    <option value="approve_reject">Approve / Reject</option>
                    <option value="edit_content">Content Editing</option>
                </select>
            </div>

            <div className="space-y-1.5">
                <div className="flex justify-between text-xs">
                    <label className="font-semibold text-gray-500 uppercase">Timeout (Minutes)</label>
                    <span className="font-mono text-gray-700 dark:text-gray-300">{data.timeout_minutes || 30}m</span>
                </div>
                <input
                    type="range"
                    min="5"
                    max="60"
                    step="5"
                    value={data.timeout_minutes || 30}
                    onChange={(e) => onChange({ timeout_minutes: parseInt(e.target.value) })}
                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-rose-500"
                />
            </div>

            <div className="flex items-center justify-between">
                <label className="text-xs font-semibold text-gray-500 uppercase">Allow Skip</label>
                <div className="relative inline-block w-10 h-5">
                    <input
                        type="checkbox"
                        checked={data.allow_skip || false}
                        onChange={(e) => onChange({ allow_skip: e.target.checked })}
                        className="peer sr-only"
                        id="skip_switch"
                    />
                    <label htmlFor="skip_switch" className="block h-5 bg-gray-300 rounded-full cursor-pointer peer-focus:outline-none peer-checked:bg-rose-500 peer-checked:after:translate-x-full after:content-[''] after:absolute after:top-0.5 after:left-0.5 after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all"></label>
                </div>
            </div>
            <p className="text-[10px] text-gray-400">
                If enabled, automatically approves after timeout.
            </p>
        </div>
    );
};
