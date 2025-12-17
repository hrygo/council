import type { FC } from 'react';

interface StatusDotProps {
    status: string;
}

/**
 * StatusDot - Visual indicator for execution/node status
 * Reusable across components that need status visualization.
 */
export const StatusDot: FC<StatusDotProps> = ({ status }) => {
    let color = 'bg-gray-400';
    if (status === 'running') color = 'bg-blue-500 animate-pulse';
    if (status === 'completed') color = 'bg-green-500';
    if (status === 'failed') color = 'bg-red-500';
    if (status === 'paused') color = 'bg-yellow-500';
    if (status === 'suspended') color = 'bg-purple-500 animate-pulse';
    if (status === 'pending') color = 'bg-gray-300';

    return <div className={`w-2.5 h-2.5 rounded-full ${color}`} />;
};
