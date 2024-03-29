<template>
  <div>
    <v-layout>
      <v-flex xs6>
        <v-layout column>
          <date-time-picker-text-field
            v-field="value.tstart"
            v-validate="tstartRules"
            :label="$t('common.startDate')"
            :date-object-preparer="startDateObjectPreparer"
            :round-hours="roundHours"
            name="tstart"
            @update:objectValue="$emit('update:startObjectValue', $event)"
          />
          <date-time-picker-text-field
            v-field="value.tstop"
            v-validate="tstopRules"
            :label="$t('common.endDate')"
            :date-object-preparer="stopDateObjectPreparer"
            :round-hours="roundHours"
            name="tstop"
            @update:objectValue="$emit('update:stopObjectValue', $event)"
          />
        </v-layout>
      </v-flex>
      <v-flex
        class="pl-1"
        xs6
      >
        <c-quick-date-interval-type-field
          v-model="range"
          :ranges="quickRanges"
          :item-value="getQuickRangeItemValue"
          return-object
        />
        <v-select
          v-field="value.time_field"
          :items="intervalFields"
          :label="$t('quickRanges.timeField')"
          clearable
        />
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { TIME_UNITS, ALARM_INTERVAL_FIELDS, DATETIME_INTERVAL_TYPES, DATETIME_FORMATS } from '@/constants';

import { convertDateIntervalToDateObject, getValueFromQuickRange } from '@/helpers/date/date-intervals';

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
    quickRanges: {
      type: Array,
      required: false,
    },
  },
  computed: {
    range: {
      get() {
        return { start: this.value?.tstart, stop: this.value?.tstop };
      },
      set(range) {
        if (range.start !== this.range.start || range.stop !== this.range.stop) {
          this.updateModel({
            ...this.value,
            ...getValueFromQuickRange(range),
          });
        }
      },
    },

    intervalFields() {
      const messages = this.$t('quickRanges.intervalFields');

      return Object.values(ALARM_INTERVAL_FIELDS).map(value => ({
        value,
        text: messages[value],
      }));
    },

    unit() {
      return this.roundHours
        ? TIME_UNITS.hour
        : TIME_UNITS.minute;
    },
  },
  methods: {
    getQuickRangeItemValue(item) {
      return `${item.start}|${item.stop}`;
    },

    startDateObjectPreparer(date) {
      return convertDateIntervalToDateObject(
        date,
        DATETIME_INTERVAL_TYPES.start,
        DATETIME_FORMATS.dateTimePicker,
        this.unit,
      );
    },

    stopDateObjectPreparer(date) {
      return convertDateIntervalToDateObject(
        date,
        DATETIME_INTERVAL_TYPES.stop,
        DATETIME_FORMATS.dateTimePicker,
        this.unit,
      );
    },
  },
};
</script>
