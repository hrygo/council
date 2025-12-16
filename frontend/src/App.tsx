import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';
import './i18n';
import './index.css';
import { useConfigStore } from './stores/useConfigStore';
import { MeetingRoom } from './components/layout/MeetingRoom';
import { WorkflowEditor } from './components/layout/WorkflowEditor';

function App() {
  // const { t } = useTranslation();
  const theme = useConfigStore((state) => state.theme);
  const [mode, setMode] = useState<'build' | 'run'>('run');

  useEffect(() => {
    document.documentElement.className = theme === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : '')
      : theme;
  }, [theme]);

  return (
    <div className="h-screen w-screen overflow-hidden bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100">
      {/* Dev Mode Switcher */}
      <div className="fixed bottom-4 left-4 z-50 flex gap-2 p-1 bg-white/50 dark:bg-black/50 backdrop-blur rounded shadow">
        <button onClick={() => setMode('build')} className={`px-3 py-1 text-xs font-medium rounded transition-colors ${mode === 'build' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}>Build</button>
        <button onClick={() => setMode('run')} className={`px-3 py-1 text-xs font-medium rounded transition-colors ${mode === 'run' ? 'bg-green-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}>Run</button>
      </div>

      {mode === 'run' ? <MeetingRoom /> : <WorkflowEditor />}
    </div>
  );
}

export default App;
