<template lang="pug">
  div.c-quick-interval
    date-picker-field(
      :value="intervalFromString",
      :label="$t('common.from')",
      :allowed-dates="isAllowedFromDate",
      hide-details,
      @input="updateFromDate"
    )
      v-icon(slot="append", color="black") calendar_today
    date-picker-field.ml-4(
      :value="intervalToString",
      :label="$t('common.to')",
      :allowed-dates="isAllowedToDate",
      hide-details,
      @input="updateToDate"
    )
      v-icon(slot="append", color="black") calendar_today
    div.c-quick-interval__range
      c-quick-date-interval-type-field.ml-4(
        v-model="range",
        :custom-filter="isAllowedQuickRange",
        hide-details
      )
</template>

<script>
import { DATETIME_FORMATS, QUICK_RANGES } from '@/constants';

import { convertDateToString, convertDateToTimestamp, getNowTimestamp, getWeekdayNumber } from '@/helpers/date/date';
import {
  findQuickRangeValue,
  convertStartDateIntervalToTimestamp,
  convertStopDateIntervalToTimestamp,
} from '@/helpers/date/date-intervals';

import { formMixin } from '@/mixins/form';

import DatePickerField from '@/components/forms/fields/date-picker/date-picker-field.vue';

export default {
  components: {
    DatePickerField,
  },
  mixins: [formMixin],
  model: {
    event: 'input',
    prop: 'interval',
  },
  props: {
    interval: {
      type: Object,
      default: () => ({}),
    },
    accumulatedBefore: {
      type: Number,
      required: false,
    },
    min: {
      type: Number,
      required: false,
    },
  },
  computed: {
    quickRanges() {
      return Object.values(QUICK_RANGES)
        .filter(this.isAllowedQuickRange)
        .map(range => ({
          ...range,
          text: this.$t(`quickRanges.types.${range.value}`),
        }));
    },

    intervalFromAsTimestamp() {
      return convertStartDateIntervalToTimestamp(this.interval.from);
    },

    intervalToAsTimestamp() {
      return convertStopDateIntervalToTimestamp(this.interval.to);
    },

    intervalFromString() {
      return convertDateToString(this.intervalFromAsTimestamp, DATETIME_FORMATS.datePicker);
    },

    intervalToString() {
      return convertDateToString(this.intervalToAsTimestamp, DATETIME_FORMATS.datePicker);
    },

    range: {
      get() {
        const range = findQuickRangeValue(this.interval.from, this.interval.to);

        return this.quickRanges.find(({ value }) => value === range.value);
      },

      set({ start, stop }) {
        if (start && stop) {
          this.updateModel({
            ...this.interval,
            from: start,
            to: stop,
          });
        }
      },
    },
  },
  methods: {
    isGreaterMinDate(dateTimestamp) {
      if (this.min) {
        return dateTimestamp >= this.min;
      }

      return true;
    },

    isLessToDate(dateTimestamp) {
      return dateTimestamp < this.intervalToAsTimestamp;
    },

    isLessNowDate(dateTimestamp) {
      return dateTimestamp <= getNowTimestamp();
    },

    isGreaterFromDate(dateTimestamp) {
      return dateTimestamp > this.intervalFromAsTimestamp;
    },

    isAllowedAccumulatedFromDate(dateTimestamp) {
      return this.accumulatedBefore > dateTimestamp
        /**
         * NOTE: If the date is before the accumulation date, the data is grouped by week.
         * In this case, we can only select Monday.
         */
        ? getWeekdayNumber(dateTimestamp) === 1
        : true;
    },

    isAllowedAccumulatedToDate(dateTimestamp) {
      return this.accumulatedBefore > dateTimestamp
        /**
         * NOTE: If the date is before the accumulation date, the data is grouped by week.
         * In this case, we can only select Sunday.
         */
        ? getWeekdayNumber(dateTimestamp) === 7
        : true;
    },

    isAllowedFromDate(date) {
      const dateTimestamp = convertDateToTimestamp(date);

      return this.isLessToDate(dateTimestamp)
        && this.isGreaterMinDate(dateTimestamp)
        && this.isAllowedAccumulatedFromDate(dateTimestamp);
    },

    isAllowedToDate(date) {
      const dateTimestamp = convertDateToTimestamp(date);

      return this.isGreaterFromDate(dateTimestamp)
        && this.isLessNowDate(dateTimestamp)
        && this.isAllowedAccumulatedToDate(dateTimestamp);
    },

    isAllowedQuickRange({ start, stop }) {
      if (!start || !stop) {
        return true;
      }

      const startTimestamp = convertStartDateIntervalToTimestamp(start);
      const stopTimestamp = convertStopDateIntervalToTimestamp(start);

      return this.isGreaterMinDate(startTimestamp)
        && this.isAllowedAccumulatedFromDate(startTimestamp)
        && this.isLessNowDate(stopTimestamp)
        && this.isAllowedAccumulatedToDate(stopTimestamp);
    },

    updateFromDate(from) {
      this.updateField('from', from);
    },

    updateToDate(to) {
      this.updateField('to', to);
    },
  },
};
</script>

<style scoped lang="scss">
.c-quick-interval {
  display: inline-flex;

  &__range {
    display: flex;
    max-width: 180px;
  }
}
</style>
