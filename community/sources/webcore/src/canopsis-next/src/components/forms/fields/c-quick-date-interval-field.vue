<template lang="pug">
  v-layout(row, align-center, wrap)
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
    v-flex(xs2)
      c-quick-date-interval-type-field.ml-4(
        v-model="range",
        :custom-filter="isAllowedQuickRange",
        hide-details
      )
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, QUICK_RANGES } from '@/constants';

import { getNowTimestamp } from '@/helpers/date/date';
import { convertDateIntervalToMoment } from '@/helpers/date/date-intervals';

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

    intervalFromAsMoment() {
      return this.convertIntervalFromFieldToMoment(this.interval.from);
    },

    intervalToAsMoment() {
      return this.convertIntervalToFieldToMoment(this.interval.to);
    },

    intervalFromString() {
      return this.intervalFromAsMoment.format(DATETIME_FORMATS.datePicker);
    },

    intervalToString() {
      return this.intervalToAsMoment.format(DATETIME_FORMATS.datePicker);
    },

    range: {
      get() {
        return {
          start: this.interval.from,
          stop: this.interval.to,
        };
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
    convertIntervalFieldToMoment(date, type = DATETIME_INTERVAL_TYPES.start) {
      return convertDateIntervalToMoment(date, type, DATETIME_FORMATS.datePicker);
    },

    convertIntervalFromFieldToMoment(date) {
      return this.convertIntervalFieldToMoment(date, DATETIME_INTERVAL_TYPES.start);
    },

    convertIntervalToFieldToMoment(date) {
      return this.convertIntervalFieldToMoment(date, DATETIME_INTERVAL_TYPES.stop);
    },

    isGreaterMinDate(dateTimestamp) {
      if (this.min) {
        return dateTimestamp >= this.min;
      }

      return true;
    },

    isLessToDate(dateTimestamp) {
      return dateTimestamp < this.intervalToAsMoment.unix();
    },

    isAllowedAccumulatedFromDate(dateTimestamp, weekday) {
      return this.accumulatedBefore > dateTimestamp
        /**
         * NOTE: If the date is before the accumulation date, the data is grouped by week.
         * In this case, we can only select Monday.
         */
        ? weekday === 1
        : true;
    },

    isAllowedFromDate(date) {
      const dateMoment = moment(date);
      const dateTimestamp = dateMoment.unix();

      return this.isLessToDate(dateTimestamp)
        && this.isGreaterMinDate(dateTimestamp)
        && this.isAllowedAccumulatedFromDate(dateTimestamp, dateMoment.isoWeekday());
    },

    isAllowedAccumulatedToDate(dateTimestamp, weekday) {
      return this.accumulatedBefore > dateTimestamp
        /**
         * NOTE: If the date is before the accumulation date, the data is grouped by week.
         * In this case, we can only select Sunday.
         */
        ? weekday === 7
        : true;
    },

    isLessNowDate(dateTimestamp) {
      return dateTimestamp <= getNowTimestamp();
    },

    isGreaterFromDate(dateTimestamp) {
      return dateTimestamp > this.intervalFromAsMoment.unix();
    },

    isAllowedToDate(date) {
      const dateMoment = moment(date);
      const dateTimestamp = dateMoment.unix();

      return this.isGreaterFromDate(dateTimestamp)
        && this.isLessNowDate(dateTimestamp)
        && this.isAllowedAccumulatedToDate(dateTimestamp, dateMoment.isoWeekday());
    },

    isAllowedQuickRange({ start, stop }) {
      if (!start || !stop) {
        return true;
      }

      const startMoment = this.convertIntervalFromFieldToMoment(start);
      const stopMoment = this.convertIntervalToFieldToMoment(stop);
      const startTimestamp = startMoment.unix();
      const stopTimestamp = stopMoment.unix();

      return this.isGreaterMinDate(startTimestamp)
        && this.isAllowedAccumulatedFromDate(startTimestamp, startMoment.isoWeekday())
        && this.isLessNowDate(stopTimestamp)
        && this.isAllowedAccumulatedToDate(stopTimestamp, stopMoment.isoWeekday());
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
