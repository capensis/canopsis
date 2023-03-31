<template lang="pug">
  span.ds-calendar-app-period-picker
    template(v-if="!isYearType")
      v-select.ds-calendar-app-period-picker__week(
        v-if="isWeekType",
        v-model="week",
        :items="weeks",
        :menu-props="menuProps"
      )
      template(v-else)
        v-select.ds-calendar-app-period-picker__day(
          v-if="!isMonthType",
          v-model="day",
          :items="days",
          :menu-props="menuProps"
        )
        v-select.ds-calendar-app-period-picker__month(
          v-model="month",
          :items="months",
          :menu-props="menuProps"
        )
    v-select.ds-calendar-app-period-picker__year(
      v-model="year",
      :items="years",
      :menu-props="menuProps"
    )
</template>

<script>
import { Calendar, Units } from 'dayspan';
import { range } from 'lodash';

import { TIME_UNITS } from '@/constants';

import { getDiffBetweenDates } from '@/helpers/date/date';

export default {
  props: {
    calendar: {
      type: Calendar,
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
  computed: {
    menuProps() {
      return {
        auto: true,
      };
    },

    isWeekType() {
      return this.calendar.type === Units.WEEK;
    },

    isYearType() {
      return this.calendar.type === Units.YEAR;
    },

    isMonthType() {
      return this.calendar.type === Units.MONTH;
    },

    start() {
      return this.calendar.start;
    },

    startMoment() {
      return this.start.date;
    },

    day: {
      get() {
        return this.start.dayOfMonth;
      },
      set(value) {
        const selectDate = this.startMoment.toDate();
        selectDate.setDate(value);

        this.updateCalendar(selectDate);
      },
    },

    month: {
      get() {
        return this.start.month;
      },
      set(value) {
        const selectDate = this.startMoment.toDate();
        const lastDay = this.startMoment.clone()
          .month(value)
          .daysInMonth();

        if (selectDate.getDate() > lastDay) {
          selectDate.setDate(lastDay);
        }

        selectDate.setMonth(value);

        this.updateCalendar(selectDate);
      },
    },

    year: {
      get() {
        return this.start.year;
      },
      set(value) {
        const selectDate = this.startMoment.toDate();
        selectDate.setYear(value);

        this.updateCalendar(selectDate);
      },
    },

    week: {
      get() {
        return this.start.week;
      },
      set(value) {
        const selectDate = this.startMoment.clone().week(value).toDate();

        this.updateCalendar(selectDate);
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
      return range(1, this.startMoment.daysInMonth() + 1);
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
      const start = this.startMoment.clone()
        .week(week)
        .startOf(TIME_UNITS.week);
      const end = this.startMoment.clone()
        .week(week)
        .endOf(TIME_UNITS.week);

      const format = start.isSame(end, TIME_UNITS.year) ? 'MMMM Do' : 'MMMM Do YYYY';

      return `${start.format(format)} - ${end.format(format)}`;
    },

    updateCalendar(date) {
      const unit = {
        [Units.DAY]: TIME_UNITS.day,
        [Units.WEEK]: TIME_UNITS.week,
        [Units.MONTH]: TIME_UNITS.month,
      }[this.calendar.type];

      const diff = getDiffBetweenDates(date, this.startMoment, unit);

      this.$emit('change', diff);
    },
  },
};
</script>

<style lang="scss">
.ds-calendar-app-period-picker {
  display: inline-flex;
  gap: 10px;

  &__day {
    width: 50px;
  }

  &__week {
    width: 350px;
  }

  &__month {
    width: 100px;
  }

  &__year {
    width: 80px;
  }
}
</style>
