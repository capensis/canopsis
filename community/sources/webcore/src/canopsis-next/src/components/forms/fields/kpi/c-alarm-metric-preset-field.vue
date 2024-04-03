<template>
  <v-layout
    class="c-alarm-metric-preset-field"
    column
  >
    <c-alarm-external-metric-parameters-field
      v-if="preset.external || preset.auto"
      :value="preset.metric"
      :label="preset.auto ? $t('kpi.addMetricMask') : $t('kpi.selectMetric')"
      :addable="preset.auto"
      :name="`${name}.metric`"
      required
      @input="updateMetric"
    />
    <c-alarm-metric-parameters-field
      v-else
      :value="preset.metric"
      :label="$t('kpi.selectMetric')"
      :parameters="parameters"
      :disabled-parameters="disabledParameters"
      :name="`${name}.metric`"
      required
      @input="updateMetric"
    />
    <c-name-field
      v-if="!preset.auto && preset.metric && preset.external"
      v-field="preset.label"
      :label="$t('kpi.displayedLabel')"
      :name="`${name}.label`"
    />
    <template v-if="preset.metric">
      <v-layout
        v-if="withColor && !preset.auto"
        align-center
        justify-space-between
      >
        <v-switch
          :label="$t('kpi.customColor')"
          :input-value="!!preset.color"
          color="primary"
          @change="enableColor($event)"
        />
        <c-color-picker-field
          v-show="preset.color"
          v-field="preset.color"
          class="c-alarm-metric-preset-field__color-picker"
        />
      </v-layout>
      <c-alarm-metric-aggregate-function-field
        v-if="withAggregateFunction || isExternalMetric"
        v-field="preset.aggregate_func"
        :aggregate-functions="aggregateFunctions"
        :label="$t('kpi.calculationMethod')"
      />
    </template>
  </v-layout>
</template>

<script>
import { getMetricColor } from '@/helpers/entities/metric/color';
import { getAggregateFunctionsByMetric, getDefaultAggregateFunctionByMetric } from '@/helpers/entities/metric/list';

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
      return this.preset.external || this.preset.auto;
    },

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
