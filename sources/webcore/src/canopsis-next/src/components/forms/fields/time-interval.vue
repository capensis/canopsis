<template lang="pug">
  v-layout.mb-3(align-center)
    v-flex(xs6)
      v-text-field(
        v-field="interval.interval",
        v-validate="'required|numeric|min_value:1'",
        :error-messages="errors.collect('interval')",
        :min="1",
        name="interval",
        type="number",
        hide-details
      )
    v-flex(xs6)
      v-select(
        v-field="interval.unit",
        v-validate="'required'",
        :items="availableUnits",
        :error-messages="errors.collect('unit')",
        name="unit",
        hide-details
      )
</template>

<script>
import { PERIODIC_REFRESH_UNITS, DEFAULT_TIME_INTERVAL } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'interval',
    event: 'input',
  },
  props: {
    interval: {
      type: Object,
      default: () => ({ ...DEFAULT_TIME_INTERVAL }),
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    availableUnits() {
      return Object.values(PERIODIC_REFRESH_UNITS).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.interval.interval),
      }));
    },
  },
};
</script>
