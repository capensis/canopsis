<template>
  <c-select-field
    v-field="value"
    :items="availableMetrics"
    :label="label || $t('kpi.metrics.parameter')"
    :name="name"
    :required="required"
    :hide-details="hideDetails"
  />
</template>

<script>
import { KPI_RATING_ENTITY_METRICS, KPI_RATING_SETTINGS_TYPES, KPI_RATING_USER_METRICS } from '@/constants';

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
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'metric',
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    metrics: {
      type: Array,
      required: false,
    },
  },
  computed: {
    availableMetrics() {
      if (this.metrics) {
        return this.metrics;
      }

      const metrics = this.type === KPI_RATING_SETTINGS_TYPES.entity
        ? KPI_RATING_ENTITY_METRICS
        : KPI_RATING_USER_METRICS;

      return metrics.map((value) => {
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
