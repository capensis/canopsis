import { isNil, isNumber } from 'lodash';
import moment from 'moment';

import 'moment-duration-format';

import {
  AVAILABLE_SORTED_TIME_UNITS,
  DATETIME_FORMATS,
  DAYS_IN_MONTH,
  DAYS_IN_WEEK,
  DAYS_IN_YEAR,
  DEFAULT_DURATION_FORMAT,
  MONTHS_IN_YEAR,
  TIME_UNITS,
} from '@/constants';

import { isValidTimeUnit } from '@/helpers/date/date';

/**
 * @typedef { "y" | "M" | "w" | "d" | "h" | "m" | "s" } DurationUnit
 */

/**
 * @typedef {Object} Duration
 * @property {number} value
 * @property {DurationUnit} unit
 */

/**
 * @typedef {Duration} DurationWithEnabled
 * @property {boolean} enabled
 */

/**
 * Check unit is valid
 *
 * @param {DurationUnit} unit
 * @returns {boolean}
 */
export const isValidUnit = unit => Object.values(TIME_UNITS).includes(unit);

/**
 * Convert duration from one unit to another unit
 *
 * @param {number} value
 * @param {DurationUnit} [fromUnit = TIME_UNITS.second]
 * @param {DurationUnit} [toUnit = TIME_UNITS.second]
 * @returns {number}
 */
export const convertUnit = (value, fromUnit = TIME_UNITS.second, toUnit = TIME_UNITS.second) => {
  if (fromUnit === toUnit) {
    return value;
  }

  /**
   * Using of the moment, we are faced with the problems of converting any units into months and years.
   * TODO: Should be removed after change format from seconds to value on backend size.
   * @link https://github.com/moment/moment/issues/5892
   */
  if (fromUnit === TIME_UNITS.year) {
    if (toUnit === TIME_UNITS.month) {
      return value * MONTHS_IN_YEAR;
    }

    if (toUnit === TIME_UNITS.week) {
      return value * (DAYS_IN_YEAR / DAYS_IN_WEEK);
    }

    return moment.duration(value * DAYS_IN_YEAR, TIME_UNITS.day).as(toUnit);
  }

  if (fromUnit === TIME_UNITS.month) {
    if (toUnit === TIME_UNITS.year) {
      return value / MONTHS_IN_YEAR;
    }

    return moment.duration(value * DAYS_IN_MONTH, TIME_UNITS.day).as(toUnit);
  }

  const momentDuration = moment.duration(value, fromUnit);

  if (toUnit === TIME_UNITS.month) {
    return momentDuration.as(TIME_UNITS.day) / DAYS_IN_MONTH;
  }

  if (toUnit === TIME_UNITS.year) {
    return momentDuration.as(TIME_UNITS.day) / DAYS_IN_YEAR;
  }

  return momentDuration.as(toUnit);
};

/**
 * Convert duration from unit to "seconds"
 *
 * @param {number} value
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {number}
 */
export const toSeconds = (value, unit = TIME_UNITS.second) => convertUnit(value, unit, TIME_UNITS.second);

/**
 * Convert duration from "seconds" to unit
 *
 * @param {number} value
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {number}
 */
export const fromSeconds = (value, unit = TIME_UNITS.second) => convertUnit(value, TIME_UNITS.second, unit);

/**
 * Convert duration to seconds
 *
 * @param {Duration} duration
 * @return {number}
 */
export const durationToSeconds = ({ value, unit } = {}) => toSeconds(value, unit);

/**
 * Convert Duration object to form
 *
 * @param {Duration} [duration = { value: 1, unit: TIME_UNITS.second }]
 * @returns {Duration}
 */
export const durationToForm = duration => (duration ? { ...duration } : { value: 1, unit: TIME_UNITS.second });

/**
 * Convert DurationWithEnabled object to form
 *
 * @param {DurationWithEnabled} duration
 * @return {DurationWithEnabled}
 */
export const durationWithEnabledToForm = ({ value, unit, enabled = false } = {}) => ({
  ...durationToForm({ value, unit }),

  enabled,
});

/**
 * Get max available interval value
 *
 * @param {Duration} [duration = { value: 0, unit: TIME_UNITS.second }]
 * @param {DurationUnit[]} [availableUnits = AVAILABLE_SORTED_TIME_UNITS]
 * @return {Duration}
 */
export const convertDurationToMaxUnitDuration = (
  duration = { value: 0, unit: TIME_UNITS.second },
  availableUnits = AVAILABLE_SORTED_TIME_UNITS,
) => {
  const { value, unit } = duration;
  let unitValue = value;

  const maxUnit = availableUnits.find((availableUnit) => {
    unitValue = Math.floor(convertUnit(value, unit, availableUnit));
    return unitValue;
  });

  return {
    value: unitValue,
    unit: maxUnit || unit,
  };
};

/**
 * Filter for getting max available interval value from unit
 *
 * @param {number|string} [value = 0]
 * @param {string} [unit = TIME_UNITS.second]
 * @param {string[]} [availableUnits = AVAILABLE_SORTED_TIME_UNITS]
 * @return {string}
 */
export const convertDurationToMaxUnitDurationString = (
  value = 0,
  unit = TIME_UNITS.second,
  availableUnits = AVAILABLE_SORTED_TIME_UNITS,
) => {
  const durationForm = convertDurationToMaxUnitDuration({ value, unit }, availableUnits);

  return `${durationForm.value}${durationForm.unit}`;
};

/**
 * Convert duration to more readable format
 *
 * @param {number | Duration} duration
 * @param {string} [format = DEFAULT_DURATION_FORMAT]
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {string}
 */
export const convertDurationToString = (duration, format = DEFAULT_DURATION_FORMAT, unit = TIME_UNITS.second) => {
  if (isNil(duration)) {
    return '';
  }

  const resultFormat = DATETIME_FORMATS[format] || format;
  const preparedDuration = isNumber(duration)
    ? duration
    : duration?.value;

  return moment
    .duration(preparedDuration, duration?.unit ?? unit)
    .format(resultFormat, { trim: 'both final' }) || '0s';
};

/**
 * Check duration is valid
 *
 * @param {Duration | *} duration
 * @return {boolean}
 */
export const isValidDuration = duration => isNumber(duration?.value) && isValidTimeUnit(duration.unit);
