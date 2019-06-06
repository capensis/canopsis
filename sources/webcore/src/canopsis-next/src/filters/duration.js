import moment from 'moment';
import 'moment-duration-format';

/**
 *
 * @param {Number} value - Numeric value to format
 * @param {String} format - Duration format
 *
 * @returns {String}
 */
export default function (value, format = 'HH[h] mm[m] ss[s]') {
  return moment.duration(value, 'seconds').format(format, { trim: 'both final' }) || '0s';
}
