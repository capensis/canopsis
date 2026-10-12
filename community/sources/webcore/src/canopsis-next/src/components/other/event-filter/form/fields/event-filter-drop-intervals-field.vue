<template>
  <v-layout class="gap-2">
    <date-time-picker-text-field
      v-validate="startRules"
      :value="startString"
      :label="$t('common.start')"
      :error-message="validator.errors.collect('start')"
      name="start"
      @input="updateStartDate"
    />
    <date-time-picker-text-field
      v-validate="stopRules"
      :value="stopString"
      :label="$t('common.stop')"
      :error-message="validator.errors.collect('stop')"
      name="stop"
      @input="updateStopDate"
    />
  </v-layout>
</template>

<script>
import { computed, watch } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import { convertDateToDateObject, convertDateToString, convertDateToTimestamp } from '@/helpers/date/date';

import { useModelField } from '@/hooks/form/model-field';
import { useValidator } from '@/hooks/validator/validator';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  inject: ['$validator'],
  components: { DateTimePickerTextField },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();
    const { updateField } = useModelField(props, emit);

    const startString = computed(
      () => props.form.start && convertDateToString(props.form.start, DATETIME_FORMATS.dateTimePicker),
    );

    const stopString = computed(
      () => props.form.stop && convertDateToString(props.form.stop, DATETIME_FORMATS.dateTimePicker),
    );

    const startRules = computed(() => ({
      required: props.required,
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const stopRules = computed(() => ({
      required: props.required,
      after: [convertDateToString(props.form.start, DATETIME_FORMATS.dateTimePicker)],
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    /**
     * Parses a date-time picker string into a Unix timestamp for the form model
     *
     * @param {string|null|undefined} date - Raw value from the date-time picker
     * @returns {number|null}
     */
    const prepareDate = date => (
      date ? convertDateToTimestamp(convertDateToDateObject(date, DATETIME_FORMATS.dateTimePicker)) : null
    );

    /**
     * Writes the interval start timestamp to `form.start` and emits `input`
     *
     * @param {string|null|undefined} value - Raw value from the start date-time picker
     */
    const updateStartDate = value => updateField('start', prepareDate(value));

    /**
     * Writes the interval end timestamp to `form.stop` and emits `input`
     *
     * @param {string|null|undefined} value - Raw value from the stop date-time picker
     */
    const updateStopDate = value => updateField('stop', prepareDate(value));

    watch(() => props.form.rrule, () => {
      validator.errors.remove('start');
      validator.errors.remove('stop');
    });

    return {
      validator,
      startString,
      stopString,
      startRules,
      stopRules,
      updateStartDate,
      updateStopDate,
    };
  },
};
</script>
