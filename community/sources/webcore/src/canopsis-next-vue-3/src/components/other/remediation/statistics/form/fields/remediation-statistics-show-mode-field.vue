<template>
  <c-select
    :model-value="modelValue"
    :items="showingTypes"
    :label="$t('common.show')"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { REMEDIATION_STATISTICS_CHART_DATA_TYPE } from '@/constants';

export default {
  props: {
    modelValue: {
      type: String,
      required: true,
    },
  },

  emits: ['update:modelValue'],

  computed: {
    showingTypes() {
      return [
        {
          value: REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent,
          text: this.$t('common.percent'),
        },
        {
          value: REMEDIATION_STATISTICS_CHART_DATA_TYPE.remediation,
          text: this.$t('remediation.statistic.remediation'),
        },
      ];
    },
  },
};
</script>
