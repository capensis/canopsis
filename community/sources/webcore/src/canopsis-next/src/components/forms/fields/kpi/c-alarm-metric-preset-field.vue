<template lang="pug">
  v-layout.c-alarm-metric-preset-field(column)
    c-alarm-metric-parameters-field(
      v-field="preset.metric",
      :label="$t('kpi.selectMetric')"
    )
    v-layout(v-if="withColor", align-center, justify-space-between)
      v-flex(xs11)
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
    v-radio-group.mt-0(
      v-if="withAggregateFunction",
      v-field="preset.aggregate_func",
      :label="$t('kpi.calculationMethod')"
    )
      v-radio(
        v-for="aggregateFunction in availableAggregateFunctions",
        :key="aggregateFunction.value",
        :label="aggregateFunction.label",
        :value="aggregateFunction.value",
        color="primary"
      )
</template>

<script>
import { AGGREGATE_FUNCTIONS } from '@/constants';

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
  computed: {
    availableAggregateFunctions() {
      return Object.values(AGGREGATE_FUNCTIONS).map(value => ({
        value,
        label: this.$t(`kpi.aggregateFunctions.${value}`),
      }));
    },
  },
  methods: {
    enableColor(value) {
      this.updateField('color', value ? '#fff' : '');
    },
  },
};
</script>

<style lang="scss">
.c-alarm-metric-preset-field {
  &__color-picker {
    width: max-content;
  }
}
</style>
