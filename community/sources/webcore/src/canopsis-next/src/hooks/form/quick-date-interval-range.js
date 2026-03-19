import { computed, unref } from 'vue';

import { PATTERN_QUICK_RANGES_WITHOUT_CUSTOM } from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { useI18n } from '@/hooks/i18n';

/**
 * Hook providing shared logic for quick date interval range selection (from/to).
 * Handles preparation of ranges with labels and validation rules for from/to constraints.
 *
 * @param {Object} options - Options for the hook.
 * @param {Array} [options.fromRanges] - Ranges for "from" field.
 * @param {Array} [options.toRanges] - Ranges for "to" field.
 * @param {Object} [options.value] - Current value { from, to } (reactive ref or plain).
 * @returns {Object} Prepared ranges and disabled checkers.
 */
export const useQuickDateIntervalRange = ({
  fromRanges,
  toRanges,
  value = {},
} = {}) => {
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

  const allQuickRangesIntervals = computed(() => [
    ...fromQuickRanges.value,
    ...toQuickRanges.value,
  ].reduce((acc, range) => {
    acc[range.value] = {
      start: convertStartDateIntervalToTimestamp(range.start),
      stop: convertStopDateIntervalToTimestamp(range.stop),
    };

    return acc;
  }, {}));

  /**
   * Determines if a "from" range item should be disabled.
   * Disables items whose start is after or equal to the current "to" range start.
   *
   * @param {Object} item - The range item to check.
   * @param {string} item.value - The range value key.
   * @returns {boolean} True if the item should be disabled.
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
   * Determines if a "to" range item should be disabled.
   * Disables items whose start is before or equal to the current "from" range start.
   *
   * @param {Object} item - The range item to check.
   * @param {string} item.value - The range value key.
   * @returns {boolean} True if the item should be disabled.
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
