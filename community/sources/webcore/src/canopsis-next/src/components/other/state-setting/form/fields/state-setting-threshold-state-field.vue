<template lang="pug">
  v-layout(row)
    v-flex.mt-2(xs4) {{ label }}
    v-flex(xs4)
      c-number-field.ma-1.pa-0(v-if="isNumberType", v-field="value", :name="name")
      c-percents-field.ma-1.pa-0(v-else, v-field="value", :name="name")
</template>

<script>
import { STATE_SETTING_THRESHOLD_TYPES } from '@/constants';

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
      default: STATE_SETTING_THRESHOLD_TYPES.number,
    },
    name: {
      type: String,
      default: 'state',
    },
  },
  computed: {
    isNumberType() {
      return this.type === STATE_SETTING_THRESHOLD_TYPES.number;
    },
  },
};
</script>
