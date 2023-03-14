<template lang="pug">
  c-movable-card-iterator-field(v-field="metrics", @add="add")
    template(#item="{ item, index }")
      c-alarm-metric-preset-field(
        v-field="metrics[index]",
        :with-color="withColor",
        :with-aggregate-function="withAggregateFunction",
        :parameters="parameters",
        :disabled-parameters="disabledParameters"
      )
</template>

<script>
import { ALARM_METRIC_PARAMETERS, AGGREGATE_FUNCTIONS } from '@/constants';

import { metricPresetToForm } from '@/helpers/forms/metric';
import { isRatioMetric, isTimeMetric } from '@/helpers/metrics';

import { formArrayMixin } from '@/mixins/form';

export default {
  mixins: [formArrayMixin],
  model: {
    prop: 'metrics',
    event: 'input',
  },
  props: {
    metrics: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'metrics',
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
      default: () => Object.values(ALARM_METRIC_PARAMETERS),
    },
    onlyGroup: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    excludedParameters() {
      const [firstMetric] = this.metrics;

      if (!this.onlyGroup || !firstMetric?.metric) {
        return [];
      }

      if (isRatioMetric(firstMetric.metric)) {
        return this.parameters.filter(metric => !isRatioMetric(metric));
      }

      if (isTimeMetric(firstMetric.metric)) {
        return this.parameters.filter(metric => !isTimeMetric(metric));
      }

      return this.parameters.filter(metric => isTimeMetric(metric) || isRatioMetric(metric));
    },

    disabledParameters() {
      return [
        ...this.metrics.map(({ metric }) => metric),
        ...this.excludedParameters,
      ];
    },
  },
  methods: {
    add() {
      this.addItemIntoArray(metricPresetToForm({
        aggregate_func: this.withAggregateFunction ? AGGREGATE_FUNCTIONS.avg : '',
      }));
    },
  },
};
</script>
