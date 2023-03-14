<template lang="pug">
  v-radio-group.mt-0(v-field="value", :name="name", :label="label")
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

export default {
  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'aggregate_func',
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
};
</script>
