import { useState } from 'react';
import { useWorkflowRunStore } from '../../../stores/useWorkflowRunStore';
import { AlertCircle, CheckCircle, XCircle, Clock } from 'lucide-react';
import { clsx } from 'clsx';

export const HumanReviewModal = () => {
    const humanReview = useWorkflowRunStore(state => state.humanReview);
    const submitHumanReview = useWorkflowRunStore(state => state.submitHumanReview);
    const [submitting, setSubmitting] = useState(false);
    const [feedback, setFeedback] = useState('');

    if (!humanReview) return null;

    const handleAction = async (action: 'approve' | 'reject') => {
        setSubmitting(true);
        try {
            await submitHumanReview(humanReview, action, { feedback }); // Sending feedback as data
        } catch (error) {
            console.error(error);
            // Optionally show error toast
        } finally {
            setSubmitting(false);
            setFeedback('');
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4">
            <div className="bg-white rounded-lg shadow-xl w-full max-w-lg overflow-hidden border border-gray-200">
                {/* Header */}
                <div className="bg-purple-50 p-4 border-b border-purple-100 flex items-center justify-between">
                    <div className="flex items-center gap-2 text-purple-800">
                        <AlertCircle className="w-5 h-5" />
                        <h3 className="font-semibold text-lg">Human Review Required</h3>
                    </div>
                </div>

                {/* Body */}
                <div className="p-6 space-y-4">
                    <div className="p-3 bg-gray-50 rounded border border-gray-200">
                        <div className="text-sm text-gray-500 mb-1">Reason</div>
                        <div className="font-medium text-gray-800">{humanReview.reason}</div>
                    </div>

                    <div className="flex items-center gap-2 text-sm text-gray-600">
                        <Clock className="w-4 h-4" />
                        <span>Timeout: <b>{humanReview.timeout}m</b></span>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">Feedback / Context (Optional)</label>
                        <textarea
                            className="w-full text-sm border-gray-300 rounded-md shadow-sm focus:border-purple-500 focus:ring-purple-500 p-2 border"
                            rows={3}
                            placeholder="Add notes or adjust context..."
                            value={feedback}
                            onChange={e => setFeedback(e.target.value)}
                        />
                    </div>
                </div>

                {/* Footer */}
                <div className="p-4 bg-gray-50 border-t border-gray-200 flex justify-end gap-3">
                    <button
                        onClick={() => handleAction('reject')}
                        disabled={submitting}
                        className={clsx(
                            "flex items-center gap-2 px-4 py-2 bg-white border border-red-300 text-red-700 rounded-md hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-red-500",
                            submitting && "opacity-50 cursor-not-allowed"
                        )}
                    >
                        <XCircle className="w-4 h-4" />
                        Reject (Stop)
                    </button>
                    <button
                        onClick={() => handleAction('approve')}
                        disabled={submitting}
                        className={clsx(
                            "flex items-center gap-2 px-4 py-2 bg-purple-600 text-white rounded-md hover:bg-purple-700 focus:outline-none focus:ring-2 focus:ring-purple-500",
                            submitting && "opacity-50 cursor-not-allowed"
                        )}
                    >
                        <CheckCircle className="w-4 h-4" />
                        Approve (Resume)
                    </button>
                </div>
            </div>
        </div>
    );
};
