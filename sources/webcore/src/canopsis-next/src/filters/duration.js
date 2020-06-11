import moment from 'moment';
import 'moment-duration-format';

/**
 * Duration filter
 *
 * @param {Number} value - Numeric value to format
 * @param {String} format - Duration format
 * @returns {String}
 */
export default function (value = 0, format = 'D __ H _ m _ s _') {
  const durationValue = value.value ? value.value : value;

  return moment.duration(durationValue, 'seconds').format(format, { trim: 'both final' }) || '0s';
}
