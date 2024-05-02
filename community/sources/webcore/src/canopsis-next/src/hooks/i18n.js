import Vue, { computed } from 'vue';

import i18nInstance from '@/i18n';

import { useComponentInstance } from '@/hooks/vue';

/**
 * Provides a reactive i18n interface for handling internationalization within Vue components.
 * This hook utilizes the `useComponentInstance` to access the current Vue instance and sets up
 * a reactive `locale` property. It also binds the translation methods `$t`, `$tc`, and `$te` from
 * the Vue instance to ensure they can be used reactively in the component.
 *
 * @returns {Object} An object containing:
 * - `locale`: a reactive property for getting and setting the current locale.
 * - `t`: a method bound to the current Vue instance for translating text.
 * - `tc`: a method bound to the current Vue instance for translating text with pluralization.
 * - `te`: a method bound to the current Vue instance for checking if a translation exists.
 */
export const useI18n = () => {
  const instance = useComponentInstance();
  const vm = instance || new Vue({});

  const locale = computed({
    get() {
      return i18nInstance.locale;
    },
    set(value) {
      i18nInstance.locale = value;
    },
  });

  return {
    locale,
    t: vm.$t.bind(vm),
    tc: vm.$tc.bind(vm),
    te: vm.$te.bind(vm),
  };
};
