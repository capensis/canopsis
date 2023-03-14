<template lang="pug">
  v-layout.c-alarm-metric-preset-field(column)
    c-alarm-metric-parameters-field.mb-4(
      v-field="preset.metric",
      :label="$t('kpi.selectMetric')",
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
      :label="$t('kpi.calculationMethod')"
    )
</template>

<script>
import { getMetricColor } from '@/helpers/color';

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
  },
  methods: {
    enableColor(value) {
      this.updateField('color', value ? getMetricColor(this.preset.metric) : '');
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
