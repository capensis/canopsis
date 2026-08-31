import { computed, unref } from 'vue';

import { PATTERN_QUICK_RANGES_WITHOUT_CUSTOM } from '@/constants';

import {
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';

/**
 * Returns prepared quick range lists and disables invalid from/to combinations.
 *
 * @param {Object} [options = {}]
 * @param {Array|import('vue').Ref<Array>} [options.fromRanges]
 * @param {Array|import('vue').Ref<Array>} [options.toRanges]
 * @param {Object|import('vue').Ref<Object>} [options.value]
 * @returns {{
 *   fromQuickRanges: import('vue').ComputedRef<Array>,
 *   toQuickRanges: import('vue').ComputedRef<Array>,
 *   itemFromDisabled: Function,
 *   itemToDisabled: Function
 * }}
 */
export const useQuickDateIntervalRange = ({ fromRanges, toRanges, value = {} } = {}) => {
  const { t } = useI18n();

  const fromQuickRanges = computed(() => {
    const ranges = unref(fromRanges) ?? PATTERN_QUICK_RANGES_WITHOUT_CUSTOM;

    return ranges.map(range => ({
      ...range,
      text: range.text ?? t(`quickRanges.types.${range.value}`),
    }));
  });

  const toQuickRanges = computed(() => {
    const ranges = unref(toRanges) ?? PATTERN_QUICK_RANGES_WITHOUT_CUSTOM;

    return ranges.map(range => ({
      ...range,
      text: range.text ?? t(`quickRanges.types.${range.value}`),
    }));
  });

  const allQuickRangesIntervals = computed(() => [...fromQuickRanges.value, ...toQuickRanges.value].reduce((acc, range) => {
      acc[range.value] = {
        start: convertStartDateIntervalToTimestamp(range.start),
        stop: convertStopDateIntervalToTimestamp(range.stop),
      };

      return acc;
    }, {}),
  );

  /**
   * Checks whether a from item should be disabled for the selected to range.
   *
   * @param {Object} item
   * @param {string} item.value
   * @returns {boolean}
   */
  const itemFromDisabled = (item) => {
    const currentValue = unref(value);

    if (!currentValue?.to) {
      return false;
    }

    const { start: itemStart } = allQuickRangesIntervals.value[item.value] ?? {};
    const { start: toStart } = allQuickRangesIntervals.value[currentValue.to] ?? {};

    return toStart <= itemStart;
  };

  /**
   * Checks whether a to item should be disabled for the selected from range.
   *
   * @param {Object} item
   * @param {string} item.value
   * @returns {boolean}
   */
  const itemToDisabled = (item) => {
    const currentValue = unref(value);

    if (!currentValue?.from) {
      return false;
    }

    const { start: fromStart } = allQuickRangesIntervals.value[currentValue.from] ?? {};
    const { start: itemStart } = allQuickRangesIntervals.value[item.value] ?? {};

    return fromStart >= itemStart;
  };

  return {
    fromQuickRanges,
    toQuickRanges,

    itemFromDisabled,
    itemToDisabled,
  };
};
