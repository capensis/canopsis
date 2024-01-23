<template>
  <v-radio-group
    v-field="value"
    :name="name"
    :label="label"
    class="mt-0"
  >
    <v-radio
      v-for="aggregateFunction in availableAggregateFunctions"
      :key="aggregateFunction.value"
      :label="aggregateFunction.label"
      :value="aggregateFunction.value"
      color="primary"
    />
  </v-radio-group>
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
    aggregateFunctions: {
      type: Array,
      default: () => Object.values(AGGREGATE_FUNCTIONS),
    },
  },
  computed: {
    availableAggregateFunctions() {
      return this.aggregateFunctions.map(value => ({
        value,
        label: this.$t(`kpi.aggregateFunctions.${value}`),
      }));
    },
  },
};
</script>
