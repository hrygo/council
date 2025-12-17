import { useState, type FC } from 'react';
import { ChevronDown, ChevronRight, Settings } from 'lucide-react';
import { useConfigStore } from '../../../stores/useConfigStore';
import type { ModelConfig } from '../../../types/agent';

interface ModelSelectorProps {
    value: ModelConfig;
    onChange: (config: ModelConfig) => void;
    showAdvanced?: boolean;
}

type ModelProvider = ModelConfig['provider'];

const PROVIDERS: Record<ModelProvider, { name: string; icon: string; models: string[] }> = {
    openai: {
        name: 'OpenAI',
        icon: '🟢',
        models: ['gpt-4o', 'gpt-4o-mini', 'o1-preview', 'o1-mini'],
    },
    anthropic: {
        name: 'Anthropic',
        icon: '🟠',
        models: ['claude-3-5-sonnet-latest', 'claude-3-opus-latest', 'claude-3-haiku-20240307'],
    },
    google: {
        name: 'Google',
        icon: '🔵',
        models: ['gemini-1.5-pro', 'gemini-1.5-flash', 'gemini-2.0-flash-exp'],
    },
    deepseek: {
        name: 'DeepSeek',
        icon: '🟣',
        models: ['deepseek-chat', 'deepseek-reasoner'],
    },
    dashscope: {
        name: 'DashScope',
        icon: '🟡',
        models: ['qwen-max', 'qwen-plus', 'qwen-turbo'],
    },
};

export const ModelSelector: FC<ModelSelectorProps> = ({ value, onChange, showAdvanced }) => {
    const { godMode } = useConfigStore();
    const showParams = showAdvanced || godMode;
    const [isExpanded, setIsExpanded] = useState(false);

    const handleProviderChange = (provider: ModelProvider) => {
        onChange({
            ...value,
            provider,
            model: PROVIDERS[provider].models[0],
        });
    };

    return (
        <div className="space-y-4 p-4 bg-gray-50 dark:bg-gray-800/50 rounded-xl border border-gray-100 dark:border-gray-800">
            <div className="flex items-center gap-2 mb-2 text-sm font-semibold text-gray-700 dark:text-gray-300">
                <Settings size={16} />
                <h3>Model Configuration</h3>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                    <label className="text-xs font-medium text-gray-500 uppercase">Provider</label>
                    <div className="relative">
                        <select
                            value={value.provider}
                            onChange={(e) => handleProviderChange(e.target.value as ModelProvider)}
                            className="w-full appearance-none px-3 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all text-sm"
                        >
                            {Object.entries(PROVIDERS).map(([key, p]) => (
                                <option key={key} value={key}>
                                    {p.icon} {p.name}
                                </option>
                            ))}
                        </select>
                        <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none" size={14} />
                    </div>
                </div>

                <div className="space-y-2">
                    <label className="text-xs font-medium text-gray-500 uppercase">Model</label>
                    <div className="relative">
                        <select
                            value={value.model}
                            onChange={(e) => onChange({ ...value, model: e.target.value })}
                            className="w-full appearance-none px-3 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 transition-all text-sm"
                        >
                            {PROVIDERS[value.provider as ModelProvider]?.models.map((m) => (
                                <option key={m} value={m}>
                                    {m}
                                </option>
                            ))}
                        </select>
                        <ChevronDown className="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 pointer-events-none" size={14} />
                    </div>
                </div>
            </div>

            {showParams && (
                <div className="border-t border-gray-200 dark:border-gray-700 pt-3">
                    <button
                        type="button"
                        onClick={() => setIsExpanded(!isExpanded)}
                        className="flex items-center gap-1 text-xs font-medium text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 transition-colors mb-4"
                    >
                        {isExpanded ? <ChevronDown size={14} /> : <ChevronRight size={14} />}
                        Advanced Parameters
                    </button>

                    {isExpanded && (
                        <div className="space-y-4 animate-in fade-in slide-in-from-top-2 duration-200">
                            <div className="space-y-2">
                                <div className="flex justify-between text-xs">
                                    <label className="text-gray-600 dark:text-gray-400">Temperature Details</label>
                                    <span className="font-mono text-gray-900 dark:text-gray-100">{value.temperature}</span>
                                </div>
                                <input
                                    type="range"
                                    min="0"
                                    max="2"
                                    step="0.1"
                                    value={value.temperature}
                                    onChange={(e) => onChange({ ...value, temperature: parseFloat(e.target.value) })}
                                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer"
                                />
                                <div className="flex justify-between text-[10px] text-gray-400">
                                    <span>Precise (0.0)</span>
                                    <span>Balanced (1.0)</span>
                                    <span>Creative (2.0)</span>
                                </div>
                            </div>

                            <div className="space-y-2">
                                <div className="flex justify-between text-xs">
                                    <label className="text-gray-600 dark:text-gray-400">Top P</label>
                                    <span className="font-mono text-gray-900 dark:text-gray-100">{value.top_p}</span>
                                </div>
                                <input
                                    type="range"
                                    min="0"
                                    max="1"
                                    step="0.05"
                                    value={value.top_p}
                                    onChange={(e) => onChange({ ...value, top_p: parseFloat(e.target.value) })}
                                    className="w-full h-1.5 bg-gray-200 dark:bg-gray-700 rounded-lg appearance-none cursor-pointer"
                                />
                            </div>

                            <div className="space-y-2">
                                <div className="flex justify-between text-xs">
                                    <label className="text-gray-600 dark:text-gray-400">Max Tokens</label>
                                    <span className="font-mono text-gray-900 dark:text-gray-100">{value.max_tokens}</span>
                                </div>
                                <input
                                    type="number"
                                    min="100"
                                    max="128000"
                                    step="100"
                                    value={value.max_tokens}
                                    onChange={(e) => onChange({ ...value, max_tokens: parseInt(e.target.value) || 4096 })}
                                    className="w-full px-3 py-1.5 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg text-sm text-center"
                                />
                            </div>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};
