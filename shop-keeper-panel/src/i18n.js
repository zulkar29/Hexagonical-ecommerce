import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import enCommon from './locales/en/common.json';
import bnCommon from './locales/bn/common.json';

i18n
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: enCommon },
      bn: { translation: bnCommon },
    },
    lng: 'en',
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false
    }
  });

export default i18n; 