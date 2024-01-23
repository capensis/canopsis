<template>
  <v-radio-group
    v-field="value"
    :label="label"
    class="mt-0"
    color="primary"
    hide-details
    mandatory
  >
    <v-radio
      :label="$t('common.variableTypes.number')"
      :value="numberValue"
      color="primary"
    />
    <v-radio
      :label="$tc('common.percent')"
      :value="percentValue"
      color="primary"
    />
  </v-radio-group>
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
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    label: {
      type: String,
      default: '',
    },
    numberValue: {
      type: [Number, String],
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.number,
    },
    percentValue: {
      type: [Number, String],
      default: JUNIT_STATE_SETTING_THRESHOLDS_TYPES.percent,
    },
    name: {
      type: String,
      default: 'type',
    },
  },
};
</script>
