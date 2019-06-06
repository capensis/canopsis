import moment from 'moment';
import 'moment-duration-format';
import momentDurationFrLocale from '@/i18n/moment-duration-fr';

import { DEFAULT_LOCALE } from '@/config';

moment.updateLocale('fr', momentDurationFrLocale);

/**
 *
 * @param {Number} value - Numeric value to format
 * @param {String} format - Duration format
 *
 * @returns {String}
 */
export default function ({
  value = 0,
  locale = DEFAULT_LOCALE,
  format = 'D __ H _ m _ s _',
}) {
  moment.locale(locale);
  return moment.duration(value, 'seconds').format(format, { trim: 'both final' }) || 0;
}
