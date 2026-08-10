<template>
  <span class="calendar-period-picker">
    <c-select
      v-if="isWeekType"
      :model-value="week"
      :items="weeks"
      class="calendar-period-picker__week pt-0 mt-0"
      item-title="text"
      item-value="value"
      hide-details
      @update:model-value="week = $event"
    />

    <template v-else>
      <c-select
        v-if="!isMonthType"
        :model-value="day"
        :items="days"
        class="calendar-period-picker__day pt-0 mt-0"
        hide-details
        @update:model-value="day = $event"
      />

      <c-select
        :model-value="month"
        :items="months"
        class="calendar-period-picker__month pt-0 mt-0"
        item-title="text"
        item-value="value"
        hide-details
        @update:model-value="month = $event"
      />
    </template>

    <c-select
      :model-value="year"
      :items="years"
      class="calendar-period-picker__year pt-0 mt-0"
      hide-details
      @update:model-value="year = $event"
    />
  </span>
</template>

<script>
import { range } from 'lodash';

import { CALENDAR_TYPES, TIME_UNITS } from '@/constants';
import {
  convertDateToString,
  getDateByMonthNumber,
  getDateByWeekNumber,
  getDaysInMonth,
  getWeekEndDay,
  getWeekNumber,
  getWeekStartDay,
  isSameDates,
} from '@/helpers/date/date';

export default {
  props: {
    focus: {
      type: Date,
      required: true,
    },
    type: {
      type: String,
      required: true,
    },
    previousYearsOffset: {
      type: Number,
      default: 20,
    },
    nextYearsOffset: {
      type: Number,
      default: 50,
    },
  },

  emits: ['update:focus'],

  computed: {
    isWeekType() {
      return this.type === CALENDAR_TYPES.week;
    },

    isMonthType() {
      return this.type === CALENDAR_TYPES.month;
    },

    day: {
      get() {
        return this.focus.getDate();
      },
      set(value) {
        const selectDate = new Date(this.focus);

        selectDate.setDate(value);

        this.updateFocus(selectDate);
      },
    },

    month: {
      get() {
        return this.focus.getMonth();
      },
      set(value) {
        const selectDate = new Date(this.focus);
        const lastDay = getDaysInMonth(getDateByMonthNumber(selectDate, value));

        if (selectDate.getDate() > lastDay) {
          selectDate.setDate(lastDay);
        }

        selectDate.setMonth(value);

        this.updateFocus(selectDate);
      },
    },

    year: {
      get() {
        return this.focus.getFullYear();
      },
      set(value) {
        const selectDate = new Date(this.focus);

        selectDate.setFullYear(value);

        this.updateFocus(selectDate);
      },
    },

    week: {
      get() {
        return getWeekNumber(this.focus);
      },
      set(value) {
        this.updateFocus(getDateByWeekNumber(this.focus, value));
      },
    },

    months() {
      return [
        { text: this.$t('common.months.january'), value: 0 },
        { text: this.$t('common.months.february'), value: 1 },
        { text: this.$t('common.months.march'), value: 2 },
        { text: this.$t('common.months.april'), value: 3 },
        { text: this.$t('common.months.may'), value: 4 },
        { text: this.$t('common.months.june'), value: 5 },
        { text: this.$t('common.months.july'), value: 6 },
        { text: this.$t('common.months.august'), value: 7 },
        { text: this.$t('common.months.september'), value: 8 },
        { text: this.$t('common.months.october'), value: 9 },
        { text: this.$t('common.months.november'), value: 10 },
        { text: this.$t('common.months.december'), value: 11 },
      ];
    },

    years() {
      const nowYear = new Date().getFullYear();
      const start = nowYear - this.previousYearsOffset;
      const end = nowYear + this.nextYearsOffset;

      return range(start, end);
    },

    days() {
      return range(1, getDaysInMonth(this.focus) + 1);
    },

    weeks() {
      return range(1, 54).map(week => ({
        text: this.getWeekLabelByNumber(week),
        value: week,
      }));
    },
  },

  methods: {
    getWeekLabelByNumber(week) {
      const start = getWeekStartDay(this.focus, week);
      const end = getWeekEndDay(this.focus, week);

      const format = isSameDates(start, end, TIME_UNITS.year) ? 'MMMM Do' : 'MMMM Do YYYY';

      return `${convertDateToString(start, format)} - ${convertDateToString(end, format)}`;
    },

    updateFocus(date) {
      this.$emit('update:focus', date);
    },
  },
};
</script>

<style lang="scss">
.calendar-period-picker {
  display: inline-flex;
  gap: 10px;

  &__day {
    width: 70px;
  }

  &__week {
    width: 350px;
  }

  &__month {
    width: 140px;
  }

  &__year {
    width: 100px;
  }
}
</style>
