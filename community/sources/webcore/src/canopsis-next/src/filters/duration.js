import moment from 'moment';
import { isNil } from 'lodash';

import 'moment-duration-format';

import { DATETIME_FORMATS, DEFAULT_DURATION_FORMAT, TIME_UNITS } from '@/constants';

/**
 * Duration filter
 *
 * @param {number | Duration | DurationForm} duration
 * @param {string} [format = DEFAULT_DURATION_FORMAT]
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {string}
 */
export default function (duration, format = DEFAULT_DURATION_FORMAT, unit = TIME_UNITS.second) {
  if (isNil(duration)) {
    return '';
  }

  /**
   * TODO: Should be removed after duration refactoring
   */
  const durationValue = duration ? (duration.seconds || duration.value || duration) : duration;
  const resultFormat = DATETIME_FORMATS[format] || format;

  return moment.duration(durationValue, unit).format(resultFormat, { trim: 'both final' }) || '0s';
}
