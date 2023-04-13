<template lang="pug">
  v-layout.c-alarm-metric-preset-field(column)
    c-alarm-metric-parameters-field(
      :value="preset.metric",
      :label="preset.auto ? $t('kpi.addMetricMask') : $t('kpi.selectMetric')",
      :parameters="parameters",
      :disabled-parameters="disabledParameters",
      :addable="preset.auto",
      required,
      with-external,
      @input="updateMetric"
    )
    c-name-field(
      v-if="!preset.auto && preset.metric && isExternalMetric",
      v-field="preset.label",
      :label="$t('kpi.displayedLabel')",
      required
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
    isExternalMetric() {
      return !this.isInternalMetric(this.preset.metric);
    },

    aggregateFunctions() {
      return getAggregateFunctionsByMetric(this.preset.metric);
    },
  },
  methods: {
    isInternalMetric() {
      return this.parameters.includes(this.preset.metric);
    },

    enableColor(value) {
      this.updateField('color', value ? getMetricColor(this.preset.metric) : '');
    },

    updateMetric(metric) {
      if (this.withAggregateFunction) {
        this.updateModel({
          ...this.preset,
          metric,
          aggregate_func: getDefaultAggregateFunctionByMetric(metric),
          label: this.isInternalMetric(metric) ? '' : this.preset.label,
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
