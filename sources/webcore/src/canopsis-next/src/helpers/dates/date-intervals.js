import moment from 'moment-timezone';

import { STATS_DURATION_UNITS, DATETIME_FORMATS, STATS_QUICK_RANGES } from '@/constants';

/**
 * Convert a date interval string to moment date object
 *
 * @param {String} dateString
 * @param {String} type
 * @returns {Moment}
 */
export function parseStringToDateInterval(dateString, type) {
  const preparedDateString = dateString.toLowerCase().trim();
  const matches = preparedDateString.match(/^now(([+--])(\d+)([hdwmy]{1}))?(\/([hdwmy]{1}))?$/);

  if (matches) {
    const result = moment().utc();
    const operator = matches[2];
    const deltaValue = matches[3];
    let roundUnit = matches[6];
    let deltaUnit = matches[4];

    const roundMethod = type === 'start' ? 'startOf' : 'endOf';
    const deltaMethod = operator === '+' ? 'add' : 'subtract';

    if (roundUnit) {
      if (roundUnit === STATS_DURATION_UNITS.month) {
        roundUnit = roundUnit.toUpperCase();
      }

      result[roundMethod](roundUnit);
    }


    if (deltaValue && deltaUnit) {
      if (deltaUnit === STATS_DURATION_UNITS.month) {
        deltaUnit = deltaUnit.toUpperCase();
      }

      result[deltaMethod](deltaValue, deltaUnit);
    }

    return result;
  }

  throw new Error('Date string pattern not recognized');
}

/**
 * Parse date in every formats to moment object
 *
 * @param date
 * @param type
 * @param format
 * @return {*}
 */
export function dateParse(date, type, format) {
  if (typeof date === 'number') {
    return moment.unix(date);
  }

  const momentDate = moment(date, format, true);

  if (!momentDate.isValid()) {
    return parseStringToDateInterval(date, type);
  }

  return momentDate;
}

/**
 * Prepare date to date object
 *
 * @param {number} date
 * @param {string} type
 * @param {string} [unit = 'hour']
 * @param {string} [format = DATETIME_FORMATS.dateTimePicker]
 * @returns {Date}
 */
export function prepareDateToObject(
  date,
  type,
  unit = 'hour',
  format = DATETIME_FORMATS.dateTimePicker,
) {
  const momentDate = dateParse(date, type, format);

  if (momentDate.isValid()) {
    return momentDate.startOf(unit).toDate();
  }

  return null;
}

/**
 * Prepare start of stats interval for month period unit
 *
 * @param {Date|Moment|number} start
 * @returns Moment
 */
export function prepareStatsStartForMonthPeriod(start) {
  return moment.utc(start).startOf('month').tz(moment.tz.guess());
}

/**
 * Prepare stop of stats interval for month period unit
 *
 * @param {Date|Moment|number} stop
 * @returns Moment
 */
export function prepareStatsStopForMonthPeriod(stop) {
  const startOfStopMonthUtc = moment.utc(stop).startOf('month');
  const startOfCurrentMonthUtc = moment.utc().startOf('month');

  if (startOfStopMonthUtc.isSame(startOfCurrentMonthUtc)) {
    return startOfStopMonthUtc.add(1, 'month').tz(moment.tz.guess());
  }

  return startOfStopMonthUtc.tz(moment.tz.guess());
}

/**
 * Find range object
 *
 * @param {string} start
 * @param {string} stop
 * @param {Object} ranges
 * @param {Object} defaultValue
 * @returns {Object}
 */
export function findRange(start, stop, ranges = STATS_QUICK_RANGES, defaultValue = STATS_QUICK_RANGES.custom) {
  return Object.values(ranges)
    .find(range => start === range.start && stop === range.stop) || defaultValue;
}
