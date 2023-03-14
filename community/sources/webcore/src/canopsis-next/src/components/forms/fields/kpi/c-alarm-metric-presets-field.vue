<template lang="pug">
  c-movable-card-iterator-field(v-field="presets", @add="add")
    template(#item="{ item, index }")
      c-alarm-metric-preset-field(
        v-field="presets[index]",
        :with-color="withColor",
        :with-aggregate-function="withAggregateFunction"
      )
</template>

<script>
import { metricPresetToForm } from '@/helpers/forms/metric';

import { formArrayMixin } from '@/mixins/form';
import { AGGREGATE_FUNCTIONS } from '@/constants';

export default {
  mixins: [formArrayMixin],
  model: {
    prop: 'presets',
    event: 'input',
  },
  props: {
    presets: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'presets',
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
    add() {
      this.addItemIntoArray(metricPresetToForm({
        aggregate_func: this.withAggregateFunction ? AGGREGATE_FUNCTIONS.avg : '',
      }));
    },
  },
};
</script>
