import { type FC } from 'react';
import { FileText, Target, Users, ArrowRight } from 'lucide-react';
import type { Template } from '../../types/template';

interface ConfirmPreviewProps {
    template: Template;
    documentContent: string;
    objective: string;
}

export const ConfirmPreview: FC<ConfirmPreviewProps> = ({ template, documentContent, objective }) => {
    // Count Agents logic (if not in template directly, deduce from graph)
    const agentCount = Object.values(template.graph.nodes || {}).filter((n) => n.type === 'agent').length;
    const docLength = documentContent.length;

    return (
        <div className="space-y-4 animate-in fade-in slide-in-from-bottom-2 duration-300">
            <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-100 dark:border-blue-800 rounded-xl p-4">
                <h3 className="text-sm font-semibold text-blue-900 dark:text-blue-100 flex items-center gap-2 mb-3">
                    <Target size={16} />
                    Session Summary
                </h3>

                <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                        <span className="text-gray-500 dark:text-gray-400 block text-xs">Template</span>
                        <span className="font-medium text-gray-900 dark:text-gray-100">{template.name}</span>
                    </div>
                    <div>
                        <span className="text-gray-500 dark:text-gray-400 block text-xs">Agents</span>
                        <span className="font-medium text-gray-900 dark:text-gray-100 flex items-center gap-1">
                            <Users size={12} /> {agentCount}
                        </span>
                    </div>
                </div>
            </div>

            <div className="space-y-3">
                <div className="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-100 dark:border-gray-700">
                    <FileText className="text-gray-400 mt-0.5" size={16} />
                    <div className="flex-1 min-w-0">
                        <div className="text-xs font-medium text-gray-500 uppercase mb-1">Document to Review</div>
                        <div className="text-sm text-gray-700 dark:text-gray-300 line-clamp-2 font-mono">
                            {docLength > 0 ? documentContent.slice(0, 150) + (docLength > 150 ? '...' : '') : <span className="text-gray-400 italic">No document provided</span>}
                        </div>
                        <div className="mt-1 text-xs text-gray-400">
                            {docLength} characters
                        </div>
                    </div>
                </div>

                <div className="flex items-start gap-3 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-lg border border-gray-100 dark:border-gray-700">
                    <Target className="text-gray-400 mt-0.5" size={16} />
                    <div className="flex-1 min-w-0">
                        <div className="text-xs font-medium text-gray-500 uppercase mb-1">Optimization Objective</div>
                        <div className="text-sm text-gray-700 dark:text-gray-300">
                            {objective || <span className="text-gray-400 italic">No specific objective defined</span>}
                        </div>
                    </div>
                </div>
            </div>

            <div className="flex items-center gap-2 text-xs text-gray-500 justify-center pt-2">
                <span>Configure</span>
                <ArrowRight size={12} />
                <span>Input</span>
                <ArrowRight size={12} />
                <span className="font-bold text-blue-600 dark:text-blue-400">Launch</span>
            </div>
        </div>
    );
};
