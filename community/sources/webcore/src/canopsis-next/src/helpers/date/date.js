import moment from 'moment-timezone';

import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';

/**
 * Convert timestamps/Date to moment
 *
 * @param {Date | number | string | moment.Moment} date
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
 * Convert timestamp from source timezone to local timezone with time keeping
 *
 * @param {number} timestamp
 * @param {string} sourceTimezone
 * @param {string} [localTimezone = moment.tz.guess()]
 * @returns {Object}
 */
export const convertTimestampToMomentByTimezone = (
  timestamp,
  sourceTimezone = moment.tz.guess(),
  localTimezone = moment.tz.guess(),
) => {
  const dateObject = convertDateToMoment(timestamp);

  if (sourceTimezone === localTimezone) {
    return dateObject;
  }

  return dateObject.tz(sourceTimezone).tz(localTimezone, true);
};

/**
 * Convert date to timestamp with keep time
 *
 * @param {Date|number|moment.Moment} date
 * @param {string} timezone
 * @returns {number}
 */
export const convertDateToTimestampByTimezone = (date, timezone = moment.tz.guess()) => convertDateToMoment(date)
  .tz(timezone, true)
  .unix();

/**
 * Check if date is start of day
 *
 * @param {Date|moment.Moment} date
 * @param {string} unit
 * @returns {boolean}
 */
export const isStartOfDay = (date, unit = 'seconds') => {
  const dateMoment = moment(date);

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
  const dateMoment = moment(date);

  return dateMoment.clone().endOf('day').diff(dateMoment, unit) === 0;
};

/**
 * Convert date to string format
 *
 * @param {Date|number|moment.Moment} date
 * @param {string} [format = DATETIME_FORMATS.long]
 * @param {string} [defaultValue = '']
 * @return {string}
 */
export const convertDateToString = (date, format = DATETIME_FORMATS.long, defaultValue = '') => {
  if (!date) {
    return defaultValue;
  }

  const dateObject = convertDateToMoment(date);

  if (!dateObject?.isValid()) {
    console.warn('Could not build a valid `moment` object from input.');

    return date;
  }

  return dateObject.format(DATETIME_FORMATS[format] ?? format);
};

/**
 * Convert date to string. If the date is today, only the time is returned.
 *
 * @param {Date|number|moment.Moment} date
 * @param {string} [format]
 * @param {string} [defaultValue]
 */
export const convertDateToStringWithFormatForToday = (date, format, defaultValue) => {
  const dateObject = convertDateToMoment(date);
  const resultFormat = dateObject.isSame(new Date(), 'day') ? DATETIME_FORMATS.time : format;

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
  convertTimestampToMomentByTimezone(date, timezone),
  format,
  defaultValue,
);

export const convertDateToStartOfUnitMoment = (date, unit) => convertDateToMoment(date).startOf(unit);

export const convertDateToStartOfUnitString = (date, unit, format) => convertDateToString(
  convertDateToStartOfUnitMoment(date, unit),
  format,
);

export const convertDateToEndOfUnitMoment = (date, unit) => convertDateToMoment(date).endOf(unit);

export const convertDateToEndOfUnitString = (date, unit, format) => convertDateToString(
  convertDateToEndOfUnitMoment(date, unit),
  format,
);

/**
 * Return moment with start of day timestamp
 *
 * @param {Date|number|moment.Moment} date
 */
export const convertDateToStartOfDayMoment = (date) => {
  const startOfMoment = convertDateToStartOfUnitMoment(date, TIME_UNITS.day);
  /* Format date to string without time and timezone */
  const formattedStartOfMoment = startOfMoment.format(DATETIME_FORMATS.long);

  /* Format to moment object */
  return convertDateToMoment(formattedStartOfMoment, DATETIME_FORMATS.long);
};

export const convertDateToStartOfDayTimestamp = date => convertDateToStartOfDayMoment(date).unix();

export const convertDateToStartOfDayDateObject = date => convertDateToStartOfDayMoment(date).toDate();

/**
 * Return moment with end of day timestamp
 *
 * @param {Date|number|moment.Moment} date
 */
export const convertDateToEndOfDayMoment = (date) => {
  const endOfMoment = convertDateToEndOfUnitMoment(date, TIME_UNITS.day);
  /* Format date to string without time and timezone */
  const formattedEndOfMoment = endOfMoment.format(DATETIME_FORMATS.long);

  /* Format to moment object */
  return convertDateToMoment(formattedEndOfMoment, DATETIME_FORMATS.long);
};

export const convertDateToEndOfDayTimestamp = date => convertDateToEndOfDayMoment(date).unix();

export const convertDateToEndOfDayDateObject = date => convertDateToEndOfDayMoment(date).toDate();

/**
 * Subtract value from date by unit
 *
 * @param {Date|number|moment.Moment} date
 * @param {number} [value = 0]
 * @param {string} [unit = TIME_UNITS.second]
 * @return {number}
 */
export const subtractUnitFromDate = (date, value = 0, unit = TIME_UNITS.second) => convertDateToMoment(date)
  .clone()
  .subtract(value, unit)
  .unix();

/**
 * Format date/timestamp/unix/moment to string format
 *
 * @param {Date|number|moment.Moment} date
 * @param {string} format
 * @return {string}
 */
export const formatDate = (date, format) => convertDateToMoment(date).format(format);

/**
 * Return all available timezones
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

export const getDiffBetweenDate = (left, right, unit = TIME_UNITS.second) => convertDateToMoment(left)
  .diff(convertDateToMoment(right), unit);

/**
 * Getting a now timestamp
 *
 * @param {Date|number|moment.Moment} [date]
 * @return {number}
 */
export const getDateTimestamp = date => convertDateToMoment(date).unix();

/**
 * Getting a now timestamp
 *
 * @return {number}
 */
export const getNowTimestamp = () => getDateTimestamp();

export const getWeekdayNumber = date => convertDateToMoment(date).isoWeekday();

export const convertDateToDateObject = (date, format) => convertDateToMoment(date, format).toDate();

export const getLocalTimezone = () => moment.tz.guess();

export default convertDateToMoment;
