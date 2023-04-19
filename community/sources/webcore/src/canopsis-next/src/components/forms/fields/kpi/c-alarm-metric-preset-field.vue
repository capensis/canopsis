<template lang="pug">
  v-layout.c-alarm-metric-preset-field(column)
    c-alarm-metric-parameters-field(
      :value="preset.metric",
      :label="preset.auto ? $t('kpi.addMetricMask') : $t('kpi.selectMetric')",
      :parameters="preset.auto || onlyExternal ? [] : parameters",
      :disabled-parameters="disabledParameters",
      :addable="preset.auto",
      :name="`${name}.metric`",
      :with-external="withExternal",
      required,
      @input="updateMetric"
    )
    c-name-field(
      v-if="!preset.auto && preset.metric && isExternalMetric",
      v-field="preset.label",
      :label="$t('kpi.displayedLabel')",
      :name="`${name}.label`",
      required
    )
    template(v-if="preset.metric")
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
        v-if="withAggregateFunction || isExternalMetric",
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
    withExternal: {
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
    onlyExternal: {
      type: Boolean,
      default: false,
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
    isInternalMetric(metric) {
      return this.parameters.includes(metric);
    },

    enableColor(value) {
      this.updateField('color', value ? getMetricColor(this.preset.metric) : '');
    },

    getNewAggregatedFunction(metric) {
      if (this.isInternalMetric(metric)) {
        return this.withAggregateFunction ? getDefaultAggregateFunctionByMetric(metric) : undefined;
      }

      return getDefaultAggregateFunctionByMetric(metric);
    },

    updateMetric(metric) {
      this.updateModel({
        ...this.preset,
        metric,
        aggregate_func: this.getNewAggregatedFunction(metric),
        label: this.isInternalMetric(metric) ? '' : this.preset.label,
      });
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
