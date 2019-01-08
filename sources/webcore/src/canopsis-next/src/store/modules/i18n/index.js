import Vue from 'vue';

import i18n from '@/i18n';
import localStorageDataSource from '@/services/local-storage-data-source';

export default {
  namespaced: true,
  actions: {
    setLocale(context, locale) {
      Vue.set(i18n, 'locale', locale);
      localStorageDataSource.setItem('locale', locale);
    },
  },
};
