<template lang="pug">
  v-layout.c-alarm-metric-preset-field(column)
    c-alarm-metric-parameters-field(
      :value="preset.metric",
      :label="$t('kpi.selectMetric')",
      :parameters="parameters",
      :disabled-parameters="disabledParameters",
      required,
      @input="updateMetric"
    )
    v-layout(v-if="withColor", align-center, justify-space-between)
      v-switch(
        :label="$t('kpi.customColor')",
        :input-value="!!preset.color",
        color="primary",
        @change="enableColor($event)"
      )
      c-color-picker-field.c-alarm-metric-preset-field__color-picker(
        v-show="preset.color",
        v-field="preset.color"
      )
    c-alarm-metric-aggregate-function-field(
      v-if="withAggregateFunction",
      v-field="preset.aggregate_func",
      :aggregate-functions="aggregateFunctions",
      :label="$t('kpi.calculationMethod')"
    )
</template>

<script>
import { getMetricColor } from '@/helpers/color';
import { getAggregateFunctionsByMetric, getDefaultAggregateFunctionByMetric } from '@/helpers/metrics';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [
    formMixin,
  ],
  model: {
    prop: 'preset',
    event: 'input',
  },
  props: {
    preset: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: 'preset',
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
      required: false,
    },
    disabledParameters: {
      type: Array,
      required: false,
    },
  },
  computed: {
    aggregateFunctions() {
      return getAggregateFunctionsByMetric(this.preset.metric);
    },
  },
  methods: {
    enableColor(value) {
      this.updateField('color', value ? getMetricColor(this.preset.metric) : '');
    },

    updateMetric(metric) {
      if (this.withAggregateFunction) {
        this.updateModel({
          ...this.preset,
          metric,
          aggregate_func: getDefaultAggregateFunctionByMetric(metric),
        });
      } else {
        this.updateField('metric', metric);
      }
    },
  },
};
</script>

<style lang="scss">
.c-alarm-metric-preset-field {
  &__color-picker {
    width: max-content;
    flex: unset;
  }
}
</style>
