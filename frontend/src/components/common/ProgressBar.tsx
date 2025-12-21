import { type FC } from 'react';

interface ProgressBarProps {
    current: number;
    total: number;
    label?: string;
    color?: string;
}

export const ProgressBar: FC<ProgressBarProps> = ({
    current,
    total,
    label,
    color = 'bg-blue-600'
}) => {
    const percentage = Math.min(100, Math.max(0, (total > 0 ? (current / total) * 100 : 0)));

    return (
        <div className="w-full">
            {label && (
                <div className="flex justify-between text-xs mb-1">
                    <span className="text-gray-600 dark:text-gray-400 font-medium">{label}</span>
                    <span className="text-gray-500">{Math.round(percentage)}%</span>
                </div>
            )}
            <div className="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-1.5 overflow-hidden">
                <div
                    className={`h-1.5 rounded-full transition-all duration-500 ease-out ${color}`}
                    style={{ width: `${percentage}%` }}
                />
            </div>
        </div>
    );
};
