<template lang="pug">
  ds-calendar-app(
    :calendar="calendar",
    :events="events",
    :pending="pending",
    read-only,
    fluid,
    fill-height,
    current-time-for-today,
    @change="fetchTimespans"
  )
</template>

<script>
import { Calendar } from 'dayspan';

import { COLORS } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import { getScheduleForSpan, getSpanForTimestamps } from '@/helpers/calendar/dayspan';
import { pbehaviorToTimespanRequest } from '@/helpers/forms/timespans-pbehavior';
import { convertDateToMoment, convertDateToString, convertDateToTimestampByTimezone } from '@/helpers/date/date';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesPbehaviorTimespansMixin } from '@/mixins/entities/pbehavior/timespans';

export default {
  inject: ['$system'],
  mixins: [
    entitiesInfoMixin,
    entitiesPbehaviorTimespansMixin,
  ],
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      timespans: [],
      calendar: Calendar.months(),
    };
  },
  computed: {
    events() {
      return this.timespans.map((timespan, index) => {
        const daySpan = getSpanForTimestamps({
          start: timespan.from,
          end: timespan.to,
          timezone: this.$system.timezone,
        });
        const fromString = convertDateToString(timespan.from, DATETIME_FORMATS.medium);
        const toString = convertDateToString(timespan.to, DATETIME_FORMATS.medium);

        return {
          id: index,
          data: {
            ...this.$dayspan.getDefaultEventDetails(),
            title: `${fromString} - ${toString}`,
            color: COLORS.secondary,
          },
          schedule: getScheduleForSpan(daySpan),
        };
      });
    },
  },
  mounted() {
    this.calendar.set({ around: convertDateToMoment(this.pbehavior.tstart) });

    this.fetchTimespans();
  },
  methods: {
    async fetchTimespans() {
      this.pending = true;

      const from = convertDateToTimestampByTimezone(this.calendar.filled.start.date, this.$system.timezone);
      const to = convertDateToTimestampByTimezone(this.calendar.filled.end.date, this.$system.timezone);

      const timespan = pbehaviorToTimespanRequest({
        pbehavior: this.pbehavior,
        from,
        to,
      });

      this.timespans = await this.fetchTimespansListWithoutStore({ data: timespan });

      this.pending = false;
    },
  },
};
</script>
