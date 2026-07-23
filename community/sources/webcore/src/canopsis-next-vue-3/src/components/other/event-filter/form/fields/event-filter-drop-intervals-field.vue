<template>
  <c-layout gap="3">
    <date-time-picker-text-field
      v-validate="startRules"
      :value="startString"
      :label="$t('common.start')"
      :error-message="validator.errors.collect('start')"
      name="start"
      @update:value="updateStartDate"
    />
    <date-time-picker-text-field
      v-validate="stopRules"
      :value="stopString"
      :label="$t('common.stop')"
      :error-message="validator.errors.collect('stop')"
      name="stop"
      @update:value="updateStopDate"
    />
  </c-layout>
</template>

<script>
import { computed, watch } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import {
  convertDateToDateObject,
  convertDateToString,
  convertDateToTimestamp,
} from '@/helpers/date/date';

import { useValidator } from '@/hooks/validator/validator';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  emits: ['update:start', 'update:stop'],
  components: { DateTimePickerTextField },
  props: {
    start: {
      type: Number,
      default: null,
    },
    stop: {
      type: Number,
      default: null,
    },
    required: {
      type: Boolean,
      default: false,
    },
    rrule: {
      type: Object,
      default: null,
    },
  },
  setup(props, { emit }) {
    const validator = useValidator();

    const startString = computed(
      () => props.start && convertDateToString(props.start, DATETIME_FORMATS.dateTimePicker),
    );

    const stopString = computed(
      () => props.stop && convertDateToString(props.stop, DATETIME_FORMATS.dateTimePicker),
    );

    const startRules = computed(() => ({
      required: props.required,
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    const stopRules = computed(() => ({
      required: props.required,
      after: [convertDateToString(props.start, DATETIME_FORMATS.dateTimePicker)],
      date_format: DATETIME_FORMATS.veeValidateDateTimeFormat,
    }));

    /**
     * Parses a date-time picker string into a Unix timestamp for the form model
     *
     * @param {string|null|undefined} date - Raw value from the date-time picker
     * @returns {number|null}
     */
    const prepareDate = (date) =>
      date
        ? convertDateToTimestamp(convertDateToDateObject(date, DATETIME_FORMATS.dateTimePicker))
        : null;

    /**
     * Writes the interval start timestamp and emits `update:start`
     *
     * @param {string|null|undefined} value - Raw value from the start date-time picker
     */
    const updateStartDate = (value) => emit('update:start', prepareDate(value));

    /**
     * Writes the interval end timestamp and emits `update:stop`
     *
     * @param {string|null|undefined} value - Raw value from the stop date-time picker
     */
    const updateStopDate = (value) => emit('update:stop', prepareDate(value));

    watch(
      () => props.rrule,
      () => {
        validator.errors.remove('start');
        validator.errors.remove('stop');
      },
    );

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
