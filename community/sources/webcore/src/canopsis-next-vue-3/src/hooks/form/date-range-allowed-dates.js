import { computed, unref } from 'vue';

import { convertDateToTimestamp } from '@/helpers/date/date';

/**
 * Returns date-range guards for from/to date pickers.
 *
 * @param {Object|import('vue').Ref<Object>} value
 * @returns {{ isAllowedFromDate: Function, isAllowedToDate: Function }}
 */
export const useDateRangeAllowedDates = (value) => {
  const fromTimestamp = computed(() => {
    const unwrappedValue = unref(value);

    return unwrappedValue?.from ? convertDateToTimestamp(unwrappedValue.from) : null;
  });

  const toTimestamp = computed(() => {
    const unwrappedValue = unref(value);

    return unwrappedValue?.to ? convertDateToTimestamp(unwrappedValue.to) : null;
  });

  const isAllowedFromDate = date => !toTimestamp.value || convertDateToTimestamp(date) < toTimestamp.value;
  const isAllowedToDate = date => !fromTimestamp.value || convertDateToTimestamp(date) > fromTimestamp.value;

  return {
    isAllowedFromDate,
    isAllowedToDate,
  };
};
