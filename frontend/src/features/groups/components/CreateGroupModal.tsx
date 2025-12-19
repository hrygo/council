import { useState } from 'react';
import { X, Loader2, Check } from 'lucide-react';
import type { CreateGroupInput, Group } from '../../../types/group';
import { useAgents } from '../../../hooks/useAgents';

interface CreateGroupModalProps {
    open: boolean;
    onClose: () => void;
    onSubmit: (data: CreateGroupInput) => void;
    isLoading?: boolean;
    initialData?: Group | null;
}

const ICONS = ['🏢', '🏠', '💼', '🎯', '⚙️', '📊', '🧪', '🎨', '🚀', '💡', '🤖', '🌐'];

export function CreateGroupModal({ open, onClose, onSubmit, isLoading, initialData }: CreateGroupModalProps) {
    const { data: agents } = useAgents();
    const [name, setName] = useState(() => initialData?.name || '');
    const [icon, setIcon] = useState(() => initialData?.icon || ICONS[0]);
    const [systemPrompt, setSystemPrompt] = useState(() => initialData?.system_prompt || '');
    const [selectedAgents, setSelectedAgents] = useState<string[]>(() => initialData?.default_agent_ids || []);

    if (!open) return null;

    const handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        onSubmit({
            name,
            icon,
            system_prompt: systemPrompt,
            default_agent_ids: selectedAgents
        });
    };

    const toggleAgent = (agentId: string) => {
        setSelectedAgents(prev =>
            prev.includes(agentId)
                ? prev.filter(id => id !== agentId)
                : [...prev, agentId]
        );
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
            <div className="relative w-full max-w-2xl bg-white dark:bg-gray-900 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-800 animate-in zoom-in-95 duration-200 flex flex-col max-h-[90vh]">
                <div className="flex items-center justify-between p-6 border-b border-gray-100 dark:border-gray-800 shrink-0">
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

                <div className="overflow-y-auto flex-1 p-6">
                    <form id="group-form" onSubmit={handleSubmit} className="space-y-6">
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
                            <div className="grid grid-cols-6 sm:grid-cols-12 gap-2">
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
                                rows={3}
                                className="w-full px-4 py-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all resize-none"
                            />
                        </div>

                        <div className="space-y-3">
                            <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex justify-between items-center">
                                <span>Initial Agents</span>
                                <span className="text-xs font-normal text-gray-500">{selectedAgents.length} selected</span>
                            </label>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 max-h-60 overflow-y-auto pr-1">
                                {agents?.map(agent => {
                                    const isSelected = selectedAgents.includes(agent.id);
                                    return (
                                        <div
                                            key={agent.id}
                                            onClick={() => toggleAgent(agent.id)}
                                            className={`flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-all select-none ${isSelected
                                                ? 'bg-blue-50 dark:bg-blue-900/20 border-blue-500 dark:border-blue-500'
                                                : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-700'
                                                }`}
                                        >
                                            <div className={`w-5 h-5 rounded border flex items-center justify-center transition-colors ${isSelected
                                                ? 'bg-blue-600 border-blue-600 text-white'
                                                : 'border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800'
                                                }`}>
                                                {isSelected && <Check size={12} />}
                                            </div>
                                            <div className="flex items-center gap-2 overflow-hidden">
                                                <span className="text-xl flex-shrink-0">{agent.avatar || '🤖'}</span>
                                                <div className="flex flex-col min-w-0">
                                                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{agent.name}</span>
                                                    <span className="text-xs text-gray-500 truncate">{agent.model_config.model}</span>
                                                </div>
                                            </div>
                                        </div>
                                    );
                                })}
                                {(!agents || agents.length === 0) && (
                                    <div className="col-span-full py-8 text-center text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-dashed border-gray-200 dark:border-gray-700">
                                        No agents available. Create some agents first!
                                    </div>
                                )}
                            </div>
                        </div>
                    </form>
                </div>

                <div className="flex items-center justify-end gap-3 p-6 border-t border-gray-100 dark:border-gray-800 shrink-0">
                    <button
                        type="button"
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        type="submit"
                        form="group-form"
                        disabled={isLoading}
                        className="flex items-center gap-2 px-6 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        {isLoading && <Loader2 size={16} className="animate-spin" />}
                        {initialData ? 'Save Changes' : 'Create Group'}
                    </button>
                </div>
            </div>
        </div>
    );
}
