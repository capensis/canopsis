import Vue, { computed, getCurrentInstance } from 'vue';
import VueI18n from 'vue-i18n';
import { merge } from 'lodash';

import { DEFAULT_LOCALE } from '@/config';

import featuresService from '@/services/features';

import { updateDateLocaleMessages } from '@/helpers/date/date';

import durationFrMessages from './duration-fr';
import messages from './messages';

updateDateLocaleMessages('fr', durationFrMessages);

Vue.use(VueI18n);

const i18nInstance = new VueI18n({
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  messages: merge(messages, featuresService.get('i18n')),
});

/**
 * USe I18n hoo for composition API
 *
 * @return {Object}
 */
export const useI18n = () => {
  if (!i18nInstance) throw new Error('vue-i18n not initialized');

  const instance = getCurrentInstance();
  const vm = instance?.proxy || instance || new Vue({});

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

export default i18nInstance;
