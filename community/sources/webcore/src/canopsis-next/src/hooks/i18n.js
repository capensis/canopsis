import Vue, { computed } from 'vue';

import i18nInstance from '@/i18n';

import { useComponentInstance } from '@/hooks/vue';

/**
 * USe I18n hoo for composition API
 *
 * @return {Object}
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
