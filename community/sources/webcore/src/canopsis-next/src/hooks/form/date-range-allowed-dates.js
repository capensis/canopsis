import { computed, unref } from 'vue';

import { convertDateToTimestamp } from '@/helpers/date/date';

/**
 * Hook providing date range validation for "from" and "to" date fields.
 * Returns validator functions that ensure "from" is before "to" and "to" is after "from".
 *
 * @param {Object|Ref} value - The date range value { from, to } (reactive ref, computed, or plain).
 * @returns {Object} Object with isAllowedFromDate and isAllowedToDate validator functions.
 */
export const useDateRangeAllowedDates = (value) => {
  const fromTimestamp = computed(
    () => {
      const val = unref(value);

      return val?.from ? convertDateToTimestamp(val.from) : null;
    },
  );

  const toTimestamp = computed(
    () => {
      const val = unref(value);

      return val?.to ? convertDateToTimestamp(val.to) : null;
    },
  );

  /**
   * Validates if a given date is allowed as a "from" date.
   * A date is allowed if there's no upper bound (toTimestamp) or if the date is before the upper bound.
   *
   * @param {Date|string|number} date - The date to validate
   * @returns {boolean} True if the date is allowed as a "from" date, false otherwise
   */
  const isAllowedFromDate = date => !toTimestamp.value || convertDateToTimestamp(date) < toTimestamp.value;

  /**
   * Validates if a given date is allowed as a "to" date.
   * A date is allowed if there's no lower bound (fromTimestamp) or if the date is after the lower bound.
   *
   * @param {Date|string|number} date - The date to validate
   * @returns {boolean} True if the date is allowed as a "to" date, false otherwise
   */
  const isAllowedToDate = date => !fromTimestamp.value || convertDateToTimestamp(date) > fromTimestamp.value;

  return {
    isAllowedFromDate,
    isAllowedToDate,
  };
};
