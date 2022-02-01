<template lang="pug">
  v-select.kpi-rating-metric-field(
    v-field="value",
    :items="availableMetrics",
    :label="$t('kpiMetrics.parameter')",
    hide-details
  )
    template(#selection="{ item }")
      span.ellipsis {{ item.text }}
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
          const alarmKey = `alarmList.metrics.${value}`;

          return {
            value,
            text: this.$te(alarmKey) ? this.$t(alarmKey) : this.$t(`users.metrics.${value}`),
          };
        });
    },
  },
};
</script>

<style lang="scss">
$selectIconWidth: 24px;

.kpi-rating-metric-field {
  .v-select__selections {
    width: calc(100% - #{$selectIconWidth});
    flex-wrap: nowrap;
  }
}
</style>
