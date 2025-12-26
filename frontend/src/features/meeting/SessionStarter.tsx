import { type FC, useState } from 'react';
import { Play, Sparkles, LayoutTemplate, MessageSquare, ArrowRight, ArrowLeft } from 'lucide-react';
import { useTemplates } from '../../hooks/useTemplates';
import { useSessionStore } from '../../stores/useSessionStore';
import { useConnectStore } from '../../stores/useConnectStore';
import { useNavigate } from 'react-router-dom';
import { useWorkflowRunStore } from '../../stores/useWorkflowRunStore';
import { FileUploadZone } from '../../components/common/FileUploadZone';
import { ConfirmPreview } from './ConfirmPreview';

interface SessionStarterProps {
    onStarted: () => void;
}

type Step = 'template' | 'input' | 'confirm';

export const SessionStarter: FC<SessionStarterProps> = ({ onStarted }) => {
    const navigate = useNavigate();
    const { data: templates, isLoading } = useTemplates();
    const initSession = useSessionStore(state => state.initSession);

    const [step, setStep] = useState<Step>('template');
    const [selected_template_uuid, setSelectedTemplateId] = useState<string>('');

    // Inputs
    const [documentContent, setDocumentContent] = useState('');
    const [objective, setObjective] = useState('');

    const [isStarting, setIsStarting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const selectedTemplate = templates?.find(t => t.template_uuid === selected_template_uuid);

    const handleNext = () => {
        if (step === 'template' && selected_template_uuid) {
            setStep('input');
        } else if (step === 'input') {
            // Validate input if needed
            setStep('confirm');
        }
    };

    const handleBack = () => {
        if (step === 'input') setStep('template');
        else if (step === 'confirm') setStep('input');
    };

    const handleStart = async () => {
        if (!selectedTemplate) return;

        setIsStarting(true);
        setError(null);

        try {
            // 1. Prepare Payload
            const payload = {
                graph: selectedTemplate.graph,
                input: {
                    document_content: documentContent,
                    optimization_objective: objective
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
            useWorkflowRunStore.getState().setGraphFromTemplate(selectedTemplate);

            const rawNodes = selectedTemplate.graph?.nodes || {};
            // eslint-disable-next-line @typescript-eslint/no-explicit-any
            const nodes = Object.values(rawNodes).map((n: any) => ({
                node_id: n.node_id,
                name: n.name || n.node_id,
                type: n.type || 'agent'
            }));

            initSession({
                session_uuid: data.session_uuid,
                workflow_id: selectedTemplate.template_uuid,
                group_uuid: 'default',
                nodes: nodes
            });

            // 4. Establish WebSocket connection
            const wsHost = window.location.host;
            const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = `${wsProtocol}//${wsHost}/ws`;
            useConnectStore.getState().connect(wsUrl);

            // 5. Navigate
            onStarted();
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

    const renderStepContent = () => {
        switch (step) {
            case 'template':
                return (
                    <div className="space-y-3 animate-in fade-in slide-in-from-right-4 duration-300">
                        <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
                            <LayoutTemplate size={16} />
                            Select Workflow Template
                        </label>

                        {isLoading ? (
                            <div className="animate-pulse h-10 bg-gray-100 dark:bg-gray-700 rounded-lg" />
                        ) : (
                            <div className="grid grid-cols-1 gap-3 max-h-64 overflow-y-auto pr-1 custom-scrollbar">
                                {systemTemplates.map(t => (
                                    <div
                                        key={t.template_uuid}
                                        onClick={() => setSelectedTemplateId(t.template_uuid)}
                                        className={`p-3 rounded-xl border cursor-pointer transition-all flex items-center justify-between ${selected_template_uuid === t.template_uuid
                                            ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 ring-1 ring-blue-500'
                                            : 'border-gray-200 dark:border-gray-700 hover:border-blue-300 dark:hover:border-blue-700'
                                            }`}
                                    >
                                        <div>
                                            <div className="font-medium text-sm text-gray-900 dark:text-gray-100">{t.name}</div>
                                            <div className="text-xs text-gray-500 dark:text-gray-400 line-clamp-1">{t.description}</div>
                                        </div>
                                        {selected_template_uuid === t.template_uuid && (
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
                );
            case 'input':
                return (
                    <div className="space-y-4 animate-in fade-in slide-in-from-right-4 duration-300">
                        {/* Document Upload */}
                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
                                <span className="flex items-center justify-center w-5 h-5 rounded-full bg-blue-100 text-blue-600 text-xs font-bold">1</span>
                                Upload Document
                            </label>
                            <FileUploadZone
                                onContentChange={(content) => setDocumentContent(content)}
                                initialContent={documentContent}
                            />
                        </div>

                        {/* Objective Input */}
                        <div className="space-y-2">
                            <label className="text-sm font-medium text-gray-700 dark:text-gray-300 flex items-center gap-2">
                                <span className="flex items-center justify-center w-5 h-5 rounded-full bg-blue-100 text-blue-600 text-xs font-bold">2</span>
                                Optimization Objective
                                <MessageSquare size={14} className="text-gray-400" />
                            </label>
                            <textarea
                                value={objective}
                                onChange={(e) => setObjective(e.target.value)}
                                placeholder="e.g. Optimize for clarity and conciseness, maintain professional tone..."
                                className="w-full h-24 px-4 py-3 bg-gray-50 dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl focus:ring-2 focus:ring-blue-500 focus:outline-none transition-all resize-none text-sm"
                            />
                        </div>
                    </div>
                );
            case 'confirm':
                return selectedTemplate ? (
                    <ConfirmPreview
                        template={selectedTemplate}
                        documentContent={documentContent}
                        objective={objective}
                    />
                ) : null;
        }
    };

    const getButtonText = () => {
        if (isStarting) return 'Starting Session...';
        if (step === 'confirm') return 'Start Council Session';
        return 'Next Step';
    };

    const isNextDisabled = () => {
        if (isStarting) return true;
        if (step === 'template') return !selected_template_uuid;
        // In input step, technically inputs are optional if not enforced, 
        // but for better UX let's say at least one is needed? Or not.
        // Let's allow empty inputs for broad compatibility.
        return false;
    };

    return (
        <div className="absolute inset-0 z-40 flex items-center justify-center bg-gray-100/50 dark:bg-gray-900/50 backdrop-blur-sm p-4 animate-in fade-in duration-300">
            <div className="w-full max-w-lg bg-white dark:bg-gray-800 rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col max-h-[90vh]">
                <div className="p-6 border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex-none">
                    <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                        <Sparkles className="text-blue-500" size={24} />
                        Start New Session
                    </h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
                        {step === 'template' && 'Select a workflow template to begin.'}
                        {step === 'input' && 'Provide context for the AI agents.'}
                        {step === 'confirm' && 'Review configuration before launch.'}
                    </p>
                </div>

                <div className="p-6 flex-1 overflow-y-auto">
                    {renderStepContent()}

                    {error && (
                        <div className="mt-4 p-3 bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm rounded-lg">
                            {error}
                        </div>
                    )}
                </div>

                <div className="p-6 border-t border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/50 flex-none flex gap-3">
                    {step !== 'template' && (
                        <button
                            onClick={handleBack}
                            disabled={isStarting}
                            className="px-4 py-3 rounded-xl font-medium text-gray-600 dark:text-gray-300 hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors flex items-center gap-2"
                        >
                            <ArrowLeft size={16} /> Back
                        </button>
                    )}

                    <button
                        onClick={step === 'confirm' ? handleStart : handleNext}
                        disabled={isNextDisabled()}
                        className={`flex-1 py-3 rounded-xl font-semibold text-white shadow-lg shadow-blue-500/20 flex items-center justify-center gap-2 transition-all ${isNextDisabled()
                            ? 'bg-gray-300 dark:bg-gray-700 cursor-not-allowed'
                            : 'bg-blue-600 hover:bg-blue-700 hover:scale-[1.02] active:scale-[0.98]'
                            }`}
                    >
                        {isStarting ? (
                            <div className="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin" />
                        ) : (
                            <>
                                {step === 'confirm' ? <Play size={20} className="fill-current" /> : null}
                                {getButtonText()}
                                {step !== 'confirm' && <ArrowRight size={16} />}
                            </>
                        )}
                    </button>
                </div>
            </div>
        </div>
    );
};
