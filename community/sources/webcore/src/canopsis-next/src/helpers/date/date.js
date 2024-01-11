import { isNumber } from 'lodash';
import moment from 'moment-timezone';

import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';

/**
 * @typedef {Date | number | string | moment.Moment} LocalDate
 */

/**
 * Convert timestamps/Date to moment
 *
 * @param {LocalDate} date
 * @param {string} [format]
 * @returns {moment.Moment}
 */
export const convertDateToMoment = (date, format) => {
  /**
   * NOTE: If it's unix timestamp in seconds
   */
  if (typeof date === 'number' && date < 100000000000) {
    return moment.unix(date);
  }

  if (format && typeof date === 'string') {
    return moment(date, format);
  }

  return moment(date);
};

/**
 * Return all available timezones
 *
 * @return {Array}
 */
export const getTimezones = () => moment.tz.names();

/**
 * Set locale for each dates
 *
 * @param {string} locale
 */
export const setDateLocale = (locale) => {
  moment.locale(locale);
};

/**
 * Get current dates locale
 */
export const getDateLocale = () => moment.locale();

/**
 * Update messages for locale
 *
 * @param {string} language
 * @param {Object} messages
 */
export const updateDateLocaleMessages = (language, messages) => {
  const currentLocale = getDateLocale();

  moment.updateLocale(language, messages);

  setDateLocale(currentLocale);
};

/**
 * Return diff between two dates
 *
 * @param {LocalDate} left
 * @param {LocalDate} right
 * @param {string} unit
 * @return {number}
 */
export const getDiffBetweenDates = (left, right, unit = TIME_UNITS.second) => convertDateToMoment(left)
  .diff(convertDateToMoment(right), unit);

/**
 * Getting a now timestamp
 *
 * @param {LocalDate} [date]
 * @return {number}
 */
export const convertDateToTimestamp = date => convertDateToMoment(date).unix();

/**
 * Getting a now timestamp
 *
 * @return {number}
 */
export const getNowTimestamp = () => convertDateToTimestamp();

/**
 * Return day of week number
 *
 * @param {LocalDate} date
 * @return {number}
 */
export const getWeekdayNumber = date => convertDateToMoment(date).isoWeekday();

/**
 * Return local timezone
 *
 * @return {string}
 */
export const getLocaleTimezone = () => moment.tz.guess();

/**
 * Subtract value from date by unit
 *
 * @param {LocalDate} date
 * @param {number} [value = 0]
 * @param {string} [unit = TIME_UNITS.second]
 * @return {number}
 */
export const subtractUnitFromDate = (date, value = 0, unit = TIME_UNITS.second) => convertDateToMoment(date)
  .clone()
  .subtract(value, unit)
  .unix();

/**
 * Subtract value from date by unit
 *
 * @param {LocalDate} date
 * @param {number} [value = 0]
 * @param {string} [unit = TIME_UNITS.second]
 * @return {number}
 */
export const addUnitToDate = (date, value = 0, unit = TIME_UNITS.second) => convertDateToMoment(date)
  .clone()
  .add(value, unit)
  .unix();

/**
 * Convert date to native date object
 *
 * @param {LocalDate} date
 * @param {string} [format]
 * @return {Date}
 */
export const convertDateToDateObject = (date, format) => convertDateToMoment(date, format).toDate();

/**
 * Convert date from source timezone to local timezone with time keeping
 *
 * @param {LocalDate} timestamp
 * @param {string} sourceTimezone
 * @param {string} [targetTimezone = getLocaleTimezone()]
 * @returns {Object}
 */
export const convertDateToMomentByTimezone = (
  timestamp,
  sourceTimezone = getLocaleTimezone(),
  targetTimezone = getLocaleTimezone(),
) => {
  const dateObject = convertDateToMoment(timestamp);

  if (sourceTimezone === targetTimezone) {
    return dateObject;
  }

  return dateObject.tz(sourceTimezone).tz(targetTimezone, true);
};

/**
 * Convert date from source timezone to local timezone date object
 *
 * @param {LocalDate} timestamp
 * @param {string} sourceTimezone
 * @param {string} [targetTimezone]
 * @return {Date}
 */
export const convertDateToDateObjectByTimezone = (
  timestamp,
  sourceTimezone,
  targetTimezone,
) => convertDateToMomentByTimezone(timestamp, sourceTimezone, targetTimezone).toDate();

/**
 * Convert date to timestamp with keep time
 *
 * @param {Date|number|moment.Moment} date
 * @param {string} [timezone = getLocaleTimezone()]
 * @returns {number}
 */
