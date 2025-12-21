import { type FC } from 'react';

interface AgentAvatarProps {
    name: string;
    size?: 'sm' | 'md' | 'lg';
}

const COLORS = [
    'bg-red-500', 'bg-orange-500', 'bg-amber-500',
    'bg-green-500', 'bg-emerald-500', 'bg-teal-500',
    'bg-cyan-500', 'bg-sky-500', 'bg-blue-500',
    'bg-indigo-500', 'bg-violet-500', 'bg-purple-500',
    'bg-fuchsia-500', 'bg-pink-500', 'bg-rose-500'
];

function stringToColorIndex(str: string): number {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    return Math.abs(hash % COLORS.length);
}

const SIZES = {
    sm: 'w-6 h-6 text-xs',
    md: 'w-8 h-8 text-sm',
    lg: 'w-10 h-10 text-base'
};

export const AgentAvatar: FC<AgentAvatarProps> = ({ name, size = 'md' }) => {
    // Determine color deterministically based on name
    const colorIndex = stringToColorIndex(name);
    const bgClass = COLORS[colorIndex];
    const initial = name.charAt(0).toUpperCase();

    // Special case for 'system' or 'user' roles if needed, 
    // but usually we rely on name.

    return (
        <div
            className={`${SIZES[size]} rounded-full flex items-center justify-center text-white font-bold shadow-sm ${bgClass} shrink-0`}
            title={name}
        >
            {initial}
        </div>
    );
};
