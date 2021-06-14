import { TIME_UNITS, AVAILABLE_SORTED_TIME_UNITS } from '@/constants';

import { formToMaxByAvailableUnitsForm } from '@/helpers/date/duration';

/**
 * Filter for getting max available interval value from unit
 *
 * @param {number|string} [value = 0]
 * @param {string} [unit = TIME_UNITS.second]
 * @param {string[]} [availableUnits = AVAILABLE_SORTED_TIME_UNITS]
 * @return {string}
 */
export default function (
  value = 0,
  unit = TIME_UNITS.second,
  availableUnits = AVAILABLE_SORTED_TIME_UNITS,
) {
  const durationForm = formToMaxByAvailableUnitsForm({ value, unit }, availableUnits);

  return `${durationForm.value}${durationForm.unit}`;
}
