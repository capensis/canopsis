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
      v-select.ml-4(
        v-model="range",
        :items="quickRanges",
        :label="$t('quickRanges.title')",
        hide-details,
        return-object
      )
</template>

<script>
import moment from 'moment';

import { DATETIME_FORMATS, DATETIME_INTERVAL_TYPES, QUICK_RANGES } from '@/constants';

import { findRange, dateParse } from '@/helpers/date/date-intervals';

import DatePickerField from '@/components/forms/fields/date-picker/date-picker-field.vue';

export default {
  components: {
    DatePickerField,
  },
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
        const range = findRange(this.interval.from, this.interval.to);

        return this.quickRanges.find(({ value }) => value === range.value);
      },
      set({ start, stop }) {
        if (start && stop) {
          this.$emit('input', {
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
      return dateParse(date, type, DATETIME_FORMATS.datePicker);
    },

    convertIntervalFromFieldToMoment(date) {
      return this.convertIntervalFieldToMoment(date, DATETIME_INTERVAL_TYPES.start);
    },

    convertIntervalToFieldToMoment(date) {
      return this.convertIntervalFieldToMoment(date, DATETIME_INTERVAL_TYPES.stop);
    },

    isAllowedQuickRange({ start, stop }) {
      if (!start || !stop) {
        return true;
      }

      const startMoment = this.convertIntervalFromFieldToMoment(start);
      const stopMoment = this.convertIntervalToFieldToMoment(stop);

      return this.isAllowedAccumulatedFromDate(startMoment) && this.isAllowedAccumulatedToDate(stopMoment);
    },

    isAllowedAccumulatedFromDate(dateMoment) {
      return this.accumulatedBefore > dateMoment.unix()
        /**
         * NOTE: If the date is before the accumulation date, the data is grouped by week.
         * In this case, we can only select Monday.
         */
        ? dateMoment.isoWeekday() === 1
        : true;
    },

    isAllowedFromDate(date) {
      const dateMoment = moment(date);
      const dateTimestamp = dateMoment.unix();
      const toTimestamp = this.intervalToAsMoment.unix();

      if (dateTimestamp > toTimestamp) {
        return false;
      }

      return this.isAllowedAccumulatedFromDate(dateMoment);
    },

    isAllowedAccumulatedToDate(dateMoment) {
      return this.accumulatedBefore > dateMoment.unix()
        /**
         * NOTE: If the date is before the accumulation date, the data is grouped by week.
         * In this case, we can only select Sunday.
         */
        ? dateMoment.isoWeekday() === 7
        : true;
    },

    isAllowedToDate(date) {
      const dateMoment = moment(date);
      const dateTimestamp = dateMoment.unix();
      const fromTimestamp = this.intervalFromAsMoment.unix();

      if (dateTimestamp < fromTimestamp) {
        return false;
      }

      return this.isAllowedAccumulatedToDate(dateMoment);
    },

    updateFromDate(from) {
      this.$emit('input', {
        ...this.interval,
        from,
      });
    },

    updateToDate(to) {
      this.$emit('input', {
        ...this.interval,
        to,
      });
    },
  },
};
</script>
