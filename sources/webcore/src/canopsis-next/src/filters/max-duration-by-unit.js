import { AVAILABLE_SORTED_TIME_UNITS, TIME_UNITS } from '@/constants';

import { getMaxAvailableIntervalFromUnit } from '@/helpers/time';

/**
 * Filter for getting max available interval value from unit
 *
 * @param {number|string} [value]
 * @param {string} [fromUnit]
 * @param {string[]} [availableUnits]
 * @return {string}
 */
export default function (
  value = 0,
  fromUnit = TIME_UNITS.second,
  availableUnits = AVAILABLE_SORTED_TIME_UNITS,
) {
  const { unit, interval } = getMaxAvailableIntervalFromUnit(value, fromUnit, availableUnits);

  return `${interval}${unit}`;
}
