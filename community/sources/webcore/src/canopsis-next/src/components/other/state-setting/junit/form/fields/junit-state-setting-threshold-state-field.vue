<template>
  <v-layout>
    <v-flex
      class="mt-2"
      xs4
    >
      {{ label }}
    </v-flex>
    <v-flex xs4>
      <c-number-field
        class="ma-1 pa-0"
        v-if="isNumberType"
        v-field="value"
        :name="name"
      />
      <c-percents-field
        class="ma-1 pa-0"
        v-else
        v-field="value"
        :name="name"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { JUNIT_STATE_SETTING_THRESHOLDS_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Number, String],
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    type: {
      type: Number,
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    name: {
      type: String,
      default: 'state',
    },
  },
  computed: {
    isNumberType() {
      return this.type === JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number;
    },
  },
};
</script>