export const convertDateToTimestampByTimezone = (date, timezone = getLocaleTimezone()) => convertDateToMoment(date)
  .tz(timezone, true)
  .unix();

/**
 * Check if date is start of day
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @returns {boolean}
 */
export const isStartOfDay = (date, unit = 'seconds') => {
  const dateMoment = convertDateToMoment(date);

  return dateMoment.clone().startOf('day').diff(dateMoment, unit) === 0;
};

/**
 * Check if date is end of day
 *
 * @param {Date|moment.Moment} date
 * @param {string} unit
 * @returns {boolean}
 */
export const isEndOfDay = (date, unit = 'seconds') => {
  const dateMoment = convertDateToMoment(date);

  return dateMoment.clone().endOf('day').diff(dateMoment, unit) === 0;
};

/**
 * Convert date to string format
 *
 * @param {LocalDate} date
 * @param {!string} [format = DATETIME_FORMATS.long]
 * @param {string | null} [defaultValue = '']
 * @return {string}
 */
export const convertDateToString = (date, format = DATETIME_FORMATS.long, defaultValue = '') => {
  if (!date) {
    return defaultValue;
  }

  const dateObject = convertDateToMoment(date);

  if (!dateObject?.isValid()) {
    console.warn('Could not build a valid `moment` object from input.');

    return defaultValue ?? date;
  }

  return dateObject.format(DATETIME_FORMATS[format] ?? format);
};

/**
 * Convert date to string. If the date is today, only the time is returned.
 *
 * @param {LocalDate} date
 * @param {string} [format]
 * @param {string} [defaultValue]
 */
export const convertDateToStringWithFormatForToday = (date, format, defaultValue) => {
  const dateObject = convertDateToMoment(date);
  const resultFormat = dateObject.isSame(Date.now(), 'day') ? DATETIME_FORMATS.time : format;

  return convertDateToString(date, resultFormat, defaultValue);
};

/**
 * Convert date to timezone date string
 *
 * @param date
 * @param timezone
 * @param [format]
 * @param [defaultValue]
 * @return {string}
 */
export const convertDateToTimezoneDateString = (date, timezone, format, defaultValue) => convertDateToString(
  convertDateToMomentByTimezone(date, timezone),
  format,
  defaultValue,
);

/**
 * Convert date to start of unit as moment
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @return {moment.Moment}
 */
export const convertDateToStartOfUnitMoment = (date, unit) => convertDateToMoment(date).startOf(unit);

/**
 * Convert date to start of unit as timestamp
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @return {number}
 */
export const convertDateToStartOfUnitTimestamp = (date, unit) => convertDateToTimestamp(
  convertDateToStartOfUnitMoment(date, unit),
);

/**
 * Convert date to start of unit as native date
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @return {Date}
 */
export const convertDateToStartOfUnitDateObject = (date, unit) => convertDateToStartOfUnitMoment(date, unit).toDate();

/**
 * Convert date to start of unit as moment
 *
 * @param date
 * @param unit
 * @return {moment.Moment}
 */
export const convertDateToEndOfUnitMoment = (date, unit) => convertDateToMoment(date).endOf(unit);

/**
 * Convert date to end of unit as timestamp
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @return {number}
 */
export const convertDateToEndOfUnitTimestamp = (date, unit) => convertDateToTimestamp(
  convertDateToEndOfUnitMoment(date, unit),
);

/**
 * Convert date to end of unit as timestamp
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @return {Date}
 */
export const convertDateToEndOfUnitDateObject = (date, unit) => convertDateToDateObject(
  convertDateToEndOfUnitMoment(date, unit),
);

/**
 * Convert date to start of unit as formatted string
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @param {?string} [format]
 * @return {string}
 */
export const convertDateToStartOfUnitString = (date, unit, format) => convertDateToString(
  convertDateToStartOfUnitMoment(date, unit),
  format,
);

/**
 * Return moment with start of day timestamp
 *
 * @param {LocalDate} date
 */
export const convertDateToStartOfDayMoment = (date) => {
  const startOfMoment = convertDateToStartOfUnitMoment(date, TIME_UNITS.day);
  /* Format date to string without time and timezone */
  const formattedStartOfMoment = startOfMoment.format(DATETIME_FORMATS.long);

  /* Format to moment object */
  return convertDateToMoment(formattedStartOfMoment, DATETIME_FORMATS.long);
};

/**
 * Convert date to start of day as timestamp
 *
 * @param {LocalDate} date
 * @return {number}
 */
export const convertDateToStartOfDayTimestamp = date => convertDateToStartOfDayMoment(date).unix();

