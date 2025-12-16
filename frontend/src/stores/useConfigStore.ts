import { create } from 'zustand';
import { persist, createJSONStorage } from 'zustand/middleware';
import i18n from '../i18n'; // Ensure this path is correct

interface ConfigState {
    theme: 'light' | 'dark' | 'system';
    language: 'zh-CN' | 'en-US';
    godMode: boolean;

    setTheme: (theme: 'light' | 'dark' | 'system') => void;
    setLanguage: (lang: 'zh-CN' | 'en-US') => void;
    toggleGodMode: () => void;
}

export const useConfigStore = create<ConfigState>()(
    persist(
        (set) => ({
            theme: 'system',
            language: (localStorage.getItem('i18nextLng') as 'zh-CN' | 'en-US') || 'zh-CN',
            godMode: false,

            setTheme: (theme) => set({ theme }),
            setLanguage: (lang) => {
                i18n.changeLanguage(lang);
                set({ language: lang });
            },
            toggleGodMode: () => set((state) => ({ godMode: !state.godMode })),
        }),
        {
            name: 'council-config',
            storage: createJSONStorage(() => localStorage),
        }
    )
);
