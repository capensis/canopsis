<template>
  <c-layout class="date-interval-selector" gap="3">
    <c-layout gap="3" column>
      <date-time-picker-text-field
        :value="value.tstart"
        :label="$t('common.startDate')"
        :date-object-preparer="startDateObjectPreparer"
        name="tstart"
        round-hours
        autofocus
        @update:value="updateDate('tstart', $event)"
      />
      <date-time-picker-text-field
        :value="value.tstop"
        :label="$t('common.endDate')"
        :date-object-preparer="stopDateObjectPreparer"
        name="tstop"
        round-hours
        @update:value="updateDate('tstop', $event)"
      />
    </c-layout>

    <c-layout class="date-interval-selector__right" gap="3" column>
      <c-quick-date-interval-type-field
        :value="range"
        :ranges="quickRanges"
        :label="$t('quickRanges.title')"
        return-object
        @input="updateRange"
      />
      <c-select
        :model-value="value.time_field"
        :items="intervalFields"
        :label="$t('quickRanges.timeField')"
        item-title="text"
        item-value="value"
        clearable
        @update:model-value="updateField('time_field', $event)"
      />
    </c-layout>
  </c-layout>
</template>

<script>
import {
  ALARM_INTERVAL_FIELDS,
  DATETIME_FORMATS,
  DATETIME_INTERVAL_TYPES,
  TIME_UNITS,
} from '@/constants';

import {
  convertDateIntervalToDateObject,
  getValueFromQuickRange,
} from '@/helpers/date/date-intervals';

import DateTimePickerTextField from './fields/date-time-picker/date-time-picker-text-field.vue';

export default {
  components: { DateTimePickerTextField },

  props: {
    value: {
      type: Object,
      required: true,
    },
    quickRanges: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:value'],

  computed: {
    range() {
      return { start: this.value.tstart, stop: this.value.tstop };
    },

    intervalFields() {
      return Object.values(ALARM_INTERVAL_FIELDS).map(value => ({
        value,
        text: this.$t(`quickRanges.intervalFields[${JSON.stringify(value)}]`),
      }));
    },
  },

  methods: {
    startDateObjectPreparer(date) {
      return convertDateIntervalToDateObject(
        date,
        DATETIME_INTERVAL_TYPES.start,
        DATETIME_FORMATS.dateTimePicker,
        TIME_UNITS.hour,
      );
    },

    stopDateObjectPreparer(date) {
      return convertDateIntervalToDateObject(
        date,
        DATETIME_INTERVAL_TYPES.stop,
        DATETIME_FORMATS.dateTimePicker,
        TIME_UNITS.hour,
      );
    },

    updateField(field, value) {
      this.$emit('update:value', { ...this.value, [field]: value });
    },

    updateDate(field, value) {
      this.updateField(field, value || '');
    },

    updateRange(range) {
      if (range?.start && range?.stop) {
        this.$emit('update:value', {
          ...this.value,
          ...getValueFromQuickRange(range),
        });
      }
    },
  },
};
</script>

<style scoped>
.date-interval-selector {
  display: flex;
}

.date-interval-selector > :deep(*) {
  flex: 1 1 50%;
  min-width: 0;
}

.date-interval-selector__right {
  padding-left: 4px;
}
</style>
