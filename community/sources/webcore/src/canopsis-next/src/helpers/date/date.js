import moment from 'moment-timezone';

import { DATETIME_FORMATS, TIME_UNITS } from '@/constants';

/**
 * @typedef {Date | number | string | moment.Moment} LocalDate
 */

/**
 * Convert timestamps/Date to moment
 *
 * @param {Date|number|moment.Moment} timestamp
 * @returns {moment.Moment}
 */
export const convertDateToMoment = (timestamp) => {
  /**
   * NOTE: If it's unix timestamp in seconds
   */
  if (typeof timestamp === 'number' && timestamp < 100000000000) {
    return moment.unix(timestamp);
  }

  return moment(timestamp);
};

/**
 * Convert duration to interval object
 *
 * @param duration
 * @return {{unit: string, interval: number}}
 */
export const convertDurationToIntervalObject = (duration) => {
  const durationUnits = [
    TIME_UNITS.year,
    TIME_UNITS.month,
    TIME_UNITS.week,
    TIME_UNITS.week,
    TIME_UNITS.day,
    TIME_UNITS.hour,
    TIME_UNITS.minute,
    TIME_UNITS.second,
  ];

  const durationType = durationUnits.find(unit => moment.duration(duration, 'seconds').as(unit) % 1 === 0);

  return {
    interval: moment.duration(duration, 'seconds').as(durationType),
    unit: durationType,
  };
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
export const convertDateToTimestampByTimezone = (date, timezone = moment.tz.guess()) =>
  convertDateToMoment(date).tz(timezone, true).unix();

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
 * @param {string} format
 * @param {boolean} [ignoreTodayChecker]
 * @param {string} [defaultValue]
 * @return {string}
 */
export const convertDateToString = (date, format, ignoreTodayChecker, defaultValue) => {
  let momentFormat = DATETIME_FORMATS[format] || format;

  if (!date) {
    return defaultValue || date;
  }

  const dateObject = convertDateToMoment(date);

  if (!dateObject || !dateObject.isValid()) {
    console.warn('Could not build a valid `moment` object from input.');
    return date;
  }

  if (!ignoreTodayChecker && dateObject.isSame(new Date(), 'day')) {
    momentFormat = DATETIME_FORMATS.time;
  }

  return dateObject.format(momentFormat);
};

/**
 * Convert date to start of unit as moment
 *
 * @param {LocalDate} date
 * @param {string} unit
 * @return {moment.Moment}
 */
export const convertDateToStartOfUnitMoment = (date, unit) => convertDateToMoment(date).startOf(unit);

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
 * @param {Date|number|moment.Moment} date
 */
export const convertDateToStartOfDayMoment = (date) => {
  const startOfMoment = convertDateToMoment(date).startOf('day');
  /* Format date to string without time and timezone */
  const formattedStartOfMoment = startOfMoment.format(DATETIME_FORMATS.long);

  /* Format to moment object */
  return moment(formattedStartOfMoment, DATETIME_FORMATS.long);
};

/**
 * Return moment with end of day timestamp
 *
 * @param {Date|number|moment.Moment} date
 */
export const convertDateToEndOfDayMoment = (date) => {
  const endOfMoment = convertDateToMoment(date).endOf('day');
  /* Format date to string without time and timezone */
  const formattedEndOfMoment = endOfMoment.format(DATETIME_FORMATS.long);

  /* Format to moment object */
  return moment(formattedEndOfMoment, DATETIME_FORMATS.long);
};

/**
 * Getting a now timestamp
 *
 * @return {number}
 */
export const getNowTimestamp = () => moment().unix();

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

export default convertDateToMoment;
