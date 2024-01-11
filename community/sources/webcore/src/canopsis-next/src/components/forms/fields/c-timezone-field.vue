<template>
  <v-autocomplete
    v-field="value"
    :items="timezones"
    :label="fieldLabel"
    :disabled="disabled"
    :name="name"
  />
</template>

<script>
import { getTimezones } from '@/helpers/date/date';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'timezone',
    },
  },
  computed: {
    fieldLabel() {
      return this.label || this.$t('common.timezone');
    },

    timezones() {
      return getTimezones();
    },
  },
};
</script>
