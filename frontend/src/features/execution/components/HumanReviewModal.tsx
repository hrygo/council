import { useState } from 'react';
import { useWorkflowRunStore } from '../../../stores/useWorkflowRunStore';
import { AlertCircle, CheckCircle, XCircle, Clock, FileCode, Edit3 } from 'lucide-react';
import { clsx } from 'clsx';
import { CodeDiffViewer } from '../../../components/code/CodeDiffViewer';

interface ToolCall {
    name: string;
    arguments: {
        path?: string;
        content?: string;
        [key: string]: unknown;
    };
}

export const HumanReviewModal = () => {
    const humanReview = useWorkflowRunStore(state => state.humanReview);
    const submitHumanReview = useWorkflowRunStore(state => state.submitHumanReview);
    const [submitting, setSubmitting] = useState(false);
    const [feedback, setFeedback] = useState('');
    const [editedContent, setEditedContent] = useState<string | null>(null);

    if (!humanReview) return null;

    // Detect if this is a code review (Surgeon's write_file)
    const toolCalls: ToolCall[] = humanReview.payload?.tool_calls || [];
    const writeFileCall = toolCalls.find(tc => tc.name === 'write_file');
    const isCodeReview = !!writeFileCall;

    // Get original content from VFS if available
    const originalContent = humanReview.payload?.original_content || '';
    const newContent = writeFileCall?.arguments?.content || '';
    const filePath = writeFileCall?.arguments?.path || 'unknown';

    const handleAction = async (action: 'approve' | 'reject') => {
        setSubmitting(true);
        try {
            const data: Record<string, unknown> = { feedback };

            // If user edited the content, include it
            if (isCodeReview && editedContent !== null) {
                data.modified_content = editedContent;
            }

            await submitHumanReview(humanReview, action, data);
        } catch (error) {
            console.error(error);
        } finally {
            setSubmitting(false);
            setFeedback('');
            setEditedContent(null);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
            <div className={clsx(
                "bg-white dark:bg-gray-900 rounded-lg shadow-xl overflow-hidden border border-gray-200 dark:border-gray-700 flex flex-col",
                isCodeReview ? "w-full max-w-4xl h-[80vh]" : "w-full max-w-lg"
            )}>
                {/* Header */}
                <div className={clsx(
                    "p-4 border-b flex items-center justify-between",
                    isCodeReview
                        ? "bg-blue-50 dark:bg-blue-900/20 border-blue-100 dark:border-blue-800"
                        : "bg-purple-50 dark:bg-purple-900/20 border-purple-100 dark:border-purple-800"
                )}>
                    <div className={clsx(
                        "flex items-center gap-2",
                        isCodeReview ? "text-blue-800 dark:text-blue-300" : "text-purple-800 dark:text-purple-300"
                    )}>
                        {isCodeReview ? (
                            <>
                                <FileCode className="w-5 h-5" />
                                <h3 className="font-semibold text-lg">Code Review Required</h3>
                            </>
                        ) : (
                            <>
                                <AlertCircle className="w-5 h-5" />
                                <h3 className="font-semibold text-lg">Human Review Required</h3>
                            </>
                        )}
                    </div>
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-400">
                        <Clock className="w-4 h-4" />
                        <span>Timeout: <b>{humanReview.timeout}m</b></span>
                    </div>
                </div>

                {/* Body */}
                <div className="flex-1 overflow-hidden flex flex-col">
                    {/* Reason */}
                    <div className="p-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
                        <div className="text-sm text-gray-500 dark:text-gray-400 mb-1">Reason</div>
                        <div className="font-medium text-gray-800 dark:text-gray-200">{humanReview.reason}</div>
                    </div>

                    {/* Code Diff (if applicable) */}
                    {isCodeReview && (
                        <div className="flex-1 p-4 overflow-hidden">
                            <div className="flex items-center gap-2 mb-2 text-sm text-gray-600 dark:text-gray-400">
                                <Edit3 size={14} />
                                <span>File: <code className="bg-gray-100 dark:bg-gray-800 px-1 rounded">{filePath}</code></span>
                            </div>
                            <div className="h-[calc(100%-2rem)] rounded-lg overflow-hidden border border-gray-200 dark:border-gray-700">
                                <CodeDiffViewer
                                    originalContent={originalContent}
                                    newContent={newContent}
                                    filePath={filePath}
                                    editable={true}
                                    onContentChange={setEditedContent}
                                />
                            </div>
                        </div>
                    )}

                    {/* Feedback Input */}
                    <div className="p-4 border-t border-gray-200 dark:border-gray-700">
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
                            Feedback / Context (Optional)
                        </label>
                        <textarea
                            className="w-full text-sm border-gray-300 dark:border-gray-600 rounded-md shadow-sm focus:border-purple-500 focus:ring-purple-500 p-2 border bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                            rows={isCodeReview ? 2 : 3}
                            placeholder="Add notes or adjust context..."
                            value={feedback}
                            onChange={e => setFeedback(e.target.value)}
                        />
                    </div>
                </div>

                {/* Footer */}
                <div className="p-4 bg-gray-50 dark:bg-gray-800/50 border-t border-gray-200 dark:border-gray-700 flex justify-end gap-3">
                    <button
                        onClick={() => handleAction('reject')}
                        disabled={submitting}
                        className={clsx(
                            "flex items-center gap-2 px-4 py-2 bg-white dark:bg-gray-800 border border-red-300 dark:border-red-700 text-red-700 dark:text-red-400 rounded-md hover:bg-red-50 dark:hover:bg-red-900/20 focus:outline-none focus:ring-2 focus:ring-red-500",
                            submitting && "opacity-50 cursor-not-allowed"
                        )}
                    >
                        <XCircle className="w-4 h-4" />
                        Reject
                    </button>
                    <button
                        onClick={() => handleAction('approve')}
                        disabled={submitting}
                        className={clsx(
                            "flex items-center gap-2 px-4 py-2 text-white rounded-md focus:outline-none focus:ring-2",
                            isCodeReview
                                ? "bg-blue-600 hover:bg-blue-700 focus:ring-blue-500"
                                : "bg-purple-600 hover:bg-purple-700 focus:ring-purple-500",
                            submitting && "opacity-50 cursor-not-allowed"
                        )}
                    >
                        <CheckCircle className="w-4 h-4" />
                        {isCodeReview ? 'Apply Changes' : 'Approve'}
                    </button>
                </div>
            </div>
        </div>
    );
};
