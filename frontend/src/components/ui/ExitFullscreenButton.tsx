import type { FC } from 'react';
import { Minimize2 } from 'lucide-react';
import clsx from 'clsx';

interface ExitFullscreenButtonProps {
    onClick?: () => void;
    className?: string;
}

export const ExitFullscreenButton: FC<ExitFullscreenButtonProps> = ({ onClick, className }) => {
    if (!onClick) return null;

    return (
        <button
            onClick={onClick}
            className={clsx(
                "flex items-center gap-2 px-3 py-1.5 rounded-lg shadow-sm",
                "bg-white dark:bg-gray-800",
                "border border-gray-200 dark:border-gray-700",
                "text-gray-700 dark:text-gray-200",
                "hover:bg-gray-50 dark:hover:bg-gray-700",
                "transition-all text-sm font-medium",
                className
            )}
        >
            <Minimize2 size={16} />
            <span>Exit Fullscreen</span>
        </button>
    );
};
