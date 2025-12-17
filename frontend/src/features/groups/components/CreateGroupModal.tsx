import { useState, useEffect } from 'react';
import { X, Loader2 } from 'lucide-react';
import type { CreateGroupInput, Group } from '../../../types/group';

interface CreateGroupModalProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: CreateGroupInput) => void;
    isLoading?: boolean;
    initialData?: Group | null;
}

const ICONS = ['🏢', '🏠', '💼', '🎯', '⚙️', '📊', '🧪', '🎨', '🚀', '💡', '🤖', '🌐'];

export function CreateGroupModal({ open, onClose, onSubmit, isLoading, initialData }: CreateGroupModalProps) {
    const [name, setName] = useState('');
    const [icon, setIcon] = useState(ICONS[0]);
    const [systemPrompt, setSystemPrompt] = useState('');

    useEffect(() => {
        if (open) {
            if (initialData) {
                setName(initialData.name);
                setIcon(initialData.icon);
                setSystemPrompt(initialData.system_prompt);
            } else {
                setName('');
                setIcon(ICONS[0]);
                setSystemPrompt('');
            }
        }
    }, [open, initialData]);

    if (!open) return null;

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit({
            name,
            icon,
            system_prompt: systemPrompt,
            default_agent_ids: initialData?.default_agent_ids || [] // Preserve existing agents if editing
        });
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
            <div className="relative w-full max-w-lg bg-white dark:bg-gray-900 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-800 animate-in zoom-in-95 duration-200">
                <div className="flex items-center justify-between p-6 border-b border-gray-100 dark:border-gray-800">
                    <h2 className="text-xl font-semibold text-gray-900 dark:text-gray-100">
                        {initialData ? 'Edit Group' : 'Create New Group'}
                    </h2>
                    <button
                        onClick={onClose}
                        className="p-2 text-gray-400 hover:text-gray-500 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
                    >
                        <X size={20} />
                    </button>
                </div>

                <form onSubmit={handleSubmit} className="p-6 space-y-6">
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                            Group Name
                        </label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="e.g. Engineering Team"
                            required
                            className="w-full px-4 py-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                            Icon
                        </label>
                        <div className="grid grid-cols-6 gap-2">
                            {ICONS.map((i) => (
                                <button
                                    key={i}
                                    type="button"
                                    onClick={() => setIcon(i)}
                                    className={`p-2 text-xl rounded-lg border transition-all ${icon === i
                                        ? 'bg-blue-50 border-blue-500 ring-2 ring-blue-500/20 dark:bg-blue-900/20 dark:border-blue-500'
                                        : 'border-transparent hover:bg-gray-50 dark:hover:bg-gray-800'
                                        }`}
                                >
                                    {i}
                                </button>
                            ))}
                        </div>
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">
                            System Prompt
                        </label>
                        <textarea
                            value={systemPrompt}
                            onChange={(e) => setSystemPrompt(e.target.value)}
                            placeholder="Define the group's core purpose and context..."
                            rows={4}
                            className="w-full px-4 py-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all resize-none"
                        />
                    </div>

                    <div className="flex items-center justify-end gap-3 pt-2">
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={isLoading}
                            className="flex items-center gap-2 px-6 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {isLoading && <Loader2 size={16} className="animate-spin" />}
                            {initialData ? 'Save Changes' : 'Create Group'}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
}
