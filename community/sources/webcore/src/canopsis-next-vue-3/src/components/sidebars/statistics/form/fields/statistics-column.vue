<template>
  <c-layout column>
    <kpi-rating-metric-field
      :model-value="modelValue.metric"
      :metrics="availableMetrics"
      :type="type"
      :label="t('common.column', 1)"
      :name="`${name}.column`"
      required
      @update:model-value="updateColumnField('metric', $event)"
    />
    <c-enabled-field
      v-model="customLabel"
      :label="t('settings.columns.customLabel')"
      class="pa-0 my-2"
      hide-details
      @update:model-value="updateCustomLabel"
    />
    <c-text-field
      v-if="customLabel"
      :model-value="modelValue.label"
      :label="t('common.label')"
      :name="`${name}.label`"
      required
      @update:model-value="updateColumnField('label', $event)"
    />
    <c-enabled-field
      :model-value="modelValue.split"
      :label="t('settings.statisticsWidgetColumn.split')"
      :disabled="!hasPossibilityToSplit"
      hide-details
      @update:model-value="updateColumnField('split', $event)"
    />
    <c-select
      v-if="modelValue.split"
      :model-value="modelValue.criteria"
      :items="ratingSettings"
      :label="t('common.infos')"
      :name="`${name}.infos`"
      item-title="label"
      item-value="id"
      required
      @update:model-value="updateColumnField('criteria', $event)"
    />
  </c-layout>
</template>

<script>
import { computed, ref, watch } from 'vue';

import {
  KPI_RATING_SETTINGS_TYPES,
  STATISTICS_WIDGETS_ENTITY_METRICS,
  STATISTICS_WIDGETS_USER_METRICS,
  STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE,
} from '@/constants';

import { useModelField } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import KpiRatingMetricField from '@/components/other/kpi/charts/form/fields/kpi-rating-metric-field.vue';

export default {
  components: {
    KpiRatingMetricField,
  },

  props: {
    modelValue: {
      type: Object,
      default: () => ({}),
    },
    type: {
      type: Number,
      default: KPI_RATING_SETTINGS_TYPES.entity,
    },
    name: {
      type: String,
      default: 'column',
    },
    ratingSettings: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t, te } = useI18n();
    const { updateField } = useModelField(props, emit);
    const customLabel = ref(!!props.modelValue.label);

    const isEntityType = computed(() => props.type === KPI_RATING_SETTINGS_TYPES.entity);

    const hasColumnSelected = computed(() => !!props.modelValue?.metric);

    const hasPossibilityToSplit = computed(
      () => hasColumnSelected.value &&
        STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE.includes(props.modelValue?.metric),
    );

    const availableMetrics = computed(() => {
      const metrics = isEntityType.value
        ? STATISTICS_WIDGETS_ENTITY_METRICS
        : STATISTICS_WIDGETS_USER_METRICS;

      return metrics.map((value) => {
        const kpiKey = `kpi.statisticsWidgets.metrics.${value}`;
        const alarmKey = `alarm.metrics.${value}`;

        let text;

        if (te(kpiKey)) {
          text = t(kpiKey);
        } else if (te(alarmKey)) {
          text = t(alarmKey);
        }

        return {
          value,
          text: text ?? t(`user.metrics.${value}`),
        };
      });
    });

    const updateColumnField = (field, value) => {
      if (field === 'split') {
        if (props.modelValue.split === value) {
          return;
        }

        updateField('split', value);

        return;
      }

      const nextValue = value ?? '';

      if (props.modelValue[field] === nextValue) {
        return;
      }

      if (field === 'metric' && !props.modelValue.metric && !nextValue) {
        return;
      }

      updateField(field, nextValue);
    };

    const updateCustomLabel = (checked) => {
      if (checked) {
        return;
      }

      updateColumnField('label', '');
    };

    watch(
      hasPossibilityToSplit,
      (value) => {
        if (value || (!props.modelValue.split && !props.modelValue.criteria)) {
          return;
        }

        emit('update:modelValue', {
          ...props.modelValue,
          split: false,
          criteria: '',
        });
      },
      { flush: 'post' },
    );

    return {
      t,
      customLabel,
      hasPossibilityToSplit,
      availableMetrics,
      updateColumnField,
      updateCustomLabel,
    };
  },
};
</script>
