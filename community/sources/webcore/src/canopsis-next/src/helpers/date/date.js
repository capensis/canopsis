import moment from 'moment-timezone';

import { TIME_UNITS } from '@/constants';

/**
 * Convert timestamps/Date to moment
 *
 * @param {Date|number|moment.Moment} timestamp
 * @returns {moment.Moment}
 */
export const convertTimestampToMoment = (timestamp) => {
  let dateObject;

  // If it's unix timestamp in seconds
  if (typeof timestamp === 'number' && timestamp < 100000000000) {
    dateObject = moment.unix(timestamp);
  } else {
    dateObject = moment(timestamp);
  }

  return dateObject;
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
  const dateObject = convertTimestampToMoment(timestamp);

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
export const convertDateToTimestampByTimezone = (date, timezone = moment.tz.guess()) => convertTimestampToMoment(date)
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
 * Return moment with start of day timestamp
 *
 * @param {Date|number|moment.Moment} date
 */
export const convertDateToStartOfDayMoment = date => moment(convertTimestampToMoment(date).startOf('day').format());

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
export const subtractUnitFromDate = (date, value = 0, unit = TIME_UNITS.second) => convertTimestampToMoment(date)
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
export const formatDate = (date, format) => convertTimestampToMoment(date).format(format);

export default convertTimestampToMoment;
