import { type FC, useState } from 'react';
import { Play, Sparkles, LayoutTemplate, MessageSquare } from 'lucide-react';
import { useTemplates } from '../../hooks/useTemplates';
import { useSessionStore } from '../../stores/useSessionStore';

import { useNavigate } from 'react-router-dom';

interface SessionStarterProps {
    onStarted: () => void;
}

export const SessionStarter: FC<SessionStarterProps> = ({ onStarted }) => {
    const navigate = useNavigate();
    const { data: templates, isLoading } = useTemplates();
    const initSession = useSessionStore(state => state.initSession);

    const [selectedTemplateId, setSelectedTemplateId] = useState<string>('');
    const [topic, setTopic] = useState('');
    const [isStarting, setIsStarting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleStart = async () => {
        if (!selectedTemplateId || !topic) return;

        setIsStarting(true);
        setError(null);

        try {
            const template = templates?.find(t => t.id === selectedTemplateId);
            if (!template) throw new Error("Selected template not found");

            // 1. Prepare Payload
            const payload = {
                graph: template.graph,
                input: {
                    topic: topic
                }
            };

            // 2. Call API
            const res = await fetch('/api/v1/workflows/execute', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (!res.ok) {
                const errData = await res.json().catch(() => ({}));
                throw new Error(errData.error || 'Failed to start session');
            }

            const data = await res.json(); // { session_id, status }

            // 3. Initialize Store
            // We need a flat list of nodes for the session store to track status
            const rawNodes = template.graph?.nodes || {};
            const nodes = Object.values(rawNodes).map((n: { id: string; name?: string; type?: string }) => ({
                id: n.id,
                name: n.name || n.id,
                type: n.type || 'agent'
            }));

            initSession({
                sessionId: data.session_id,
                workflowId: template.id,
                groupId: 'default',
                nodes: nodes
            });

            // 4. Navigate to Meeting
            onStarted(); // Close modal
            // We need to use router to navigate. 
            // Since this component is inside a Router, we can use useNavigate.
            // I will add the hook import and usage.
            navigate('/meeting');

        } catch (err: unknown) {
            console.error(err);
            const message = err instanceof Error ? err.message : "Failed to start session";
            setError(message);
        } finally {
            setIsStarting(false);
        }
    };

    const systemTemplates = templates?.filter(t => t.is_system) || [];

    return (
        <div className="absolute inset-0 z-40 flex items-center justify-center bg-gray-100/50 dark:bg-gray-900/50 backdrop-blur-sm p-4 animate-in fade-in duration-300">
            <div className="w-full max-w-lg bg-white dark:bg-gray-800 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-700 overflow-hidden">
                <div className="p-6 border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50">
                    <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                        <Sparkles className="text-blue-500" size={24} />
                        Start New Session
                    </h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        Select a workflow template and define your topic to begin.
                    </p>
                </div>

                <div className="p-6 space-y-6">
                    {/* Template Selection */}
                    <div className="space-y-3">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
                            <LayoutTemplate size={16} />
                            Select Template
                        </label>

                        {isLoading ? (
                            <div className="animate-pulse h-10 bg-gray-100 dark:bg-gray-700 rounded-lg" />
                        ) : (
                            <div className="grid grid-cols-1 gap-3 max-h-48 overflow-y-auto pr-1 custom-scrollbar">
                                {systemTemplates.map(t => (
                                    <div
                                        key={t.id}
                                        onClick={() => setSelectedTemplateId(t.id)}
                                        className={`p-3 rounded-xl border cursor-pointer transition-all flex items-center justify-between ${selectedTemplateId === t.id
                                            ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 ring-1 ring-blue-500'
                                            : 'border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-700'
                                            }`}
                                    >
                                        <div>
                                            <div className="font-medium text-sm text-gray-900 dark:text-gray-100">{t.name}</div>
                                            <div className="text-xs text-gray-500 dark:text-gray-400 line-clamp-1">{t.description}</div>
                                        </div>
                                        {selectedTemplateId === t.id && (
                                            <div className="w-2 h-2 rounded-full bg-blue-500" />
                                        )}
                                    </div>
                                ))}
                                {systemTemplates.length === 0 && (
                                    <div className="text-sm text-gray-500 text-center py-4">No system templates found.</div>
                                )}
                            </div>
                        )}
                    </div>

                    {/* Topic Input */}
                    <div className="space-y-3">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
                            <MessageSquare size={16} />
                            Discussion Topic
                        </label>
                        <textarea
                            value={topic}
                            onChange={(e) => setTopic(e.target.value)}
                            placeholder="e.g. Should artificial intelligence be regulated by international governments?"
                            className="w-full h-24 px-4 py-3 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none transition-all resize-none text-sm"
                        />
                    </div>

                    {error && (
                        <div className="p-3 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm rounded-lg">
                            {error}
                        </div>
                    )}

                    <button
                        onClick={handleStart}
                        disabled={!selectedTemplateId || !topic || isStarting}
                        className={`w-full py-3 rounded-xl font-semibold text-white shadow-lg shadow-blue-500/20 flex items-center justify-center gap-2 transition-all ${!selectedTemplateId || !topic || isStarting
                            ? 'bg-gray-300 dark:bg-gray-700 cursor-not-allowed'
                            : 'bg-blue-600 hover:bg-blue-700 hover:scale-[1.02] active:scale-[0.98]'
                            }`}
                    >
                        {isStarting ? (
                            <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                        ) : (
                            <>
                                <Play size={20} className="fill-current" />
                                Start Council Session
                            </>
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
};
