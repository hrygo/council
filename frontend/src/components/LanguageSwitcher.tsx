import { useTranslation } from 'react-i18next';
import { Globe } from 'lucide-react';

const languages = [
    { code: 'en', label: 'EN' },
    { code: 'zh-CN', label: '中文' },
];

export function LanguageSwitcher() {
    const { i18n } = useTranslation();

    const handleChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
        i18n.changeLanguage(e.target.value);
    };

    return (
        <div className="relative inline-flex items-center gap-1">
            <Globe className="w-4 h-4 text-slate-500" />
            <select
                value={i18n.language}
                onChange={handleChange}
                className="appearance-none bg-transparent text-sm font-medium text-slate-700 dark:text-slate-300 cursor-pointer focus:outline-none pr-4"
                aria-label="Select language"
            >
                {languages.map(({ code, label }) => (
                    <option key={code} value={code}>
                        {label}
                    </option>
                ))}
            </select>
            <span className="pointer-events-none absolute right-0 text-slate-400 text-xs">▼</span>
        </div>
    );
}
