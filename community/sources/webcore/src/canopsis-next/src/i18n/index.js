import Vue from 'vue';
import VueI18n from 'vue-i18n';
import { merge } from 'lodash';

import { DEFAULT_LOCALE } from '@/config';

import featuresService from '@/services/features';

import { updateDateLocaleMessages } from '@/helpers/date/date';

import durationFrMessages from './duration-fr';
import messages from './messages';

updateDateLocaleMessages('fr', durationFrMessages);

Vue.use(VueI18n);

export default new VueI18n({
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  messages: merge(messages, featuresService.get('i18n')),
});
