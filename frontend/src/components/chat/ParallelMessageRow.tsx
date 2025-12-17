import type { FC } from 'react';
import type { Message } from '../../types/session';
import { ParallelMessageCard } from './ParallelMessageCard';

interface ParallelMessageRowProps {
    messages: Message[];
    maxColumns?: number;
}

const accentColors = [
    'border-t-blue-500',
    'border-t-green-500',
    'border-t-purple-500',
    'border-t-orange-500',
    'border-t-pink-500',
    'border-t-cyan-500'
];

export const ParallelMessageRow: FC<ParallelMessageRowProps> = ({
    messages, maxColumns = 3,
}) => (
    <div
        className="grid gap-4 w-full"
        style={{
            gridTemplateColumns: `repeat(${Math.min(messages.length, maxColumns)}, minmax(0, 1fr))`
        }}
    >
        {messages.map((msg, idx) => (
            <ParallelMessageCard
                key={msg.id}
                message={msg}
                index={idx}
                accentColor={accentColors[idx % accentColors.length]}
            />
        ))}
    </div>
);
