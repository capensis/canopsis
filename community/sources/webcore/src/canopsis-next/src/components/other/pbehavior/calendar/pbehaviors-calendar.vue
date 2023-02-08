<template lang="pug">
  ds-calendar-app(
    :calendar="calendar",
    :pending="pending",
    read-only,
    fluid,
    fill-height,
    current-time-for-today,
    @change="fetchEvents"
  )
    template(#calendarAppLoader="")
      c-progress-overlay(:pending="pending")
</template>

<script>
import { Calendar } from 'dayspan';

import { DATETIME_FORMATS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';
import { getPbehaviorColor } from '@/helpers/entities/pbehavior';
import { getScheduleForSpan, getSpanForTimestamps } from '@/helpers/calendar/dayspan';
import { convertDateToString, convertDateToTimestampByTimezone } from '@/helpers/date/date';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesPbehaviorTimespansMixin } from '@/mixins/entities/pbehavior/timespans';

export default {
  inject: ['$system'],
  mixins: [
    entitiesPbehaviorMixin,
    entitiesPbehaviorTimespansMixin,
  ],
  props: {
    entityId: {
      type: String,
      required: false,
    },
  },
  data() {
    const calendar = Calendar.months();

    calendar.set({ listTimes: false });

    return {
      calendar,

      pending: false,
      events: [],
    };
  },
  computed: {
    viewInterval() {
      return {
        from: convertDateToTimestampByTimezone(this.calendar.filled.start.date, this.$system.timezone),
        to: convertDateToTimestampByTimezone(this.calendar.filled.end.date, this.$system.timezone),
      };
    },
  },
  mounted() {
    this.fetchEvents();
  },
  methods: {
    async applyEvents() {
      if (this.events) {
        this.calendar.removeEvents(null, true);
        await this.calendar.addEventsAsync(this.events);
      }
    },

    convertPbehaviorsCalendarToEvents(pbehaviors) {
      const defaultEvent = this.$dayspan.getDefaultEventDetails();

      return pbehaviors.map((pbehavior, index) => {
        const daySpan = getSpanForTimestamps({
          start: pbehavior.from,
          end: pbehavior.to,
          timezone: this.$system.timezone,
        });
        const fromString = convertDateToString(pbehavior.from, DATETIME_FORMATS.medium);
        const toString = convertDateToString(pbehavior.to, DATETIME_FORMATS.medium);
        const color = getPbehaviorColor(pbehavior);
        const forecolor = getMostReadableTextColor(color, { level: 'AA', size: 'large' });

        return {
          id: index,
          data: {
            ...defaultEvent,

            color,
            forecolor,
            title: `${fromString} - ${toString} ${pbehavior.title}`,
            icon: pbehavior.type.icon_name,
          },
          schedule: getScheduleForSpan(daySpan),
        };
      });
    },

    fetchPbehaviorsCalendar() {
      const params = {
        from: this.viewInterval.from,
        to: this.viewInterval.to,
      };

      if (this.entityId) {
        params._id = this.entityId;

        return this.fetchEntitiesPbehaviorsCalendarWithoutStore({ params });
      }

      return this.fetchPbehaviorsCalendarWithoutStore({ params });
    },

    async fetchEvents() {
      this.pending = true;

      const pbehaviorsCalendar = await this.fetchPbehaviorsCalendar();

      this.events = this.convertPbehaviorsCalendarToEvents(pbehaviorsCalendar);

      await this.applyEvents();

      this.pending = false;
    },
  },
};
</script>
