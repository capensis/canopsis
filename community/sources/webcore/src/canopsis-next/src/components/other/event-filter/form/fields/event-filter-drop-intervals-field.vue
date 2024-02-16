<template>
  <v-layout>
    <date-time-picker-text-field
      v-validate="startRules"
      :value="startString"
      :label="$t('common.start')"
      :error-message="errors.collect('start')"
      name="start"
      @input="updateStartDate"
    />
    <date-time-picker-text-field
      class="ml-2"
      v-validate="stopRules"
      :value="stopString"
      :label="$t('common.stop')"
      :error-message="errors.collect('stop')"
      name="stop"
      @input="updateStopDate"
    />
  </v-layout>
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { convertDateToDateObject, convertDateToString, convertDateToTimestamp } from '@/helpers/date/date';

import { formMixin } from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerTextField },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    startString() {
      return this.form.start && convertDateToString(this.form.start, DATETIME_FORMATS.dateTimePicker);
    },

    stopString() {
      return this.form.stop && convertDateToString(this.form.stop, DATETIME_FORMATS.dateTimePicker);
    },

    startRules() {
      return {
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },

    stopRules() {
      return {
        after: [convertDateToString(this.form.start, DATETIME_FORMATS.dateTimePicker)],
        date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
      };
    },
  },
  watch: {
    'form.rrule': function watchRRule() {
      this.errors.remove('start');
      this.errors.remove('stop');
    },
  },
  methods: {
    prepareDate(date) {
      return date ? convertDateToTimestamp(convertDateToDateObject(date, DATETIME_FORMATS.dateTimePicker)) : null;
    },

    updateStartDate(value) {
      this.updateField('start', this.prepareDate(value));
    },

    updateStopDate(value) {
      this.updateField('stop', this.prepareDate(value));
    },
  },
};
</script>
