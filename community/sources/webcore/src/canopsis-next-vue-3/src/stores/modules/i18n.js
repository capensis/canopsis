import { defineStore } from 'pinia';
import { ref } from 'vue';

import { DEFAULT_LOCALE, LOCALE_PRIORITIES } from '@/config';
import { setDateLocale } from '@/helpers/date/date';
import { i18n } from '@/i18n';

export const useI18nStore = defineStore('i18n', () => {
  const localePriority = ref(0);

  /**
   * Set application locale if priority is high enough.
   *
   * @param {Object} payload
   * @param {string} [payload.locale]
   * @param {number} [payload.priority]
   */
  const setLocale = ({ locale = DEFAULT_LOCALE, priority = LOCALE_PRIORITIES.default } = {}) => {
    if (localePriority.value > priority) {
      return;
    }

    setDateLocale(locale);

    i18n.global.locale.value = locale;

    localePriority.value = priority;
  };

  /**
   * Set default locale.
   *
   * @param {string} [locale]
   */
  const setDefaultLocale = (locale = DEFAULT_LOCALE) =>
    setLocale({
      locale,
      priority: LOCALE_PRIORITIES.default,
    });

  /**
   * Set global locale.
   *
   * @param {string} locale
   */
  const setGlobalLocale = locale =>
    setLocale({
      locale,
      priority: LOCALE_PRIORITIES.global,
    });

  /**
   * Set personal locale.
   *
   * @param {string} locale
   */
  const setPersonalLocale = locale =>
    setLocale({
      locale,
      priority: LOCALE_PRIORITIES.personal,
    });

  return {
    localePriority,

    setLocale,
    setDefaultLocale,
    setGlobalLocale,
    setPersonalLocale,
  };
});
