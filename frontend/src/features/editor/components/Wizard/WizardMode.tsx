import { type FC, useState } from 'react';
import { Sparkles, ArrowRight, CheckCircle, Wand2, X, AlertCircle } from 'lucide-react';
import { useGenerateWorkflow } from '../../../../hooks/useGenerateWorkflow';
import type { BackendGraph } from '../../../../utils/graphUtils';
import type { Template } from '../../../../types/template';

interface WizardModeProps {
    open: boolean;
    onClose: () => void;
    onComplete: (graph: BackendGraph) => void;
}

export const WizardMode: FC<WizardModeProps> = ({ open, onClose, onComplete }) => {
    const { generate, isGenerating, error } = useGenerateWorkflow();
    const [step, setStep] = useState<1 | 2 | 3>(1);
    const [intent, setIntent] = useState('');
    const [generatedResult, setGeneratedResult] = useState<{ graph: BackendGraph; similar: Template[] } | null>(null);

    if (!open) return null;

    const handleGenerate = async () => {
        if (!intent.trim()) return;
        const result = await generate(intent);
        if (result && result.graph) {
            setGeneratedResult({
                graph: result.graph,
                similar: result.similar_templates || []
            });
            setStep(2);
        }
    };

    const handleSelectGraph = (graph: BackendGraph) => {
        onComplete(graph);
        onClose();
        // Reset state?
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm animate-in fade-in duration-300">
            <div className="bg-white dark:bg-gray-900 w-full max-w-2xl rounded-2xl shadow-2xl border border-gray-200 dark:border-gray-800 flex flex-col overflow-hidden h-[600px] animate-in zoom-in-95 duration-200">
                {/* Header */}
                <div className="p-6 border-b border-gray-100 dark:border-gray-800 flex items-center justify-between">
                    <div>
                        <h2 className="text-xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                            <Wand2 className="text-purple-600" />
                            Workflow Wizard
                        </h2>
                        <p className="text-sm text-gray-500 mt-1">AI-powered workflow generation</p>
                    </div>
                    <button onClick={onClose} className="p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full transition-colors">
                        <X size={20} className="text-gray-500" />
                    </button>
                </div>

                {/* Progress */}
                <div className="px-6 py-4 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-100 dark:border-gray-800 flex items-center justify-center gap-8">
                    {[1, 2, 3].map(s => (
                        <div key={s} className={`flex items-center gap-2 ${step >= s ? 'text-purple-600 font-medium' : 'text-gray-400'}`}>
                            <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs border ${step >= s ? 'bg-purple-100 border-purple-600 text-purple-700' : 'border-gray-300 text-gray-500'
                                }`}>
                                {step > s ? <CheckCircle size={14} /> : s}
                            </div>
                            <span className="text-sm text-nowrap">
                                {s === 1 ? 'Describe Intent' : s === 2 ? 'Review & Select' : 'Refine'}
                            </span>
                        </div>
                    ))}
                </div>

                {/* Content */}
                <div className="flex-1 p-8 overflow-y-auto">
                    {step === 1 && (
                        <div className="space-y-6 max-w-lg mx-auto">
                            <div className="text-center space-y-2">
                                <h3 className="text-lg font-semibold">What kind of meeting do you want to run?</h3>
                                <p className="text-gray-500 text-sm">Describe your goal, participants, and desired outcome.</p>
                            </div>

                            <textarea
                                value={intent}
                                onChange={e => setIntent(e.target.value)}
                                placeholder="Example: I need a strict code review process for a high-risk security patch involving specific agents..."
                                className="w-full h-40 p-4 bg-gray-50 dark:bg-gray-800 border-2 border-dashed border-gray-200 dark:border-gray-700 rounded-xl focus:border-purple-500 focus:outline-none resize-none transition-colors"
                            />

                            {error && (
                                <div className="p-3 bg-red-50 text-red-600 rounded-lg text-sm flex items-center gap-2">
                                    <AlertCircle size={16} />
                                    {error}
                                </div>
                            )}

                            <button
                                onClick={handleGenerate}
                                disabled={!intent.trim() || isGenerating}
                                className="w-full py-3 bg-purple-600 hover:bg-purple-700 text-white rounded-xl font-medium flex items-center justify-center gap-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed group"
                            >
                                {isGenerating ? (
                                    <>
                                        <Sparkles className="animate-spin" size={18} /> Generating Magic...
                                    </>
                                ) : (
                                    <>
                                        Generate Workflow <ArrowRight size={18} className="group-hover:translate-x-1 transition-transform" />
                                    </>
                                )}
                            </button>
                        </div>
                    )}

                    {step === 2 && generatedResult && (
                        <div className="space-y-6">
                            <div className="flex items-center justify-between">
                                <h3 className="text-lg font-medium">Generation Results</h3>
                                <div className="text-sm text-green-600 flex items-center gap-1">
                                    <Sparkles size={14} /> AI Success
                                </div>
                            </div>

                            {/* Generated Card */}
                            <div
                                onClick={() => handleSelectGraph(generatedResult.graph)}
                                className="p-5 border-2 border-purple-500 bg-purple-50 dark:bg-purple-900/10 rounded-xl cursor-pointer hover:shadow-lg transition-all relative overflow-hidden group"
                            >
                                <div className="absolute top-0 right-0 bg-purple-500 text-white text-xs px-2 py-1 rounded-bl-lg font-medium">
                                    Recommended
                                </div>
                                <h4 className="font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                                    <Wand2 size={18} className="text-purple-600" />
                                    {generatedResult.graph.name || 'AI Generated Workflow'}
                                </h4>
                                <p className="text-sm text-gray-600 dark:text-gray-300 mt-2 line-clamp-2">
                                    {generatedResult.graph.description || `Generated from: "${intent}"`}
                                </p>
                                <div className="mt-4 flex gap-4 text-xs text-gray-500">
                                    <span>{Object.keys(generatedResult.graph.nodes).length} Nodes</span>
                                    <span>Auto-layout applied</span>
                                </div>
                            </div>

                            <div className="relative">
                                <div className="absolute inset-0 flex items-center">
                                    <div className="w-full border-t border-gray-200 dark:border-gray-700"></div>
                                </div>
                                <div className="relative flex justify-center text-sm">
                                    <span className="px-2 bg-white dark:bg-gray-900 text-gray-500">Or choose similar template</span>
                                </div>
                            </div>

                            {/* Similar Templates */}
                            <div className="grid grid-cols-2 gap-4">
                                {generatedResult.similar.length > 0 ? generatedResult.similar.map(tpl => (
                                    <div
                                        key={tpl.id}
                                        onClick={() => handleSelectGraph(tpl.graph)}
                                        className="p-4 border border-gray-200 dark:border-gray-700 rounded-xl hover:border-purple-300 cursor-pointer transition-colors"
                                    >
                                        <h5 className="font-medium text-sm">{tpl.name}</h5>
                                    </div>
                                )) : (
                                    <div className="col-span-2 text-center text-gray-400 text-sm py-2">No similar templates found</div>
                                )}
                            </div>

                            <div className="mt-6 flex justify-between">
                                <button onClick={() => setStep(1)} className="text-gray-500 text-sm hover:underline">Back to Intent</button>
                            </div>
                        </div>
                    )}

                    {step === 3 && (
                        <div>Step 3 (Preview) - Logic moved to Builder</div>
                    )}
                </div>
            </div>
        </div>
    );
};
