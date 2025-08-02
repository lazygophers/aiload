import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

// 导入您的翻译文件
import enTranslation from './locales/en/translation.json';
import zhTranslation from './locales/zh/translation.json';
import zhtwTranslation from './locales/zh-tw/translation.json';
import jaTranslation from './locales/ja/translation.json';
import koTranslation from './locales/ko/translation.json';
import frTranslation from './locales/fr/translation.json';
import deTranslation from './locales/de/translation.json';
import esTranslation from './locales/es/translation.json';

const resources = {
  en: {
    translation: enTranslation,
  },
  zh: {
    translation: zhTranslation,
  },
  'zh-TW': {
    translation: zhtwTranslation,
  },
  ja: {
    translation: jaTranslation,
  },
  ko: {
    translation: koTranslation,
  },
  fr: {
    translation: frTranslation,
  },
  de: {
    translation: deTranslation,
  },
  es: {
    translation: esTranslation,
  },
};

i18n
  .use(LanguageDetector) // 自动检测浏览器语言
  .use(initReactI18next) // 将 i18n 实例传递给 react-i18next
  .init({
    resources,
    fallbackLng: 'en', // 如果检测不到语言，则默认使用英语
    interpolation: {
      escapeValue: false, // React 已经可以防止 XSS
    },
  });

export default i18n;