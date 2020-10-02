import { getMaxAvailableIntervalFromUnit } from '@/helpers/time';

/**
 * Filter for getting max available interval value from unit
 *
 * @param {number|string} value
 * @param {string} fromUnit
 * @param {string[]} availableUnits
 * @return {string}
 */
export default function (value = 0, fromUnit, availableUnits) {
  const { unit, interval } = getMaxAvailableIntervalFromUnit(value, fromUnit, availableUnits);

  return `${interval}${unit}`;
}
