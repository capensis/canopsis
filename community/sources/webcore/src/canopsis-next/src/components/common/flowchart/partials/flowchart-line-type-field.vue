<template lang="pug">
  v-layout(row, justify-space-between, align-center)
    v-label {{ label }}
    v-flex(xs3)
      v-select.mt-0.pt-0(v-field="value", :items="types", hide-details)
</template>

<script>
import { LINE_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
  },
  computed: {
    types() {
      return Object.values(LINE_TYPES).map(value => ({
        value,
        text: this.$t(`flowchart.lineTypes.${value}`),
      }));
    },
  },
};
</script>
