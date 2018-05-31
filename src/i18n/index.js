import Vue from 'vue';
import VueI18n from 'vue-i18n';

import { DEFAULT_LOCALE } from '@/config';

import localStorageDataSource from '@/services/local-storage-data-source';

import frTranslations from './messages/fr';
import enTranslations from './messages/en';

import frDateFormat from './dateTimeFormats/fr';
import enDateFormat from './dateTimeFormats/en';

Vue.use(VueI18n);

export default new VueI18n({
  locale: localStorageDataSource.getItem('locale') || DEFAULT_LOCALE,
  messages: {
    en: enTranslations,
    fr: frTranslations,
  },
  dateTimeFormats: {
    en: enDateFormat,
    fr: frDateFormat,
  },
});
