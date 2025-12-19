import { useEffect } from 'react';
import { Routes, Route, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Play, Boxes, Users, Network, type LucideIcon } from 'lucide-react';
import './i18n';
import './index.css';
import { useConfigStore } from './stores/useConfigStore';
import { MeetingRoom } from './features/meeting/MeetingRoom';
import { WorkflowEditor } from './features/editor/WorkflowEditor';
import { GroupsPage } from './features/groups/pages/GroupsPage';
import { AgentsPage } from './features/agents/pages/AgentsPage';
import { HomePage } from './features/home/HomePage';
import { LanguageSwitcher } from './components/LanguageSwitcher';
import { ToastProvider } from './components/ui/Toast';

// Simple Nav Bar for demo purposes
const NavButton = ({ path, icon: Icon, label, onClick }: { path: string, icon: LucideIcon, label: string, onClick: (path: string) => void }) => (
  <button
    onClick={() => onClick(path)}
    className="flex flex-col items-center justify-center p-2 text-gray-500 hover:text-blue-600 hover:bg-blue-50 dark:hover:bg-blue-900/20 rounded-lg transition-all gap-1 group select-none"
    title={label}
  >
    <Icon size={20} className="group-hover:scale-110 transition-transform" />
    <span className="text-[10px] font-medium">{label}</span>
  </button>
);

function App() {
  const { t } = useTranslation();
  const theme = useConfigStore((state) => state.theme);
  const navigate = useNavigate();

  useEffect(() => {
    document.documentElement.className = theme === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : '')
      : theme;
  }, [theme]);

  return (
    <ToastProvider>
      <div className="h-screen w-screen overflow-hidden bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
        {/* Sidebar / Nav */}
        <div className="fixed left-0 top-0 h-full w-16 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col items-center py-4 gap-4 z-50 shadow-sm">
          <div
            onClick={() => navigate('/')}
            className="mb-2 p-2 bg-blue-600 rounded-xl text-white shadow-lg shadow-blue-500/30 cursor-pointer hover:scale-110 transition-transform"
          >
            <Boxes size={24} />
          </div>

          <NavButton path="/chat" icon={Play} label={t('chat.title')} onClick={navigate} />
          <NavButton path="/editor" icon={Boxes} label={t('workflow.builder.title', { defaultValue: 'Builder' })} onClick={navigate} />
          <NavButton path="/groups" icon={Users} label={t('nav.groups')} onClick={navigate} />
          <NavButton path="/agents" icon={Network} label={t('nav.agents')} onClick={navigate} />

          {/* Spacer */}
          <div className="flex-1" />

          {/* Bottom section */}
          <div className="flex flex-col items-center gap-2 mb-2">
            <LanguageSwitcher />
          </div>
        </div>

        <div className="ml-16 h-full w-[calc(100vw-4rem)]">

          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/chat" element={<MeetingRoom />} />
            <Route path="/editor" element={<WorkflowEditor />} />
            <Route path="/groups" element={<GroupsPage />} />
            <Route path="/agents" element={<AgentsPage />} />
          </Routes>
        </div>
      </div>
    </ToastProvider>
  );
}

export default App;
