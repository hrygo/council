import { type FC } from 'react';
import type { FactCheckNodeData } from '../../../../../types/workflow';

interface FactCheckNodeFormProps {
    data: FactCheckNodeData;
    onChange: (data: Partial<FactCheckNodeData>) => void;
}

export const FactCheckNodeForm: FC<FactCheckNodeFormProps> = ({ data, onChange }) => {

    const handleSourceToggle = (source: 'tavily' | 'serper' | 'local_kb') => {
        const current = data.search_sources || [];
        const next = current.includes(source)
            ? current.filter(s => s !== source)
            : [...current, source];
        onChange({ search_sources: next });
    };

    return (
        <div className="space-y-4">
            <div className="space-y-2">
                <label className="text-xs font-semibold text-gray-500 uppercase">Search Sources</label>
                <div className="space-y-2">
                    <label className="flex items-center gap-2 text-sm cursor-pointer">
                        <input
                            type="checkbox"
                            checked={data.search_sources?.includes('tavily')}
                            onChange={() => handleSourceToggle('tavily')}
                            className="w-4 h-4 rounded text-teal-600 focus:ring-teal-500"
                        />
                        <span>🌐 Tavily (Web)</span>
                    </label>
                    <label className="flex items-center gap-2 text-sm cursor-pointer">
                        <input
                            type="checkbox"
                            checked={data.search_sources?.includes('serper')}
                            onChange={() => handleSourceToggle('serper')}
                            className="w-4 h-4 rounded text-teal-600 focus:ring-teal-500"
                        />
                        <span>🔍 Serper (Web)</span>
                    </label>
                    <label className="flex items-center gap-2 text-sm cursor-pointer">
                        <input
                            type="checkbox"
                            checked={data.search_sources?.includes('local_kb')}
                            onChange={() => handleSourceToggle('local_kb')}
                            className="w-4 h-4 rounded text-teal-600 focus:ring-teal-500"
                        />
                        <span>📚 Local Knowledge Base</span>
                    </label>
                </div>
            </div>

            <div className="space-y-1.5">
                <div className="flex justify-between text-xs">
                    <label className="font-semibold text-gray-500 uppercase">Max Queries</label>
                    <span className="font-mono text-gray-700 dark:text-gray-300">{data.max_queries || 3}</span>
                </div>
                <input
                    type="range"
                    min="1"
                    max="10"
                    step="1"
                    value={data.max_queries || 3}
                    onChange={(e) => onChange({ max_queries: parseInt(e.target.value) })}
                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-teal-500"
                />
            </div>

            <div className="space-y-1.5">
                <div className="flex justify-between text-xs">
                    <label className="font-semibold text-gray-500 uppercase">Confidence Threshold</label>
                    <span className="font-mono text-gray-700 dark:text-gray-300">{Math.round((data.verify_threshold || 0.7) * 100)}%</span>
                </div>
                <input
                    type="range"
                    min="50"
                    max="100"
                    step="5"
                    value={(data.verify_threshold || 0.7) * 100}
                    onChange={(e) => onChange({ verify_threshold: parseInt(e.target.value) / 100 })}
                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer accent-teal-500"
                />
            </div>
        </div>
    );
};
