import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import enTranslation from './locales/en/translation.json';
import zhCNTranslation from './locales/zh-CN/translation.json';

const resources = {
    en: {
        translation: enTranslation,
    },
    'zh-CN': {
        translation: zhCNTranslation,
    },
};

i18n
    .use(LanguageDetector)
    .use(initReactI18next)
    .init({
        resources,
        fallbackLng: 'en',
        supportedLngs: ['en', 'zh-CN'],
        detection: {
            order: ['localStorage', 'navigator', 'htmlTag'],
            caches: ['localStorage'],
            lookupLocalStorage: 'council-language',
        },
        interpolation: {
            escapeValue: false,
        },
    });

export default i18n;
