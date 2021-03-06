import moment from 'moment';

import { AVAILABLE_SORTED_TIME_UNITS, TIME_UNITS } from '@/constants';

/**
 * @typedef { "y" | "M" | "w" | "d" | "h" | "m" | "s" } DurationUnit
 */

/**
 * @typedef {Object} Duration
 * @property {number} seconds
 * @property {DurationUnit} unit
 */

/**
 * @typedef {Object} DurationForm
 * @property {number} value
 * @property {DurationUnit} unit
 */

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

  return moment.duration(value, fromUnit).as(toUnit);
};

/**
 * Convert duration from unit to "seconds"
 *
 * @param {number} value
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {number}
 */
export const toSeconds = (value, unit = TIME_UNITS.second) =>
  convertUnit(value, unit, TIME_UNITS.second);

/**
 * Convert duration from "seconds" to unit
 *
 * @param {number} value
 * @param {DurationUnit} [unit = TIME_UNITS.second]
 * @returns {number}
 */
export const fromSeconds = (value, unit = TIME_UNITS.second) =>
  convertUnit(value, TIME_UNITS.second, unit);

/**
 * Convert Duration object to DurationForm
 *
 * @param {Duration} duration
 * @param {number} [duration.seconds = 0]
 * @param {DurationUnit} [duration.unit = TIME_UNITS.second]
 * @returns {DurationForm}
 */
export const durationToForm = ({ seconds = 0, unit = TIME_UNITS.second } = {}) =>
  ({ unit, value: fromSeconds(seconds, unit) });

/**
 * Convert DurationForm object to Duration
 *
 * @param {DurationForm} duration
 * @param {number} [duration.value = 0]
 * @param {DurationUnit} [duration.unit = TIME_UNITS.second]
 * @returns {Duration}
 */
export const formToDuration = ({ value = 0, unit = TIME_UNITS.second } = {}) =>
  ({ unit, seconds: toSeconds(value, unit) });


/**
 * Get max available interval value
 *
 * @param {DurationForm} [durationForm = { value: 0, unit: TIME_UNITS.second }]
 * @param {DurationUnit[]} [availableUnits = AVAILABLE_SORTED_TIME_UNITS]
 * @return {DurationForm}
 */
export const formToMaxByAvailableUnitsForm = (
  durationForm = { value: 0, unit: TIME_UNITS.second },
  availableUnits = AVAILABLE_SORTED_TIME_UNITS,
) => {
  const { value, unit } = durationForm;
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
