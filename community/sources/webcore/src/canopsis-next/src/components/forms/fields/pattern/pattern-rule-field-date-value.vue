<template>
  <c-date-time-interval-field
    v-if="isDateRange"
    v-field="value"
    :name="name"
    :disabled="disabled"
    :is-allowed-from-date="isAllowedFromDate"
    :is-allowed-to-date="isAllowedToDate"
    :label="$t('common.value')"
  />
  <c-quick-date-interval-type-field
    v-else-if="isInterval"
    v-field="value.type"
    :name="name"
    :disabled="disabled"
    :ranges="typeQuickRanges"
    :label="$t('common.value')"
  />
  <c-quick-date-interval-type-range-field
    v-else-if="isIntervalRange"
    v-field="value"
    :name="name"
    :disabled="disabled"
    :from-ranges="fromQuickRanges"
    :to-ranges="toQuickRanges"
    :label="$t('common.value')"
  />
</template>

<script>
import { computed, watch } from 'vue';

import { PATTERN_OPERATORS, PATTERN_QUICK_RANGES_WITHOUT_CUSTOM, QUICK_RANGES, TIME_UNITS } from '@/constants';

import { convertDateToTimestamp } from '@/helpers/date/date';
import { getDefaultDateFormRange } from '@/helpers/entities/pattern/form';

import { useI18n } from '@/hooks/i18n';
import { useModelField } from '@/hooks/form/model-field';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    operator: {
      type: String,
      required: true,
    },
    intervalRanges: {
      type: Array,
      default: () => PATTERN_QUICK_RANGES_WITHOUT_CUSTOM,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      required: false,
    },
  },
  setup(props, { emit }) {
    const { t, tc } = useI18n();
    const { updateModel } = useModelField(props, emit);

    const isDateRange = computed(() => (props.operator === PATTERN_OPERATORS.inRangeDates));
    const isIntervalRange = computed(() => (props.operator === PATTERN_OPERATORS.inRangePeriod));
    const isInterval = computed(() => [
      PATTERN_OPERATORS.within,
      PATTERN_OPERATORS.olderThan,
    ].includes(props.operator));

    /**
     * Creates a custom range item object based on the provided duration
     *
     * @param {Object} duration - The duration object containing value and unit
     * @param {number} duration.value - The numeric value of the duration
     * @param {string} duration.unit - The unit of time from TIME_UNITS constant
     * @returns {Object} An object with value and localized text for the custom range
     * @returns {string} returns.value - The custom range value from QUICK_RANGES.custom.value
     * @returns {string} returns.text - Localized text in format "Last {value} {unit(s)}"
     */
    const getCustomRangeItem = (duration = {}) => {
      const { value, unit } = duration;

      const [unitKey] = Object.entries(TIME_UNITS).find(([, unitValue]) => unitValue === unit) ?? [];

      return {
        value: QUICK_RANGES.custom.value,
        text: `${t('common.last')} ${value} ${tc(`common.times.${unitKey}`, value)}`,
        start: `now-${duration.value}${duration.unit}`,
        stop: 'now',
      };
    };

    /**
     * Returns quick ranges array, optionally including a custom range item
     *
     * @param {string} key - The key to check for custom range (default: 'type')
     * @returns {Array} Array of quick range objects. If the current value[key] is custom,
     *                  includes the custom range item generated from value[`${key}Custom`]
     */
    const getQuickRangesWithCustom = (key = 'type') => {
      if (props.value[key] !== QUICK_RANGES.custom.value) {
        return props.intervalRanges;
      }

      return [
        ...props.intervalRanges,
        getCustomRangeItem(props.value[`${key}Custom`]),
      ];
    };

    const typeQuickRanges = computed(() => getQuickRangesWithCustom());
    const fromQuickRanges = computed(() => getQuickRangesWithCustom('from'));
    const toQuickRanges = computed(() => getQuickRangesWithCustom('to'));

    const fromTimestamp = computed(() => (props.value.from ? convertDateToTimestamp(props.value.from) : null));
    const toTimestamp = computed(() => (props.value.to ? convertDateToTimestamp(props.value.to) : null));

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

    /**
     * Clears the current range by resetting the model to default date form range
     */
    const clearRange = () => updateModel(getDefaultDateFormRange());

    watch(() => props.operator, clearRange);

    return {
      isDateRange,
      isIntervalRange,
      isInterval,

      typeQuickRanges,
      fromQuickRanges,
      toQuickRanges,

      isAllowedFromDate,
      isAllowedToDate,

      clearRange,
    };
  },
};
</script>
