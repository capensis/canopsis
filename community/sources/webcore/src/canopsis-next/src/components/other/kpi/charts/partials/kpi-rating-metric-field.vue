<template lang="pug">
  c-select-field(
    v-field="value",
    :items="availableMetrics",
    :label="$t('kpi.metrics.parameter')",
    name="metric",
    hide-details
  )
</template>

<script>
import { getAvailableMetricsByCriteria } from '@/helpers/metrics';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    criteria: {
      type: Object,
      required: false,
    },
  },
  computed: {
    availableMetrics() {
      return getAvailableMetricsByCriteria(this.criteria?.label)
        .map((value) => {
          const alarmKey = `alarm.metrics.${value}`;

          return {
            value,
            text: this.$te(alarmKey) ? this.$t(alarmKey) : this.$t(`user.metrics.${value}`),
          };
        });
    },
  },
};
</script>