/**
 * Convert date to start of day as timestamp
 *
 * @param {LocalDate} date
 * @param {string} [timezone = getLocaleTimezone()]
 * @return {number}
 */
export const convertDateToStartOfDayTimestampByTimezone = (date, timezone = getLocaleTimezone()) => (
  convertDateToMomentByTimezone(convertDateToStartOfDayMoment(date), timezone, getLocaleTimezone()).unix()
);

/**
 * Convert date to start of day as native date
 *
 * @param {LocalDate} date
 * @return {Date}
 */
export const convertDateToStartOfDayDateObject = date => convertDateToStartOfDayMoment(date).toDate();

/**
 * Return moment with end of day timestamp
 *
 * @param {LocalDate} date
 */
export const convertDateToEndOfDayMoment = (date) => {
  const endOfMoment = convertDateToEndOfUnitMoment(date, TIME_UNITS.day);
  /* Format date to string without time and timezone */
  const formattedEndOfMoment = endOfMoment.format(DATETIME_FORMATS.long);

  /* Format to moment object */
  return convertDateToMoment(formattedEndOfMoment, DATETIME_FORMATS.long);
};

/**
 * Convert date to end of day as native timestamp
 *
 * @param {LocalDate} date
 * @return {number}
 */
export const convertDateToEndOfDayTimestamp = date => convertDateToEndOfDayMoment(date).unix();

/**
 * Convert date to end of day as native date
 *
 * @param {LocalDate} date
 * @return {Date}
 */
export const convertDateToEndOfDayDateObject = date => convertDateToEndOfDayMoment(date).toDate();

/**
 * Convert date to special format with new timezone without strict parsing
 *
 * @param {LocalDate} date
 * @param {string} [format = DATETIME_FORMATS.long]
 * @param {string} [timezone = getLocaleTimezone()]
 * @returns {string}
 */
export const convertDateToStringWithNewTimezone = (
  date,
  format = DATETIME_FORMATS.long,
  timezone = getLocaleTimezone(),
) => (
  date
    ? convertDateToMoment(date).tz(timezone).format(format)
    : ''
);

/**
 * Check time unit is valid
 *
 * @param {string} unit
 * @return {boolean}
 */
export const isValidTimeUnit = unit => Object.values(TIME_UNITS).includes(unit);

/**
 * Check is interval valid
 *
 * @param {Interval | *} value
 * @return {boolean}
 */
export const isValidDateInterval = value => value
  && isNumber(value?.from)
  && isNumber(value?.to)
  && value.from < value.to;

/**
 * Get days in month
 *
 * @param {LocalDate} date
 * @return {number}
 */
export const getDaysInMonth = date => convertDateToMoment(date)
  .clone()
  .daysInMonth();

/**
 * Get start week day
 *
 * @param {LocalDate} date
 * @param {number} week
 * @return {number}
 */
export const getWeekStartDay = (date, week) => convertDateToMoment(date)
  .clone()
  .week(week)
  .startOf(TIME_UNITS.week)
  .toDate();

/**
 * Get end week day
 *
 * @param {LocalDate} date
 * @param {number} week
 * @return {number}
 */
export const getWeekEndDay = (date, week) => convertDateToMoment(date)
  .clone()
  .week(week)
  .endOf(TIME_UNITS.week)
  .toDate();

/**
 * Check date is same by unit
 *
 * @param {LocalDate} leftDate
 * @param {LocalDate} rightDate
 * @param {string} unit
 * @return {number}
 */
export const isSameDates = (leftDate, rightDate, unit) => convertDateToMoment(leftDate)
  .isSame(rightDate, unit);

/**
 * Check date is same by unit
 *
 * @param {LocalDate} date
 * @return {number}
 */
export const getWeekNumber = date => convertDateToMoment(date).week();

/**
 * Get date by week
 *
 * @param {LocalDate} date
 * @param {number} week
 * @return {number}
 */
export const getDateByWeekNumber = (date, week) => convertDateToMoment(date)
  .week(week)
  .toDate();

/**
 * Get date by month
 *
 * @param {LocalDate} date
 * @param {number} month
 * @return {number}
 */
export const getDateByMonthNumber = (date, month) => convertDateToMoment(date)
  .month(month)
  .toDate();

/**
 * Check is date before target date
 *
 * @param {LocalDate} date
 * @param {LocalDate} targetDate
 * @return {number}
 */
export const isDateBefore = (date, targetDate) => convertDateToMoment(date)
  .isBefore(targetDate);

export default convertDateToMoment;
