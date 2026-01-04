import { type FC, type ReactNode } from 'react';
import { Handle, Position } from '@xyflow/react';
import clsx from 'clsx';
import { type LucideIcon, CheckCircle2, XCircle, Loader2, Circle } from 'lucide-react';
import { type NodeStatus } from '../../../types/session';

interface BaseNodeProps {
    label: string;
    icon?: LucideIcon;
    selected?: boolean;
    status?: NodeStatus;
    headerColor?: string; // Tailwind class like 'bg-blue-500'
    children?: ReactNode;
    handles?: ('top' | 'associative' | 'bottom')[]; // Which handles to show
}

export const BaseNode: FC<BaseNodeProps> = ({
    label,
    icon: Icon,
    selected,
    status = 'pending',
    headerColor = 'bg-gray-100',
    children,
    handles = ['top', 'bottom']
}) => {
    // Status visual mapping
    const statusStyles = {
        pending: {
            border: "border-gray-200 dark:border-gray-700",
            icon: Circle,
            iconColor: "text-gray-400"
        },
        running: {
            border: "border-blue-500 ring-2 ring-blue-200 dark:ring-blue-900/50",
            icon: Loader2,
            iconColor: "text-blue-500 animate-spin"
        },
        completed: {
            border: "border-green-500 ring-1 ring-green-200 dark:ring-green-900/30",
            icon: CheckCircle2,
            iconColor: "text-green-500"
        },
        failed: {
            border: "border-red-500 ring-1 ring-red-200 dark:ring-red-900/30",
            icon: XCircle,
            iconColor: "text-red-500"
        }
    };

    const currentStyle = statusStyles[status] || statusStyles.pending;
    const StatusIcon = currentStyle.icon;

    // Selection overrides border if not running/completed/failed, or adds to it?
    // Let's allow status to take precedence for border color.
    // Use selection for ring if not running (running uses blue ring).
    
    // Logic:
    // If selected, purple ring.
    // If running, blue ring (overrides purple ring? or mix?).
    // Let's prioritize running ring for status visibility.
    
    let borderClass = currentStyle.border;
    if (selected) {
        if (status === 'running') {
            // retain blue ring
        } else {
             // add purple ring/border overrides
             borderClass = "border-purple-500 ring-2 ring-purple-200 dark:ring-purple-900/50";
        }
    }

    return (
        <div className={clsx(
            "min-w-[150px] bg-white dark:bg-gray-800 rounded-lg shadow-sm border transition-all",
            borderClass
        )}>
            {/* Handles */}
            {handles.includes('top') && (
                <Handle type="target" position={Position.Top} className="!w-3 !h-3 !bg-gray-300 dark:!bg-gray-500" />
            )}

            <div className={clsx("px-3 py-2 rounded-t-lg flex items-center justify-between border-b border-gray-100 dark:border-gray-700", headerColor)}>
                <div className="flex items-center gap-2 overflow-hidden">
                    {Icon && <Icon size={14} className="text-gray-700 dark:text-gray-300 shrink-0" />}
                    <span className="text-xs font-semibold text-gray-700 dark:text-gray-200 truncate">{label}</span>
                </div>
                {/* Status Icon */}
                <StatusIcon size={14} className={clsx("shrink-0", currentStyle.iconColor)} />
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
