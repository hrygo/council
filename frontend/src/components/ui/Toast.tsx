import { createContext, useContext, useState, useCallback, type ReactNode, useEffect } from 'react';
import { X, CheckCircle, AlertCircle, Info, AlertTriangle } from 'lucide-react';
import clsx from 'clsx';

export type ToastType = 'success' | 'error' | 'info' | 'warning';

export interface Toast {
    id: string;
    message: string;
    type: ToastType;
    duration?: number;
}

interface ToastContextValue {
    addToast: (message: string, type: ToastType, duration?: number) => void;
    error: (message: string, duration?: number) => void;
    success: (message: string, duration?: number) => void;
    info: (message: string, duration?: number) => void;
    warning: (message: string, duration?: number) => void;
    removeToast: (id: string) => void;
}

const ToastContext = createContext<ToastContextValue | null>(null);

export function useToast() {
    const context = useContext(ToastContext);
    if (!context) {
        throw new Error('useToast must be used within a ToastProvider');
    }
    return context;
}

interface ToastProviderProps {
    children: ReactNode;
}

export function ToastProvider({ children }: ToastProviderProps) {
    const [toasts, setToasts] = useState<Toast[]>([]);

    const removeToast = useCallback((id: string) => {
        setToasts((prev) => prev.filter((t) => t.id !== id));
    }, []);

    const addToast = useCallback((message: string, type: ToastType, duration = 3000) => {
        const id = Math.random().toString(36).substring(2, 9);
        const newToast = { id, message, type, duration };
        setToasts((prev) => [...prev, newToast]);

        if (duration > 0) {
            setTimeout(() => {
                removeToast(id);
            }, duration);
        }
    }, [removeToast]);

    const success = useCallback((msg: string, dur?: number) => addToast(msg, 'success', dur), [addToast]);
    const error = useCallback((msg: string, dur?: number) => addToast(msg, 'error', dur), [addToast]);
    const info = useCallback((msg: string, dur?: number) => addToast(msg, 'info', dur), [addToast]);
    const warning = useCallback((msg: string, dur?: number) => addToast(msg, 'warning', dur), [addToast]);

    return (
        <ToastContext.Provider value={{ addToast, removeToast, success, error, info, warning }}>
            {children}
            <div className="fixed top-4 right-4 z-[9999] flex flex-col gap-2 pointer-events-none">
                {toasts.map((toast) => (
                    <ToastItem key={toast.id} toast={toast} onDismiss={() => removeToast(toast.id)} />
                ))}
            </div>
        </ToastContext.Provider>
    );
}

function ToastItem({ toast, onDismiss }: { toast: Toast; onDismiss: () => void }) {
    useEffect(() => {
        // Animation handled by CSS classes for now or simple timeout
        // Just render
    }, []);

    const icon = {
        success: <CheckCircle className="text-green-500" size={20} />,
        error: <AlertCircle className="text-red-500" size={20} />,
        warning: <AlertTriangle className="text-yellow-500" size={20} />,
        info: <Info className="text-blue-500" size={20} />,
    }[toast.type];

    const bgColors = {
        success: 'bg-white dark:bg-gray-800 border-green-100 dark:border-green-900',
        error: 'bg-white dark:bg-gray-800 border-red-100 dark:border-red-900',
        warning: 'bg-white dark:bg-gray-800 border-yellow-100 dark:border-yellow-900',
        info: 'bg-white dark:bg-gray-800 border-blue-100 dark:border-blue-900',
    }[toast.type];

    return (
        <div
            className={clsx(
                "pointer-events-auto flex items-center gap-3 p-4 rounded-lg shadow-lg border min-w-[300px] max-w-md animate-in slide-in-from-right-full fade-in duration-300",
                bgColors
            )}
            role="alert"
        >
            {icon}
            <p className="flex-1 text-sm font-medium text-gray-900 dark:text-gray-100">{toast.message}</p>
            <button
                onClick={onDismiss}
                className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 transition-colors"
            >
                <X size={16} />
            </button>
        </div>
    );
}
