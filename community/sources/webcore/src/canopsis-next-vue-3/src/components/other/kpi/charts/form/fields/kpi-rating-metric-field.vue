<template>
  <c-select
    :model-value="value"
    :items="availableMetrics"
    :label="label || $t('kpi.metrics.parameter')"
    :name="name"
    :required="required"
    :hide-details="hideDetails"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:value', $event)"
  />
</template>

<script>
import {
  KPI_RATING_ENTITY_METRICS,
  KPI_RATING_SETTINGS_TYPES,
  KPI_RATING_USER_METRICS,
} from '@/constants'

export default {
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
      default: '',
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

  emits: ['update:value'],

  computed: {
    availableMetrics() {
      if (this.metrics) {
        return this.metrics
      }

      const metrics =
        this.type === KPI_RATING_SETTINGS_TYPES.entity
          ? KPI_RATING_ENTITY_METRICS
          : KPI_RATING_USER_METRICS

      return metrics.map((value) => {
        const alarmKey = `alarm.metrics.${value}`

        return {
          value,
          text: this.$te(alarmKey) ? this.$t(alarmKey) : this.$t(`user.metrics.${value}`),
        }
      })
    },
  },
}
</script>
