import { computed, unref } from 'vue';

import { i18n } from '@/i18n';

export const useI18n = () => ({
  locale: computed(() => unref(i18n.global.locale)),
  t: i18n.global.t.bind(i18n.global),
  te: i18n.global.te.bind(i18n.global),
  tm: i18n.global.tm.bind(i18n.global),
});
