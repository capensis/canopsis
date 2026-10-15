import bundledEn from '@vee-validate/i18n/dist/locale/en.json';
import bundledFr from '@vee-validate/i18n/dist/locale/fr.json';

import enValidation from '@/i18n/messages/en/validation';
import frValidation from '@/i18n/messages/fr/validation';

const mergeLocale = (locale, messages) => ({
  ...locale,
  messages: {
    ...locale.messages,
    ...messages,
  },
});

export const veeValidateLocales = {
  en: mergeLocale(bundledEn, enValidation.messages),
  fr: mergeLocale(bundledFr, frValidation.messages),
};
