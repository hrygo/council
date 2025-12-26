import { type FC, useState } from 'react';
import { X, LayoutTemplate, Trash2, FileCode, Briefcase, Zap, Box } from 'lucide-react';
import { useTemplates, useDeleteTemplate } from '../../../hooks/useTemplates';
import type { Template, TemplateCategory } from '../../../types/template';

interface TemplateSidebarProps {
    open: boolean;
    onClose: () => void;
    onApply: (template: Template) => void;
}

const CATEGORY_ICONS: Record<TemplateCategory | 'other', FC<{ size?: number }>> = {
    code_review: FileCode,
    business_plan: Briefcase,
    quick_decision: Zap,
    custom: Box,
    other: LayoutTemplate
};

export const TemplateSidebar: FC<TemplateSidebarProps> = ({ open, onClose, onApply }) => {
    const { data: templates, isLoading } = useTemplates();
    const deleteTemplate = useDeleteTemplate();
    const [filter, setFilter] = useState<'all' | 'system' | 'custom'>('all');

    if (!open) return null;

    const systemTemplates = templates?.filter(t => t.is_system) || [];
    const customTemplates = templates?.filter(t => !t.is_system) || [];

    const filteredSystem = (filter === 'all' || filter === 'system') ? systemTemplates : [];
    const filteredCustom = (filter === 'all' || filter === 'custom') ? customTemplates : [];

    return (
        <div className="absolute left-0 top-0 bottom-0 w-80 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 z-20 flex flex-col shadow-xl animate-in slide-in-from-left duration-200">
            <div className="flex items-center justify-between p-4 border-b border-gray-100 dark:border-gray-700">
                <h2 className="font-semibold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                    <LayoutTemplate size={20} className="text-purple-500" />
                    Template Library
                </h2>
                <button onClick={onClose} className="p-1 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
                    <X size={20} className="text-gray-500" />
                </button>
            </div>

            <div className="p-4 border-b border-gray-100 dark:border-gray-700">
                <select
                    value={filter}
                    onChange={e => setFilter(e.target.value as 'all' | 'system' | 'custom')}
                    className="w-full px-3 py-2 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-500/50"
                >
                    <option value="all">All Templates</option>
                    <option value="system">System Templates</option>
                    <option value="custom">My Templates</option>
                </select>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-6">
                {isLoading && (
                    <div className="text-center py-8 text-gray-500 space-y-2">
                        <div className="animate-spin w-6 h-6 border-2 border-purple-500 border-t-transparent rounded-full mx-auto"></div>
                        <p className="text-xs">Loading templates...</p>
                    </div>
                )}

                {filteredSystem.length > 0 && (
                    <div className="space-y-3">
                        <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider">System Templates</h3>
                        {filteredSystem.map(template => (
                            <TemplateCard
                                key={template.template_uuid}
                                template={template}
                                onApply={() => onApply(template)}
                            />
                        ))}
                    </div>
                )}

                {filteredCustom.length > 0 && (
                    <div className="space-y-3">
                        <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider">My Templates</h3>
                        {filteredCustom.map(template => (
                            <TemplateCard
                                key={template.template_uuid}
                                template={template}
                                onApply={() => onApply(template)}
                                onDelete={() => deleteTemplate.mutate(template.template_uuid)}
                            />
                        ))}
                    </div>
                )}

                {/* Empty State */}
                {!isLoading && filteredSystem.length === 0 && filteredCustom.length === 0 && (
                    <div className="text-center py-12 text-gray-400">
                        <Box size={32} className="mx-auto mb-2 opacity-50" />
                        <p className="text-sm">No templates found</p>
                    </div>
                )}
            </div>
        </div>
    );
};

const TemplateCard: FC<{ template: Template; onApply: () => void; onDelete?: () => void }> = ({ template, onApply, onDelete }) => {
    const Icon = CATEGORY_ICONS[template.category] || CATEGORY_ICONS.other;

    return (
        <div
            onClick={onApply}
            className="group relative p-3 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl hover:border-purple-500 hover:shadow-sm cursor-pointer transition-all"
        >
            <div className="flex items-start gap-3">
                <div className="shrink-0 p-2 bg-purple-50 dark:bg-purple-900/20 text-purple-600 dark:text-purple-400 rounded-lg">
                    <Icon size={18} />
                </div>
                <div className="flex-1 min-w-0">
                    <h4 className="font-medium text-sm text-gray-900 dark:text-gray-100 truncate">{template.name}</h4>
                    <p className="text-xs text-gray-500 dark:text-gray-400 line-clamp-2 mt-0.5">{template.description}</p>
                </div>
            </div>

            {onDelete && (
                <button
                    onClick={(e) => { e.stopPropagation(); onDelete(); }}
                    className="absolute top-2 right-2 p-1.5 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded opacity-0 group-hover:opacity-100 transition-opacity"
                    title="Delete Template"
                >
                    <Trash2 size={14} />
                </button>
            )}
        </div>
    );
};
