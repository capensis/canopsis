<template>
  <c-select
    :model-value="currentValue"
    :items="availableMetrics"
    :label="label || t('kpi.metrics.parameter')"
    :name="name"
    :required="required"
    :hide-details="hideDetails"
    item-title="text"
    item-value="value"
    @update:model-value="updateValue"
  />
</template>

<script>
import { computed } from 'vue';

import {
  KPI_RATING_ENTITY_METRICS,
  KPI_RATING_SETTINGS_TYPES,
  KPI_RATING_USER_METRICS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: String,
      required: false,
    },
    value: {
      type: String,
      required: false,
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

  emits: ['update:modelValue', 'update:value'],

  setup(props, { emit }) {
    const { t, te } = useI18n();

    const currentValue = computed(() => props.modelValue ?? props.value ?? '');

    const availableMetrics = computed(() => {
      if (props.metrics) {
        return props.metrics;
      }

      const metrics =
        props.type === KPI_RATING_SETTINGS_TYPES.entity
          ? KPI_RATING_ENTITY_METRICS
          : KPI_RATING_USER_METRICS;

      return metrics.map((metricValue) => {
        const alarmKey = `alarm.metrics.${metricValue}`;

        return {
          value: metricValue,
          text: te(alarmKey) ? t(alarmKey) : t(`user.metrics.${metricValue}`),
        };
      });
    });

    const updateValue = (nextValue) => {
      const normalizedValue = nextValue ?? '';

      if (normalizedValue === currentValue.value) {
        return;
      }

      emit('update:modelValue', normalizedValue);
      emit('update:value', normalizedValue);
    };

    return {
      t,
      currentValue,
      availableMetrics,
      updateValue,
    };
  },
};
</script>
