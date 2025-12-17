import React from 'react';
import type { NodeStatus } from '../types/session';

export const getNodeStatusStyles = (status: NodeStatus): React.CSSProperties => {
    switch (status) {
        case 'pending':
            return { opacity: 0.6 };
        case 'running':
            return {
                boxShadow: '0 0 0 2px #3B82F6',
                animation: 'pulse 1.5s ease-in-out infinite',
            };
        case 'completed':
            return {
                borderColor: '#10B981',
                boxShadow: '0 0 8px rgba(16, 185, 129, 0.3)',
            };
        case 'failed':
            return {
                borderColor: '#EF4444',
                boxShadow: '0 0 8px rgba(239, 68, 68, 0.3)',
            };
        default:
            return {};
    }
};

// 节点状态图标
export const getNodeStatusIcon = (status: NodeStatus): string => {
    switch (status) {
        case 'pending': return '⏳';
        case 'running': return '🔄';
        case 'completed': return '✅';
        case 'failed': return '❌';
        default: return '';
    }
};
