import { AVAILABLE_SORTED_TIME_UNITS, TIME_UNITS } from '@/constants';

import { getMaxAvailableValueFromUnit } from '@/helpers/duration';

/**
 * Filter for getting max available interval value from unit
 *
 * @param {number|string} [value = 0]
 * @param {string} [fromUnit = TIME_UNITS.second]
 * @param {string[]} [availableUnits = AVAILABLE_SORTED_TIME_UNITS]
 * @return {string}
 */
export default function (
  value = 0,
  fromUnit = TIME_UNITS.second,
  availableUnits = AVAILABLE_SORTED_TIME_UNITS,
) {
  const { unit, interval } = getMaxAvailableValueFromUnit(value, fromUnit, availableUnits);

  return `${interval}${unit}`;
}
