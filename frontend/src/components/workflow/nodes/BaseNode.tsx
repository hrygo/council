import { type FC, type ReactNode } from 'react';
import { Handle, Position } from '@xyflow/react';
import clsx from 'clsx';
import { type LucideIcon } from 'lucide-react';

interface BaseNodeProps {
    label: string;
    icon?: LucideIcon;
    selected?: boolean;
    headerColor?: string; // Tailwind class like 'bg-blue-500'
    children?: ReactNode;
    handles?: ('top' | 'associative' | 'bottom')[]; // Which handles to show
}

export const BaseNode: FC<BaseNodeProps> = ({
    label,
    icon: Icon,
    selected,
    headerColor = 'bg-gray-100',
    children,
    handles = ['top', 'bottom']
}) => {
    return (
        <div className={clsx(
            "min-w-[150px] bg-white dark:bg-gray-800 rounded-lg shadow-sm border transition-all",
            selected ? "border-purple-500 ring-2 ring-purple-200 dark:ring-purple-900/50" : "border-gray-200 dark:border-gray-700",
        )}>
            {/* Handles */}
            {handles.includes('top') && (
                <Handle type="target" position={Position.Top} className="!w-3 !h-3 !bg-gray-300 dark:!bg-gray-500" />
            )}

            <div className={clsx("px-3 py-2 rounded-t-lg flex items-center gap-2 border-b border-gray-100 dark:border-gray-700", headerColor)}>
                {Icon && <Icon size={14} className="text-gray-700 dark:text-gray-300" />}
                <span className="text-xs font-semibold text-gray-700 dark:text-gray-200 truncate">{label}</span>
            </div>

            <div className="p-3 text-xs text-gray-600 dark:text-gray-400">
                {children}
            </div>

            {handles.includes('bottom') && (
                <Handle type="source" position={Position.Bottom} className="!w-3 !h-3 !bg-gray-300 dark:!bg-gray-500" />
            )}
        </div>
    );
};
