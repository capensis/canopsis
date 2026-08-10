import { merge } from 'lodash';
import { createI18n } from 'vue-i18n';

import { DEFAULT_LOCALE } from '@/config';
import { updateDateLocaleMessages } from '@/helpers/date/date';

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
  messages: merge(
    {},
    messages,
    // TODO: Uncomment this when we have a features service
    // featuresService.get('i18n')
  ),
});
