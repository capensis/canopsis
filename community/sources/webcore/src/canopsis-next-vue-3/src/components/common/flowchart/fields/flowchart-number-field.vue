<template>
  <c-layout justify-space-between align-center>
    <v-label>{{ label }}</v-label>
    <v-col cols="3">
      <v-select
        :model-value="value"
        :items="numbers"
        class="mt-0 pt-0"
        hide-details
        @update:model-value="updateModel"
      />
    </v-col>
  </c-layout>
</template>

<script>
import { range } from 'lodash';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  emits: ['input', 'update:value'],
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
