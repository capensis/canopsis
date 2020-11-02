import moment from 'moment';
import 'moment-duration-format';

import { DATETIME_FORMATS, DEFAULT_DURATION_FORMAT, TIME_UNITS } from '@/constants';

/**
 * Duration filter
 *
 * @param {number | Duration | DurationForm} [duration = 0]
 * @param {string} [format = DEFAULT_DURATION_FORMAT]
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {string}
 */
export default function (duration = 0, format = DEFAULT_DURATION_FORMAT, unit = TIME_UNITS.second) {
  /**
   * TODO: Should be removed after duration refactoring
   */
  const durationValue = duration.seconds || duration.value || duration;
  const resultFormat = DATETIME_FORMATS[format] || format;

  return moment.duration(durationValue, unit).format(resultFormat, { trim: 'both final' }) || '0s';
}
