import { useState } from 'react';
import { Plus, Bot, Search } from 'lucide-react';
import { useAgents, useDeleteAgent } from '../../../hooks/useAgents';
import { AgentCard } from './AgentCard';
import { AgentEditDrawer } from './AgentEditDrawer';
import type { Agent } from '../../../types/agent';

export function AgentList() {
    const { data: agents, isLoading, error } = useAgents();
    const deleteAgent = useDeleteAgent();

    const [isDrawerOpen, setIsDrawerOpen] = useState(false);
    const [editingAgent, setEditingAgent] = useState<Agent | null>(null);
    const [searchQuery, setSearchQuery] = useState('');

    const openCreate = () => {
        setEditingAgent(null);
        setIsDrawerOpen(true);
    };

    const openEdit = (agent: Agent) => {
        setEditingAgent(agent);
        setIsDrawerOpen(true);
    };

    const handleDelete = (agent: Agent) => {
        if (confirm(`Are you sure you want to retire agent "${agent.name}"?`)) {
            deleteAgent.mutate(agent.agent_uuid);
        }
    };

    const filteredAgents = agents?.filter((a: Agent) =>
        a.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        a.persona_prompt.toLowerCase().includes(searchQuery.toLowerCase())
    );

    if (isLoading) {
        return (
            <div className="flex items-center justify-center p-12">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="p-4 bg-red-50 text-red-600 rounded-lg">
                Error loading agents: {error.message}
            </div>
        );
    }

    return (
        <div className="space-y-6">
            {/* Toolbar */}
            <div className="flex items-center justify-between gap-4">
                <div className="relative flex-1 max-w-md">
                    <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" size={18} />
                    <input
                        type="text"
                        placeholder="Search agents..."
                        value={searchQuery}
                        onChange={(e) => setSearchQuery(e.target.value)}
                        className="w-full pl-10 pr-4 py-2 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50"
                    />
                </div>
                <button
                    onClick={openCreate}
                    className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
                >
                    <Plus size={18} />
                    Hire Agent
                </button>
            </div>

            {/* Grid */}
            {filteredAgents?.length === 0 ? (
                <div className="text-center py-16 bg-white dark:bg-gray-800/50 rounded-2xl border border-dashed border-gray-200 dark:border-gray-700">
                    <div className="mx-auto w-16 h-16 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mb-4 text-gray-400">
                        <Bot size={32} />
                    </div>
                    <h3 className="text-lg font-medium text-gray-900 dark:text-gray-100">No agents found</h3>
                    <p className="text-gray-500 mt-1 mb-6">Create specialized AI agents to join your council.</p>
                    <button
                        onClick={openCreate}
                        className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
                    >
                        Create Agent
                    </button>
                </div>
            ) : (
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                    {filteredAgents?.map((agent: Agent) => (
                        <AgentCard
                            key={agent.agent_uuid}
                            agent={agent}
                            onClick={openEdit}
                            onDelete={handleDelete}
                        />
                    ))}
                </div>
            )}

            {isDrawerOpen && (
                <AgentEditDrawer
                    open={true}
                    onClose={() => setIsDrawerOpen(false)}
                    agent={editingAgent}
                />
            )}
        </div>
    );
}
