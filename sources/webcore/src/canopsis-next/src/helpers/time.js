import moment from 'moment';

import { TIME_UNITS } from '@/constants';

const AVAILABLE_UNITS = [
  TIME_UNITS.year,
  TIME_UNITS.month,
  TIME_UNITS.day,
  TIME_UNITS.hour,
  TIME_UNITS.minute,
  TIME_UNITS.second,
];

/**
 * Convert time in unit system to seconds
 *
 * @param {number|string} value
 * @param {string} [unit]
 * @return {number}
 */
export const getSecondsByUnit = (value, unit = TIME_UNITS.second) => moment(0).add(value, unit) / 1000;

/**
 * Convert time to time in unit system
 *
 * @param {number|string} value
 * @param {string} [fromUnit]
 * @param {string} [toUnit]
 * @return {number}
 */
export const getUnitValueFromOtherUnit = (
  value,
  fromUnit = TIME_UNITS.second,
  toUnit = TIME_UNITS.second,
) => moment.duration(value, fromUnit).as(toUnit);

/**
 * Get max available interval value
 *
 * @param {number|string} value
 * @param {string} fromUnit
 * @param {string[]} availableUnits
 * @return {Interval}
 */
export const getMaxAvailableIntervalFromUnit = (
  value = 0,
  fromUnit = TIME_UNITS.second,
  availableUnits = AVAILABLE_UNITS,
) => {
  let unitValue;

  const maxUnit = availableUnits.find((unit) => {
    unitValue = Math.floor(getUnitValueFromOtherUnit(value, fromUnit, unit));
    return unitValue;
  });

  return {
    interval: unitValue,
    unit: maxUnit || TIME_UNITS.second,
  };
};
