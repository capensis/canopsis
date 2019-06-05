import moment from 'moment';
import i18n from '@/i18n';
import 'moment-duration-format';

/**
 *
 * @param {Number} value - Numeric value to format
 * @param {String} format - Duration format
 *
 * @returns {String}
 */
export default function (
  value,
  format = `D [${i18n.tc('common.times.day', 'D')}] H[h] m[m] s[s]`,
) {
  return moment.duration(value, 'seconds').format(format, { trim: 'both final' }) || 0;
}
