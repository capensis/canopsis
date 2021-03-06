<template lang="pug">
  div
    v-layout(row)
      v-flex.pr-1(v-if="reverse && !fullDay", xs4)
        time-picker-field(
          :value="value | date('timePicker', true, null)",
          :label="label",
          :error="errors.has(name)",
          :disabled="disabled",
          hide-details,
          @input="updateTime"
        )
      v-flex(:class="datePickerFlexClass")
        date-picker-field(
          :value="value | date('YYYY-MM-DD', true, null)",
          :label="!reverse || fullDay ? label : ''",
          :error="errors.has(name)",
          :disabled="disabled",
          hide-details,
          @input="updateDate"
        )
      v-flex.pr-1(v-if="!reverse && !fullDay", xs4)
        time-picker-field(
          :value="value | date('timePicker', true, null)",
          :error="errors.has(name)",
          :disabled="disabled",
          hide-details,
          @input="updateTime"
        )
    div.v-text-field__details.mt-2
      div.v-messages.theme--light.error--text
        div.v-messages__wrapper
          div.v-messages__message(
            v-for="error in errorMessages",
            :key="error"
          ) {{ error }}
</template>

<script>
import { convertTimestampToMoment } from '@/helpers/date/date';
import { updateTime, updateDate } from '@/helpers/date/date-time-picker';

import formBaseMixin from '@/mixins/form/base';

import DatePickerField from '@/components/forms/fields/date-picker/date-picker-field.vue';
import TimePickerField from '@/components/forms/fields/time-picker/time-picker-field.vue';

export default {
  $_veeValidate: {
    value() {
      if (!this.value) {
        return this.value;
      }

      return convertTimestampToMoment(this.value).startOf('minute').toDate();
    },

    name() {
      return this.name;
    },
  },
  inject: ['$validator'],
  components: { DatePickerField, TimePickerField },
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
  },
  computed: {
    datePickerFlexClass() {
      return {
        'pr-1': !this.reverse,
        xs8: !this.fullDay,
        xs12: this.fullDay,
      };
    },

    errorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.name, this.label));
    },
  },
  methods: {
    updateTime(time) {
      this.updateModel(updateTime(this.value, time));
    },

    updateDate(date) {
      this.updateModel(updateDate(this.value, date));
    },
  },
};
</script>
