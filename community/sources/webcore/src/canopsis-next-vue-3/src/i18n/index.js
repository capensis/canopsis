import { merge } from 'lodash';
import { createI18n } from 'vue-i18n';

import { DEFAULT_LOCALE } from '@/config';

import { featuresService } from '@/services/features';

import { updateDateLocaleMessages } from '@/helpers/date/date';

import { FEATURES_SERVICE_KEYS } from '@/constants/features';

import durationFrMessages from './duration-fr';
import messages from './messages';

updateDateLocaleMessages('fr', durationFrMessages);

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: DEFAULT_LOCALE,
  fallbackLocale: DEFAULT_LOCALE,
  warnHtmlMessage: false,
  warnHtmlInMessage: 'off',
  messages: merge({}, messages, featuresService.get(FEATURES_SERVICE_KEYS.i18n, {})),
});
