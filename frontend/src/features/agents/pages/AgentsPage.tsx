import { AgentList } from '../components/AgentList';
import { Network } from 'lucide-react'; // Factory icon or network

export function AgentsPage() {
    return (
        <div className="container mx-auto max-w-7xl px-4 py-8">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-gray-900 dark:text-gray-100 flex items-center gap-3">
                    <Network size={32} className="text-blue-600" />
                    Agent Factory
                </h1>
                <p className="mt-2 text-gray-600 dark:text-gray-400">
                    Design and commission specialized AI agents for your council.
                </p>
            </div>

            <AgentList />
        </div>
    );
}
