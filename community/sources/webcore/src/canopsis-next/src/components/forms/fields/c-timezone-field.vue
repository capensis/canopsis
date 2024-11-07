<template>
  <v-autocomplete
    v-field="value"
    v-bind="$attrs"
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
  inheritAttrs: false,
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
    server: {
      type: Boolean,
      default: false,
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
