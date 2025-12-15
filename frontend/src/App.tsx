import { useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import './i18n';
import './index.css';
import { useConfigStore } from './stores/useConfigStore';

function App() {
  const { t } = useTranslation();
  const theme = useConfigStore((state) => state.theme);

  useEffect(() => {
    document.documentElement.className = theme === 'system'
      ? (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : '')
      : theme;
  }, [theme]);

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900 text-gray-900 dark:text-gray-100 flex items-center justify-center">
      <div className="p-8 bg-white dark:bg-gray-800 rounded shadow-lg">
        <h1 className="text-3xl font-bold mb-4">{t('welcome')}</h1>
        <p className="text-gray-600 dark:text-gray-400">
          Environment: {import.meta.env.MODE}
        </p>
      </div>
    </div>
  );
}

export default App;
