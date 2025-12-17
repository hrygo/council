import { useEffect } from 'react';
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom';
import { Play, Boxes, Users, Network } from 'lucide-react';
import './i18n';
import './index.css';
import { useConfigStore } from './stores/useConfigStore';
import { MeetingRoom } from './components/layout/MeetingRoom';
import { WorkflowEditor } from './components/layout/WorkflowEditor';
import { GroupsPage } from './features/groups/pages/GroupsPage';
import { AgentsPage } from './features/agents/pages/AgentsPage';

function App() {
  const theme = useConfigStore((state) => state.theme);
  const navigate = useNavigate();
  const location = useLocation();

  useEffect(() => {
    document.documentElement.className = theme === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : '')
      : theme;
  }, [theme]);

  // Simple Nav Bar for demo purposes
  const NavButton = ({ path, icon: Icon, label }: { path: string, icon: any, label: string }) => (
    <button
      onClick={() => navigate(path)}
      className={`p-2 rounded-lg transition-colors ${location.pathname === path
        ? 'bg-blue-600 text-white'
        : 'bg-white/50 dark:bg-black/50 hover:bg-gray-200 dark:hover:bg-gray-700'
        }`}
      title={label}
    >
      <Icon size={20} />
    </button>
  );

  return (
    <div className="h-screen w-screen overflow-hidden bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      {/* Navigation Bar */}
      <div className="fixed bottom-4 left-4 z-50 flex gap-2 p-1 rounded shadow-lg">
        <NavButton path="/" icon={Play} label="Run Mode" />
        <NavButton path="/editor" icon={Boxes} label="Builder Mode" />
        <NavButton path="/groups" icon={Users} label="Groups Management" />
        <NavButton path="/agents" icon={Network} label="Agent Factory" />
      </div>

      <Routes>
        <Route path="/" element={<MeetingRoom />} />
        <Route path="/editor" element={<WorkflowEditor />} />
        <Route path="/groups" element={<GroupsPage />} />
        <Route path="/agents" element={<AgentsPage />} />
      </Routes>
    </div>
  );
}

export default App;
