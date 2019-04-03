<template lang="pug">
  v-menu(
  ref="menu",
  v-model="opened",
  transition="slide-y-transition",
  :close-on-content-click="false",
  right,
  lazy
  )
    div(slot="activator")
      v-text-field(
      readonly,
      :label="label",
      :error-messages="errorMessages",
      :value="value | date('datePicker', true)",
      :append-icon="clearable ? 'close' : ''",
      @click:append="clear"
      )
    v-date-picker(
    :locale="$i18n.locale",
    :value="value | date('YYYY-MM-DD', true)",
    :title-date-format="titleDateFormat",
    color="primary",
    @input="updateDate"
    )
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS } from '@/constants';

/**
 * Date time picker component
 *
 * @prop {Date} [value=null] - Date value
 * @prop {Boolean} [clearable=false] - if it is true input field will be have cross button with clear event on click
 * @prop {string} [label=''] - Label of the input field
 * @prop {string} [name=null] - Name property in the validation object
 * @prop {Boolean} [roundHours=false] - Deny to change minutes it will be only 0
 *
 * @event value#input
 */
export default {
  $_veeValidate: {
    value() {
      if (!this.value) {
        return this.value;
      }

      return moment(this.value).startOf('day').toDate();
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  props: {
    clearable: {
      type: Boolean,
      default: false,
    },
    value: {
      type: Date,
      default: null,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: null,
    },
  },
  data() {
    return {
      opened: false,
    };
  },
  computed: {
    errorMessages() {
      if (this.$validator && this.errors && this.name) {
        return this.errors.collect(this.name);
      }

      return [];
    },
  },
  methods: {
    updateDate(date) {
      const newValue = new Date(this.value ? this.value.getTime() : null);
      const [year, month, day] = date.split('-');

      newValue.setFullYear(parseInt(year, 10), (parseInt(month, 10) - 1), parseInt(day, 10));
      newValue.setHours(0, 0, 0, 0);

      this.$emit('input', newValue);
    },

    clear() {
      this.$emit('input', null);
    },

    titleDateFormat(date) {
      return moment(date).format(DATETIME_FORMATS.datePicker);
    },
  },
};
</script>
