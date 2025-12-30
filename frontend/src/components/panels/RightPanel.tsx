import { type FC, useState } from 'react';
import { Database, FolderCode, TrendingUp, Minimize2 } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { KnowledgePanel } from '../../features/meeting-room/components/KnowledgePanel';
import { VFSExplorer } from '../vfs/VFSExplorer';
import { LoopAnalytics } from '../analytics/LoopAnalytics';
import { useSessionStore } from '../../stores/useSessionStore';

interface RightPanelProps {
    sessionId: string;
    fullscreen?: boolean;
    onExitFullscreen?: () => void;
}

type TabKey = 'knowledge' | 'codebase' | 'analytics';

interface TabConfig {
    key: TabKey;
    labelKey: string;
    icon: typeof Database;
}

const tabs: TabConfig[] = [
    { key: 'knowledge', labelKey: 'rightPanel.tabs.knowledge', icon: Database },
    { key: 'codebase', labelKey: 'rightPanel.tabs.codebase', icon: FolderCode },
    { key: 'analytics', labelKey: 'rightPanel.tabs.analytics', icon: TrendingUp },
];

export const RightPanel: FC<RightPanelProps> = ({ sessionId, fullscreen, onExitFullscreen }) => {
    const { t } = useTranslation();
    const [activeTab, setActiveTab] = useState<TabKey>('knowledge');

    // Get score history from session context (if available)
    const contextData = useSessionStore(state => state.currentSession?.contextData);
    const scoreHistory = contextData?.score_history || [];

    return (
        <div className="h-full flex flex-col bg-white dark:bg-gray-900 border-l border-gray-200 dark:border-gray-700">
            {/* Tab Header */}
            <div className="flex border-b border-gray-200 dark:border-gray-700 relative">
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
                            <span className="hidden sm:inline">{t(tab.labelKey)}</span>
                        </button>
                    );
                })}

                {/* Exit Fullscreen Button */}
                {fullscreen && onExitFullscreen && (
                    <button
                        onClick={onExitFullscreen}
                        className="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 bg-gray-100 dark:bg-gray-700 rounded hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
                        title={t('common.actions.close')}
                    >
                        <Minimize2 size={14} />
                    </button>
                )}
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
