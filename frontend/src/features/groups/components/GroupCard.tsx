import { Users, Pencil, Trash2 } from 'lucide-react';
import type { Group } from '../../../types/group';

interface GroupCardProps {
    group: Group;
    onEdit: (group: Group) => void;
    onDelete: (group: Group) => void;
}

export function GroupCard({ group, onEdit, onDelete }: GroupCardProps) {
    return (
        <div className="group relative bg-white dark:bg-gray-800 rounded-xl p-6 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-md hover:border-blue-500/50 dark:hover:border-blue-500/50 transition-all duration-200">
            <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                    <div className="w-12 h-12 rounded-lg bg-gradient-to-br from-blue-500/10 to-indigo-500/10 flex items-center justify-center text-2xl">
                        {group.icon || '📦'}
                    </div>
                    <div>
                        <h3 className="font-semibold text-lg text-gray-900 dark:text-gray-100">
                            {group.name}
                        </h3>
                        <div className="flex items-center gap-1.5 mt-1 text-sm text-gray-500 dark:text-gray-400">
                            <Users size={14} />
                            <span>{group.default_agent_uuids?.length || 0} members</span>
                        </div>
                    </div>
                </div>

                <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                        onClick={(e) => { e.stopPropagation(); onEdit(group); }}
                        className="p-2 text-gray-400 hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-colors"
                        title="Edit Group"
                    >
                        <Pencil size={16} />
                    </button>
                    <button
                        onClick={(e) => { e.stopPropagation(); onDelete(group); }}
                        className="p-2 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                        title="Delete Group"
                    >
                        <Trash2 size={16} />
                    </button>
                </div>
            </div>

            <p className="mt-4 text-sm text-gray-600 dark:text-gray-300 line-clamp-2">
                {group.system_prompt || 'No system prompt configured.'}
            </p>
        </div>
    );
}
