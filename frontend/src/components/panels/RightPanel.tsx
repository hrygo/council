import { type FC, useState } from 'react';
import { Database, FolderCode, TrendingUp } from 'lucide-react';
import { KnowledgePanel } from '../../features/meeting-room/components/KnowledgePanel';
import { VFSExplorer } from '../vfs/VFSExplorer';
import { LoopAnalytics } from '../analytics/LoopAnalytics';
import { useSessionStore } from '../../stores/useSessionStore';

interface RightPanelProps {
    sessionId: string;
}

type TabKey = 'knowledge' | 'codebase' | 'analytics';

interface TabConfig {
    key: TabKey;
    label: string;
    icon: typeof Database;
}

const tabs: TabConfig[] = [
    { key: 'knowledge', label: 'Knowledge', icon: Database },
    { key: 'codebase', label: 'Codebase', icon: FolderCode },
    { key: 'analytics', label: 'Analytics', icon: TrendingUp },
];

export const RightPanel: FC<RightPanelProps> = ({ sessionId }) => {
    const [activeTab, setActiveTab] = useState<TabKey>('knowledge');

    // Get score history from session context (if available)
    const contextData = useSessionStore(state => state.currentSession?.contextData);
    const scoreHistory = contextData?.score_history || [];

    return (
        <div className="h-full flex flex-col bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-700">
            {/* Tab Header */}
            <div className="flex border-b border-gray-200 dark:border-gray-700">
                {tabs.map((tab) => {
                    const Icon = tab.icon;
                    const isActive = activeTab === tab.key;
                    return (
                        <button
                            key={tab.key}
                            onClick={() => setActiveTab(tab.key)}
                            className={`flex-1 flex items-center justify-center gap-1.5 px-3 py-2.5 text-xs font-medium transition-colors ${isActive
                                    ? 'text-blue-600 dark:text-blue-400 border-b-2 border-blue-600 dark:border-blue-400 bg-blue-50/50 dark:bg-blue-900/20'
                                    : 'text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-800'
                                }`}
                        >
                            <Icon size={14} />
                            <span className="hidden sm:inline">{tab.label}</span>
                        </button>
                    );
                })}
            </div>

            {/* Tab Content */}
            <div className="flex-1 overflow-hidden">
                {activeTab === 'knowledge' && (
                    <KnowledgePanel sessionId={sessionId} />
                )}
                {activeTab === 'codebase' && (
                    <VFSExplorer sessionId={sessionId} />
                )}
                {activeTab === 'analytics' && (
                    <LoopAnalytics scoreHistory={scoreHistory} threshold={90} />
                )}
            </div>
        </div>
    );
};
