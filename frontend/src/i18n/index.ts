import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

// In a real app, these would be imported from JSON files
const resources = {
    'zh-CN': {
        translation: {
            "welcome": "欢迎来到理事会",
        }
    },
    'en-US': {
        translation: {
            "welcome": "Welcome to The Council",
        }
    }
};

i18n
    .use(initReactI18next)
    .init({
        resources,
        lng: 'zh-CN',
        fallbackLng: 'en-US',
        interpolation: {
            escapeValue: false,
        },
    });

export default i18n;
