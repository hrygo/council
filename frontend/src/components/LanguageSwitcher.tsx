import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';

const languages = [
    { code: 'en', label: 'EN' },
    { code: 'zh-CN', label: '中' },
];

export function LanguageSwitcher() {
    const { i18n } = useTranslation();

    // 获取当前语言的显示标签
    const currentLang = languages.find(l => l.code === i18n.language) || languages[0];

    // 一键切换到下一个语言
    const toggleLanguage = () => {
        const currentIndex = languages.findIndex(l => l.code === i18n.language);
        const nextIndex = (currentIndex + 1) % languages.length;
        i18n.changeLanguage(languages[nextIndex].code);
    };

    return (
        <button
            onClick={toggleLanguage}
            className="flex items-center gap-1 px-2 py-1.5 text-sm font-medium text-slate-600 dark:text-slate-400 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-all"
            aria-label="切换语言 / Toggle language"
            title={`当前: ${currentLang.label} (点击切换)`}
        >
            <Globe className="w-4 h-4" />
            <span className="min-w-[2ch]">{currentLang.label}</span>
        </button>
    );
}
