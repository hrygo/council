import { create } from 'zustand';
import { persist } from 'zustand/middleware';

interface ConfigState {
    theme: 'light' | 'dark' | 'system';
    language: 'zh-CN' | 'en-US';
    setTheme: (theme: 'light' | 'dark' | 'system') => void;
    setLanguage: (lang: 'zh-CN' | 'en-US') => void;
}

export const useConfigStore = create<ConfigState>()(
    persist(
        (set) => ({
            theme: 'system',
            language: 'zh-CN',
            setTheme: (theme) => set({ theme }),
            setLanguage: (language) => set({ language }),
        }),
        {
            name: 'council-config',
        }
    )
);
