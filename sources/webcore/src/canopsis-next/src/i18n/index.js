import Vue from 'vue';
import VueI18n from 'vue-i18n';
import moment from 'moment';
import 'moment/locale/fr';

import { DEFAULT_LOCALE } from '@/config';

import frTranslations from './messages/fr';
import enTranslations from './messages/en';
import momentDurationFrLocale from './moment-duration-fr';

const currentLocale = moment.locale();

moment.updateLocale('fr', momentDurationFrLocale);
moment.locale(currentLocale);

Vue.use(VueI18n);

export default new VueI18n({
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  messages: {
    en: enTranslations,
    fr: frTranslations,
  },
});
