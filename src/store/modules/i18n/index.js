import Vue from 'vue';

import localStorageDataSource from '@/services/local-storage-data-source';
import i18n from '@/i18n';

export default {
  namespaced: true,
  actions: {
    setLocale(context, locale) {
      Vue.set(i18n, 'locale', locale);
      localStorageDataSource.setItem('locale', locale);
    },
  },
};
