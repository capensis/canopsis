<template lang="pug">
  v-radio-group.mt-0(
    v-field="value",
    :label="label",
    hide-details,
    mandatory,
    row
  )
    v-radio(
      v-for="method in methods",
      :key="method",
      :label="$t(`stateSetting.calculationMethods.${method}`)",
      :value="method",
      color="primary"
    )
</template>

<script>
import { STATE_SETTING_CALCULATION_METHODS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: STATE_SETTING_CALCULATION_METHODS.share,
    },
    label: {
      type: String,
      default: '',
    },
  },
  computed: {
    methods() {
      return STATE_SETTING_CALCULATION_METHODS;
    },
  },
};
</script>
