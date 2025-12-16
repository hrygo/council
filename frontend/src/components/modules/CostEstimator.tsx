import type { FC } from 'react';

export const CostEstimator: FC = () => {
    // Mock values
    const estimatedCost = 0.05; // $
    const timeToComplete = 45; // seconds

    return (
        <div className="flex items-center gap-4 text-xs text-gray-500 mb-2 px-2">
            <div className="flex items-center gap-1">
                <span className="font-semibold">Est. Cost:</span>
                <span>${estimatedCost.toFixed(4)}</span>
            </div>
            <div className="flex items-center gap-1">
                <span className="font-semibold">Time:</span>
                <span>~{timeToComplete}s</span>
            </div>
        </div>
    );
};
