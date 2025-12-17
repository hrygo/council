import { Code, Globe, Pencil, Trash2 } from 'lucide-react'; // Removed Bot import
import type { Agent } from '../../../types/agent';


interface AgentCardProps {
    agent: Agent;
    onClick: (agent: Agent) => void;
    onDelete: (agent: Agent) => void;
}

const PROVIDER_COLORS: Record<string, string> = {
    openai: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    anthropic: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-400',
    google: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    deepseek: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
    dashscope: 'bg-yellow-100 text-yellow-700 dark:bg-yellow-900/30 dark:text-yellow-400',
};

export function AgentCard({ agent, onClick, onDelete }: AgentCardProps) {
    const providerColor = PROVIDER_COLORS[agent.model_config.provider] || 'bg-gray-100 text-gray-700';

    return (
        <div
            onClick={() => onClick(agent)}
            className="group relative bg-white dark:bg-gray-800 rounded-xl p-5 shadow-sm border border-gray-200 dark:border-gray-700 hover:shadow-md hover:border-blue-500/50 dark:hover:border-blue-500/50 transition-all duration-200 cursor-pointer"
        >
            <div className="flex items-start justify-between">
                <div className="flex items-center gap-4">
                    <div className="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500/10 to-purple-500/10 flex items-center justify-center text-2xl overflow-hidden">
                        {agent.avatar?.startsWith('http') ? (
                            <img src={agent.avatar} alt={agent.name} className="w-full h-full object-cover" />
                        ) : (
                            <span>{agent.avatar || '🤖'}</span>
                        )}
                    </div>
                    <div>
                        <h3 className="font-semibold text-lg text-gray-900 dark:text-gray-100 truncate max-w-[150px]">
                            {agent.name}
                        </h3>
                        <p className="text-sm text-gray-500 dark:text-gray-400 truncate max-w-[200px]">
                            {agent.persona_prompt.slice(0, 50)}...
                        </p>
                    </div>
                </div>

                <div className="flex items-center gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button
                        onClick={(e) => { e.stopPropagation(); onClick(agent); }}
                        className="p-1.5 text-gray-400 hover:text-blue-500 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-colors"
                    >
                        <Pencil size={14} />
                    </button>
                    <button
                        onClick={(e) => { e.stopPropagation(); onDelete(agent); }}
                        className="p-1.5 text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                    >
                        <Trash2 size={14} />
                    </button>
                </div>
            </div>

            <div className="mt-4 flex flex-wrap items-center gap-2 text-xs">
                <span className={`px-2 py-0.5 rounded-full font-medium ${providerColor}`}>
                    {agent.model_config.provider}
                </span>
                <span className="px-2 py-0.5 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 font-mono">
                    {agent.model_config.model}
                </span>
            </div>

            <div className="mt-3 flex gap-2">
                {agent.capabilities.web_search && (
                    <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-sky-50 dark:bg-sky-900/20 text-sky-600 dark:text-sky-400 text-xs border border-sky-100 dark:border-sky-800">
                        <Globe size={12} />
                        Web
                    </span>
                )}
                {agent.capabilities.code_execution && (
                    <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md bg-amber-50 dark:bg-amber-900/20 text-amber-600 dark:text-amber-400 text-xs border border-amber-100 dark:border-amber-800">
                        <Code size={12} />
                        Code
                    </span>
                )}
            </div>
        </div>
    );
}
