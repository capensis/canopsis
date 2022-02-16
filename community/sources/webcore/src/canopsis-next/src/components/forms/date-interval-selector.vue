<template lang="pug">
  div
    v-layout
      v-flex(xs6)
        v-layout(column)
          date-time-picker-text-field(
            v-field="value.tstart",
            v-validate="tstartRules",
            :label="$t('common.startDate')",
            :date-object-preparer="startDateObjectPreparer",
            :round-hours="roundHours",
            name="tstart",
            @update:objectValue="$emit('update:startObjectValue', $event)"
          )
          date-time-picker-text-field(
            v-field="value.tstop",
            v-validate="tstopRules",
            :label="$t('common.endDate')",
            :date-object-preparer="stopDateObjectPreparer",
            :round-hours="roundHours",
            name="tstop",
            @update:objectValue="$emit('update:stopObjectValue', $event)"
          )
      v-flex.pl-1(xs6)
        c-quick-date-interval-type-field(v-model="range")
        v-select(
          v-field="value.time_field",
          :items="intervalFields",
          :label="$t('quickRanges.timeField')",
          clearable
        )
</template>

<script>
import { TIME_UNITS, ALARM_INTERVAL_FIELDS, DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, QUICK_RANGES } from '@/constants';

import { convertDateIntervalToDateObject } from '@/helpers/date/date-intervals';
import { convertDateToStartOfUnitString, subtractUnitFromDate } from '@/helpers/date/date';

import { formMixin } from '@/mixins/form';

import DateTimePickerTextField from '@/components/forms/fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  inject: ['$validator'],
  components: {
    DateTimePickerTextField,
  },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    roundHours: {
      type: Boolean,
      default: false,
    },
    value: {
      type: Object,
      required: true,
    },
    tstopRules: {
      type: [String, Object],
      default: null,
    },
    tstartRules: {
      type: [String, Object],
      default: null,
    },
  },
  computed: {
    range: {
      get() {
        return this.value;
      },
      set(range) {
        if (range.value !== this.range.value) {
          this.updateModel({
            ...this.value,
            ...this.getValueFromRange(range),
          });
        }
      },
    },

    intervalFields() {
      return Object.values(ALARM_INTERVAL_FIELDS).map(value => ({
        value,
        text: value,
      }));
    },

    unit() {
      return this.roundHours
        ? TIME_UNITS.hour
        : TIME_UNITS.minute;
    },
  },
  methods: {
    getValueFromRange({ value, start, stop }) {
      if (value === QUICK_RANGES.custom.value) {
        return {
          periodUnit: TIME_UNITS.hour,
          periodValue: 1,

          tstart: convertDateToStartOfUnitString(
            subtractUnitFromDate(Date.now(), 1, TIME_UNITS.hour),
            TIME_UNITS.hour,
            DATETIME_FORMATS.dateTimePicker,
          ),

          tstop: convertDateToStartOfUnitString(
            Date.now(),
            TIME_UNITS.hour,
            DATETIME_FORMATS.dateTimePicker,
          ),
        };
      }

      return {
        tstart: start,
        tstop: stop,
      };
    },

    stopDateObjectPreparer(date) {
      return convertDateIntervalToDateObject(date, DATETIME_INTERVAL_TYPES.stop, this.unit);
    },

    startDateObjectPreparer(date) {
      return convertDateIntervalToDateObject(date, DATETIME_INTERVAL_TYPES.start, this.unit);
    },
  },
};
</script>
