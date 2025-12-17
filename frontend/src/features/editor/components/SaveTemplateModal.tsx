import { type FC, useState } from 'react';
import { useCreateTemplate } from '../../../hooks/useTemplates';
import type { BackendGraph } from '../../../utils/graphUtils';
import type { TemplateCategory } from '../../../types/template';
import { X, Loader2, Save } from 'lucide-react';

interface SaveTemplateModalProps {
    open: boolean;
    onClose: () => void;
    currentGraph: BackendGraph | null;
}

export const SaveTemplateModal: FC<SaveTemplateModalProps> = ({ open, onClose, currentGraph }) => {
    const { mutate: createTemplate, isPending } = useCreateTemplate();
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [category, setCategory] = useState<TemplateCategory>('custom');

    if (!open || !currentGraph) return null;

    const handleSave = (e: React.FormEvent) => {
        e.preventDefault();

        // Ensure graph has latest name/desc if we want to sync them, 
        // but here we are creating a specific template wrapper.

        createTemplate({
            name,
            description,
            category,
            graph: currentGraph
        }, {
            onSuccess: () => {
                alert('Template saved successfully!');
                onClose();
                setName('');
                setDescription('');
            }
        });
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
            <div className="bg-white dark:bg-gray-900 w-full max-w-md rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-800 animate-in zoom-in-95 duration-200 overflow-hidden">
                <div className="flex items-center justify-between p-4 border-b border-gray-100 dark:border-gray-800">
                    <h2 className="font-semibold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                        <Save size={18} className="text-purple-500" />
                        Save as Template
                    </h2>
                    <button onClick={onClose} className="p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                        <X size={20} />
                    </button>
                </div>

                <form onSubmit={handleSave} className="p-6 space-y-4">
                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">Template Name</label>
                        <input
                            type="text"
                            value={name}
                            onChange={e => setName(e.target.value)}
                            placeholder="e.g. Optimized Code Review"
                            className="w-full px-3 py-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500/50"
                            required
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">Description</label>
                        <textarea
                            value={description}
                            onChange={e => setDescription(e.target.value)}
                            placeholder="Describe what this workflow does..."
                            rows={3}
                            className="w-full px-3 py-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500/50 resize-none"
                        />
                    </div>

                    <div className="space-y-2">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300">Category</label>
                        <select
                            value={category}
                            onChange={e => setCategory(e.target.value as TemplateCategory)}
                            className="w-full px-3 py-2 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500/50"
                        >
                            <option value="code_review">Code Review</option>
                            <option value="business_plan">Business Plan</option>
                            <option value="quick_decision">Quick Decision</option>
                            <option value="custom">Custom</option>
                            <option value="other">Other</option>
                        </select>
                    </div>

                    <div className="bg-purple-50 dark:bg-purple-900/10 p-3 rounded-lg text-xs text-purple-700 dark:text-purple-300 border border-purple-100 dark:border-purple-900/20">
                        <strong>Preview:</strong> Saving {Object.keys(currentGraph.nodes).length} nodes configuration.
                    </div>

                    <div className="flex justify-end gap-3 pt-2">
                        <button type="button" onClick={onClose} className="px-4 py-2 text-sm text-gray-600 dark:text-gray-400 font-medium hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg">
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={isPending || !name}
                            className="px-4 py-2 text-sm text-white bg-purple-600 hover:bg-purple-700 rounded-lg font-medium flex items-center gap-2 disabled:opacity-50"
                        >
                            {isPending && <Loader2 size={16} className="animate-spin" />}
                            Save Template
                        </button>
                    </div>
                </form>
            </div>
        </div>
    );
};
