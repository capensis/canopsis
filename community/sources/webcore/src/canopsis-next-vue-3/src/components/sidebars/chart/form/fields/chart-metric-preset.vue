<template>
  <c-layout class="chart-metric-preset" column>
    <external-metric-field
      v-if="modelValue.external || modelValue.auto"
      :model-value="modelValue.metric"
      :label="modelValue.auto ? $t('kpi.addMetricMask') : $t('kpi.selectMetric')"
      :addable="modelValue.auto"
      required
      @update:model-value="updateMetric"
    />
    <alarm-metric-field
      v-else
      :value="modelValue.metric"
      :parameters="parameters"
      :disabled-parameters="disabledParameters"
      :label="$t('kpi.selectMetric')"
      @update:value="updateMetric"
    />
    <c-name-field
      v-if="modelValue.external && !modelValue.auto && modelValue.metric"
      :model-value="modelValue.label"
      :label="$t('kpi.displayedLabel')"
      @update:model-value="updatePreset({ label: $event })"
    />

    <template v-if="modelValue.metric">
      <c-layout v-if="withColor && !modelValue.auto" align="center" justify="space-between">
        <v-switch
          :model-value="Boolean(modelValue.color)"
          :label="$t('kpi.customColor')"
          color="primary"
          hide-details
          @update:model-value="toggleColor"
        />
        <c-color-picker-field
          v-show="modelValue.color"
          :model-value="modelValue.color"
          @update:model-value="updatePreset({ color: $event })"
        />
      </c-layout>

      <c-radio-group
        v-if="withAggregateFunction || modelValue.external"
        :model-value="modelValue.aggregate_func || AGGREGATE_FUNCTIONS.avg"
        :label="$t('kpi.calculationMethod')"
        hide-details
        @update:model-value="updatePreset({ aggregate_func: $event })"
      >
        <v-radio
          v-for="aggregateFunction in aggregateFunctions"
          :key="aggregateFunction"
          :value="aggregateFunction"
          :label="$t(`kpi.aggregateFunctions.${aggregateFunction}`)"
          color="primary"
        />
      </c-radio-group>
    </template>
  </c-layout>
</template>

<script>
import { AGGREGATE_FUNCTIONS } from '@/constants';
import { getMetricColor } from '@/helpers/entities/metric/color';

import AlarmMetricField from './alarm-metric-field.vue';
import ExternalMetricField from './external-metric-field.vue';

export default {
  components: { AlarmMetricField, ExternalMetricField },
  props: {
    modelValue: {
      type: Object,
      required: true,
    },
    withColor: {
      type: Boolean,
      default: false,
    },
    withAggregateFunction: {
      type: Boolean,
      default: false,
    },
    parameters: {
      type: Array,
      default: undefined,
    },
    disabledParameters: {
      type: Array,
      default: () => [],
    },
  },
  emits: ['update:modelValue'],
  computed: {
    aggregateFunctions() {
      return [
        AGGREGATE_FUNCTIONS.sum,
        AGGREGATE_FUNCTIONS.last,
        AGGREGATE_FUNCTIONS.avg,
        AGGREGATE_FUNCTIONS.min,
        AGGREGATE_FUNCTIONS.max,
      ];
    },
  },
  methods: {
    updatePreset(value) {
      this.$emit('update:modelValue', { ...this.modelValue, ...value });
    },
    updateMetric(metric) {
      this.updatePreset({ metric });
    },
    toggleColor(enabled) {
      this.updatePreset({ color: enabled ? getMetricColor(this.modelValue.metric) : '' });
    },
  },
};
</script>

<style lang="scss" scoped>
.chart-metric-preset {
  flex: 1 1 auto;
  min-width: 0;
}
</style>
