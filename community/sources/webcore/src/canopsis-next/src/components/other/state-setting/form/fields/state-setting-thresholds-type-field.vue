<template>
  <v-radio-group
    class="mt-0"
    v-field="value"
    :label="label"
    hide-details
    mandatory
    color="primary"
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
      default: STATE_SETTING_THRESHOLD_TYPES.number,
    },
    label: {
      type: String,
      default: '',
    },
    numberValue: {
      type: [Number, String],
      default: STATE_SETTING_THRESHOLD_TYPES.number,
    },
    percentValue: {
      type: [Number, String],
      default: STATE_SETTING_THRESHOLD_TYPES.percent,
    },
    name: {
      type: String,
      default: 'type',
    },
  },
};
</script>
