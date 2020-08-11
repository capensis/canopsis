<template lang="pug">
  v-autocomplete(
    v-field="value",
    :items="timezones",
    :label="fieldLabel",
    :disabled="disabled",
    :name="name"
  )
</template>

<script>
import moment from 'moment-timezone';

import formBaseMixin from '@/mixins/form/base';

export default {
  $_veeValidate: {
    value() {
      return this.value;
    },
    name() {
      return this.name;
    },
  },
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  inject: ['$validator'],
  props: {
    value: {
      type: String,
      required: true,
    },
    label: {
      type: String,
      required: false,
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
      return moment.tz.names();
    },
  },
};
</script>
