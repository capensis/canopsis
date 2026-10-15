<template>
  <v-select
    v-field="value"
    :label="$t('common.disableDuringPeriods')"
    :items="duringItems"
    :name="name"
    multiple
  />
</template>

<script>
import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: false,
    },
    name: {
      type: String,
      default: 'disable_during_periods',
    },
  },
  computed: {
    duringItems() {
      return [PBEHAVIOR_TYPE_TYPES.maintenance, PBEHAVIOR_TYPE_TYPES.pause, PBEHAVIOR_TYPE_TYPES.inactive];
    },
  },
};
</script>
