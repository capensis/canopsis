<template>
  <v-layout
    justify-space-between
    align-center
  >
    <v-label>{{ label }}</v-label>
    <v-flex xs3>
      <v-select
        v-field="value"
        :items="numbers"
        class="mt-0 pt-0"
        hide-details
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { range } from 'lodash';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    min: {
      type: Number,
      default: 1,
    },
  },
  computed: {
    numbers() {
      return [
        /* from | to | step */
        ...range(1, 12, 1),
        ...range(12, 24, 2),
        ...range(24, 52, 4),
      ];
    },
  },
};
</script>
