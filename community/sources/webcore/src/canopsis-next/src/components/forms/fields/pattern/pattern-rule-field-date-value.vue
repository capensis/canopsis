<template>
  <c-date-time-interval-field
    v-if="isDateRange"
    v-field="value"
    :name="name"
    :disabled="disabled"
    :label="$t('common.value')"
  />
  <c-quick-date-interval-type-field
    v-else-if="isInterval"
    v-field="value.type"
    :name="name"
    :disabled="disabled"
    :ranges="intervalRanges"
    :label="$t('common.value')"
  />
  <c-quick-date-interval-type-range-field
    v-else-if="isIntervalRange"
    v-field="value"
    :name="name"
    :disabled="disabled"
    :ranges="intervalRanges"
    :label="$t('common.value')"
  />
</template>

<script>
import { computed } from 'vue';

import { PATTERN_OPERATORS, QUICK_RANGES_WITHOUT_CUSTOM } from '@/constants';

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
    rule: {
      type: Object,
      required: true,
    },
    intervalRanges: {
      type: Array,
      default: () => Object.values(QUICK_RANGES_WITHOUT_CUSTOM),
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
  setup(props) {
    const isDateRange = computed(() => (props.rule.operator === PATTERN_OPERATORS.inRangeDates));
    const isIntervalRange = computed(() => (props.rule.operator === PATTERN_OPERATORS.inRangePeriod));
    const isInterval = computed(() => [
      PATTERN_OPERATORS.within,
      PATTERN_OPERATORS.olderThan,
    ].includes(props.rule.operator));

    return {
      isDateRange,
      isIntervalRange,
      isInterval,
    };
  },
};
</script>
