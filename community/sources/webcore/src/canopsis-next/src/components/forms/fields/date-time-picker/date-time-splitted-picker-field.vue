<template lang="pug">
  div
    v-layout(row, :reverse="reverse")
      v-flex(:class="datePickerFlexClass")
        c-date-picker-field(
          :value="value | date('vuetifyDatePicker', null)",
          :label="!reverse || fullDay ? label : ''",
          :error="errors.has(name)",
          :disabled="disabled",
          :min="min",
          :max="max",
          hide-details,
          @input="updateDate"
        )
      v-flex(v-if="!fullDay")
        time-picker-field(
          :value="value | date('timePicker', null)",
          :label="reverse ? label : ''",
          :error="errors.has(name)",
          :disabled="disabled",
          hide-details,
          @input="updateTime"
        )
    div.v-text-field__details.mt-2
      v-messages(:value="errors.collect(name)", color="error")
</template>

<script>
import { TIME_UNITS } from '@/constants';

import { convertDateToStartOfUnitDateObject } from '@/helpers/date/date';
import { getDateObjectByTime, getDateObjectByDate } from '@/helpers/date/date-time-picker';

import { formBaseMixin } from '@/mixins/form';

import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';

export default {
  $_veeValidate: {
    value() {
      if (!this.value) {
        return this.value;
      }

      return convertDateToStartOfUnitDateObject(this.value, TIME_UNITS.minute);
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: { TimePickerField },
  mixins: [formBaseMixin],
  props: {
    value: {
      type: Date,
      default: null,
    },
    fullDay: {
      type: Boolean,
      default: false,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'date',
    },
    reverse: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: String,
      required: false,
    },
    max: {
      type: String,
      required: false,
    },
  },
  computed: {
    datePickerFlexClass() {
      if (!this.fullDay) {
        return this.reverse ? 'pl-1' : 'pr-1';
      }

      return '';
    },
  },
  methods: {
    updateTime(time) {
      this.updateModel(getDateObjectByTime(this.value, time));
    },

    updateDate(date) {
      this.updateModel(getDateObjectByDate(this.value, date));
    },
  },
};
</script>
