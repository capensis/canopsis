<template>
  <v-radio-group
    v-field="value"
    :label="label"
    :disabled="disabled"
    hide-details
    mandatory
    row
  >
    <v-radio
      v-for="method in methods"
      :key="method"
      :label="$t(`stateSetting.thresholdMethods.${method}`)"
      :value="method"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { STATE_SETTING_THRESHOLDS_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: STATE_SETTING_THRESHOLDS_METHODS.share,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    methods() {
      return STATE_SETTING_THRESHOLDS_METHODS;
    },
  },
};
</script>
