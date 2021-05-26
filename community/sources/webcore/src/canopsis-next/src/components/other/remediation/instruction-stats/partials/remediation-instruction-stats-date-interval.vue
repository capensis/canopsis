<template lang="pug">
  v-layout(row, align-center, wrap)
    date-picker-field(
      :value="intervalFromString",
      :label="$t('common.from')",
      :allowed-dates="allowedFromDates",
      hide-details,
      @input="updateFromDate"
    )
      v-icon(slot="append", color="black") calendar_today
    date-picker-field.ml-4(
      :value="intervalToString",
      :label="$t('common.to')",
      :allowed-dates="allowedToDates",
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
  },
  computed: {
    quickRanges() {
      return Object.values(QUICK_RANGES).map(range => ({
        ...range,
        text: this.$t(`quickRanges.types.${range.value}`),
      }));
    },

    intervalFromAsMoment() {
      return dateParse(this.interval.from, DATETIME_INTERVAL_TYPES.start, DATETIME_FORMATS.datePicker);
    },

    intervalToAsMoment() {
      return dateParse(this.interval.to, DATETIME_INTERVAL_TYPES.stop, DATETIME_FORMATS.datePicker);
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
    allowedFromDates(date) {
      return moment(date).unix() < this.intervalToAsMoment.unix();
    },

    allowedToDates(date) {
      return moment(date).unix() > this.intervalFromAsMoment.unix();
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
