import { createContext } from 'react';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
    id: string;
    message: string;
    type: ToastType;
    duration?: number;
}

export interface ToastContextValue {
    addToast: (message: string, type: ToastType, duration?: number) => void;
    error: (message: string, duration?: number) => void;
    success: (message: string, duration?: number) => void;
    info: (message: string, duration?: number) => void;
    warning: (message: string, duration?: number) => void;
    removeToast: (id: string) => void;
}

export const ToastContext = createContext<ToastContextValue | null>(null);
