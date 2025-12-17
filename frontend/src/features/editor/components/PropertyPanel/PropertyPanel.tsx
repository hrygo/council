import { type FC } from 'react';
import type { WorkflowNode, NodeType } from '../../../../types/workflow';
import { X, Settings2 } from 'lucide-react';

// Import Forms (will be implemented next)
import { VoteNodeForm } from './NodeForms/VoteNodeForm';
import { LoopNodeForm } from './NodeForms/LoopNodeForm';
import { FactCheckNodeForm } from './NodeForms/FactCheckNodeForm';
import { HumanReviewNodeForm } from './NodeForms/HumanReviewNodeForm';

interface PropertyPanelProps {
    node: WorkflowNode | null;
    onUpdate: (nodeId: string, data: Record<string, unknown>) => void;
    onDelete: (nodeId: string) => void;
    onClose: () => void;
}

export const PropertyPanel: FC<PropertyPanelProps> = ({ node, onUpdate, onDelete, onClose }) => {
    if (!node) return null;

    const renderForm = () => {
        switch (node.type as NodeType) {
            case 'vote':
                // @ts-expect-error - data casting handled in form
                return <VoteNodeForm data={node.data} onChange={(d) => onUpdate(node.id, d)} />;
            case 'loop':
                // @ts-expect-error - data casting handled in form
                return <LoopNodeForm data={node.data} onChange={(d) => onUpdate(node.id, d)} />;
            case 'fact_check':
                // @ts-expect-error - data casting handled in form
                return <FactCheckNodeForm data={node.data} onChange={(d) => onUpdate(node.id, d)} />;
            case 'human_review':
                // @ts-expect-error - data casting handled in form
                return <HumanReviewNodeForm data={node.data} onChange={(d) => onUpdate(node.id, d)} />;

            // For other nodes, show default or specific placeholder
            case 'agent':
                return <div className="p-4 text-sm text-gray-500">Agent Node Configuration (Coming Soon)</div>;
            default:
                return <div className="p-4 text-sm text-gray-500">No specific configuration for this node type.</div>;
        }
    };

    return (
        <div className="absolute right-4 top-4 w-80 bg-white dark:bg-gray-800 rounded-xl shadow-xl border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col max-h-[calc(100vh-2rem)] z-50 animate-in slide-in-from-right-5 duration-200">
            <div className="flex items-center justify-between p-4 border-b border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-900/50">
                <div className="flex items-center gap-2 font-medium text-gray-900 dark:text-gray-100">
                    <Settings2 size={16} className="text-blue-500" />
                    <span>{node.data.label || 'Node Properties'}</span>
                </div>
                <button onClick={onClose} className="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                    <X size={18} />
                </button>
            </div>

            <div className="flex-1 overflow-y-auto">
                <div className="p-4 space-y-4">
                    {/* Common Name Input */}
                    <div className="space-y-1.5">
                        <label className="text-xs font-semibold text-gray-500 uppercase tracking-wider">Node Name</label>
                        <input
                            type="text"
                            value={node.data.label as string || ''}
                            onChange={(e) => onUpdate(node.id, { label: e.target.value })}
                            className="w-full px-3 py-2 bg-gray-50 dark:bg-gray-900/50 border border-gray-200 dark:border-gray-700 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-blue-500/50"
                        />
                    </div>

                    <div className="h-px bg-gray-100 dark:bg-gray-800 my-2" />

                    {renderForm()}
                </div>
            </div>

            <div className="p-4 border-t border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-900/50">
                <button
                    onClick={() => onDelete(node.id)}
                    className="w-full py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                >
                    Delete Node
                </button>
            </div>
        </div>
    );
};
