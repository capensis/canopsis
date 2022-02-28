import moment from 'moment-timezone';

import { DATETIME_FORMATS, QUICK_RANGES, DATETIME_INTERVAL_TYPES, TIME_UNITS } from '@/constants';

import { convertDateToStartOfUnitString, getLocaleTimezone, subtractUnitFromDate } from '@/helpers/date/date';

/**
 * Convert a date interval string to moment date object
 *
 * @param {string} string
 * @param {string} type
 * @returns {Moment}
 */
export const convertStringToDateInterval = (string, type) => {
  const matches = string.match(/^now(([+--])(\d+)([hdwmMy]{1}))?(\/([hdwmMy]{1}))?$/);

  if (matches) {
    const result = moment().utc();
    const operator = matches[2];
    const deltaValue = matches[3];
    const deltaUnit = matches[4];
    const roundUnit = matches[6];

    if (roundUnit) {
      if (type === DATETIME_INTERVAL_TYPES.start) {
        result.startOf(roundUnit);
      } else {
        result.endOf(roundUnit);
      }
    }

    if (deltaValue && deltaUnit) {
      if (operator === '+') {
        result.add(deltaValue, deltaUnit);
      } else {
        result.subtract(deltaValue, deltaUnit);
      }
    }

    return result;
  }

  throw new Error('Date string pattern not recognized');
};

/**
 * Parse date in every formats to moment object
 *
 * @param {number | string | moment.Moment }date
 * @param {string} type
 * @param {string} format
 * @return {number | moment.Moment}
 */
export const convertDateIntervalToMoment = (date, type, format) => {
  if (typeof date === 'number') {
    return moment.unix(date);
  }

  const momentDate = moment(date, format, true);

  if (!momentDate.isValid()) {
    return convertStringToDateInterval(date, type);
  }

  return momentDate;
};

/**
 * Convert date interval to timestamp unix system.
 *
 * @param {number} date
 * @param {string} type
 * @param {string} format
 * @return {number}
 */
export const convertDateIntervalToTimestamp = (date, type, format) => convertDateIntervalToMoment(date, type, format)
  .unix();

/**
 * Convert from value to timestamp or moment
 *
 * @param {LocalDate} date
 * @param {string} [format = DATETIME_FORMATS.datePicker]
 * @return {moment.Moment}
 */
export const convertStartDateIntervalToMoment = (
  date,
  format = DATETIME_FORMATS.datePicker,
) => convertDateIntervalToMoment(
  date,
  DATETIME_INTERVAL_TYPES.start,
  format,
);

/**
 * Convert from value to timestamp or moment
 *
 * @param {LocalDate} date
 * @param {string} [format = DATETIME_FORMATS.datePicker]
 * @return {number}
 */
export const convertStartDateIntervalToTimestamp = (
  date,
  format = DATETIME_FORMATS.datePicker,
) => convertDateIntervalToTimestamp(
  date,
  DATETIME_INTERVAL_TYPES.start,
  format,
);

/**
 * Convert to value to timestamp or moment
 *
 * @param {LocalDate} date
 * @param {string} [format = DATETIME_FORMATS.datePicker]
 * @return {moment.Moment}
 */
export const convertStopDateIntervalToMoment = (
  date,
  format = DATETIME_FORMATS.datePicker,
) => convertDateIntervalToMoment(
  date,
  DATETIME_INTERVAL_TYPES.stop,
  format,
);

/**
 * Convert to value to timestamp or moment
 *
 * @param {LocalDate} date
 * @param {string} [format = DATETIME_FORMATS.datePicker]
 * @return {number}
 */
export const convertStopDateIntervalToTimestamp = (
  date,
  format = DATETIME_FORMATS.datePicker,
) => convertDateIntervalToTimestamp(
  date,
  DATETIME_INTERVAL_TYPES.stop,
  format,
);

/**
 * Prepare date to date object
 *
 * @param {number} date
 * @param {string} type
 * @param {string} [unit = 'hour']
 * @param {string} [format = DATETIME_FORMATS.dateTimePicker]
 * @returns {Date}
 */
export const convertDateIntervalToDateObject = (
  date,
  type,
  unit = 'hour',
  format = DATETIME_FORMATS.dateTimePicker,
) => {
  const momentDate = convertDateIntervalToMoment(date, type, format);

  if (momentDate.isValid()) {
    return momentDate.startOf(unit).toDate();
  }

  return null;
};

/**
 * Prepare start of stats interval for month period unit
 *
 * @param {Date|Moment|number} start
 * @returns Moment
 */
export const prepareStatsStartForMonthPeriod = start => moment.utc(start).startOf('month').tz(getLocaleTimezone());

/**
 * Prepare stop of stats interval for month period unit
 *
 * @param {Date|Moment|number} stop
 * @returns Moment
 */
export const prepareStatsStopForMonthPeriod = (stop) => {
  const startOfStopMonthUtc = moment.utc(stop).startOf('month');
  const startOfCurrentMonthUtc = moment.utc().startOf('month');

  if (startOfStopMonthUtc.isSame(startOfCurrentMonthUtc)) {
    return startOfStopMonthUtc.add(1, 'month').tz(getLocaleTimezone());
  }

  return startOfStopMonthUtc.tz(getLocaleTimezone());
};

/**
 * Find range object
 *
 * @param {string} start
 * @param {string} stop
 * @param {Object} ranges
 * @param {Object} defaultValue
 * @returns {Object}
 */
export const findQuickRangeValue = (
  start,
  stop,
  ranges = QUICK_RANGES,
  defaultValue = QUICK_RANGES.custom,
) => Object.values(ranges)
  .find(range => start === range.start && stop === range.stop) ?? defaultValue;

/**
 * Get value from quick range period
 *
 * @param {string} value
 * @param {string | number} start
 * @param {string | number} stop
 * @return {Object}
 */
export const getValueFromQuickRange = ({ value, start, stop }) => {
  if (value === QUICK_RANGES.custom.value) {
    return {
      periodUnit: TIME_UNITS.hour,
      periodValue: 1,

      tstart: convertDateToStartOfUnitString(
        subtractUnitFromDate(Date.now(), 1, TIME_UNITS.hour),
        TIME_UNITS.hour,
        DATETIME_FORMATS.dateTimePicker,
      ),

      tstop: convertDateToStartOfUnitString(
        Date.now(),
        TIME_UNITS.hour,
        DATETIME_FORMATS.dateTimePicker,
      ),
    };
  }

  return {
    tstart: start,
    tstop: stop,
  };
};
