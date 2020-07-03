<template lang="pug">
  ds-calendar-app(
    :calendar="calendar",
    :events="events",
    :readOnly="readOnly",
    fluid,
    fillHeight,
    @change="changeCalendarHandler",
    @changed="changedEventHandler",
    @added="changedEventHandler"
  )
    ds-calendar-event-popover(
      slot="eventPopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(
        slot-scope="{ calendarEvent, close, edit }",
        :calendarEvent="calendarEvent",
        @close="close",
        @submit="edit"
      )
    ds-calendar-event-popover(
      slot="eventCreatePopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(
        slot-scope="{ calendarEvent, close, add }",
        :calendarEvent="calendarEvent",
        @close="close",
        @submit="add"
      )
</template>

<script>
import { get } from 'lodash';
import moment from 'moment';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Schedule, Day, DaySpan } from 'dayspan';

import uuid from '@/helpers/uuid';

import PbehaviorCreateEvent from './partials/pbehavior-create-event.vue';

const { mapActions } = createNamespacedHelpers('pbehaviorTimespan');

export default {
  components: { PbehaviorCreateEvent },
  props: {
    pbehaviors: {
      type: Array,
      default: () => [],
    },
    readOnly: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      calendar: Calendar.months(),
      events: [],
      changedPbehaviorsById: {},
      addedPbehaviorsById: {},
    };
  },
  computed: {
    pbehaviorsById() {
      return this.pbehaviors.reduce((acc, pbehavior) => {
        acc[pbehavior._id] = pbehavior;

        return acc;
      }, {});
    },

    allPbehaviors() {
      return Object.values({
        ...this.pbehaviorsById,
        ...this.changedPbehaviorsById,
        ...this.addedPbehaviorsById,
      });
    },
  },
  mounted() {
    this.fetchEvents();
  },
  methods: {
    ...mapActions({
      fetchTimespans: 'fetchItems',
    }),

    async fetchEvents() {
      this.pending = true;

      const promises = this.allPbehaviors.map(pbehavior => this.fetchEventsForPbehavior(pbehavior));

      await Promise.all(promises);

      this.pending = false;
    },

    async fetchEventsForPbehavior(pbehavior, color = '#F44336') {
      const timespans = await this.fetchTimespans({
        data: {
          rrule: pbehavior.rrule,
          start_at: pbehavior.tstart,
          end_at: pbehavior.tstop,
          view_from: this.calendar.filled.start.date.unix(),
          view_to: this.calendar.filled.end.date.unix(),
        },
      });

      const events = timespans.map((timespan) => {
        const startDay = new Day(moment.unix(timespan.from));
        const endDay = new Day(moment.unix(timespan.to));
        const daySpan = new DaySpan(startDay, endDay);

        return {
          id: uuid('event'),
          data: {
            color,
            pbehavior,

            title: pbehavior.name,
            description: '',
            location: '',
            forecolor: '#ffffff',
            calendar: '',
            busy: true,
            icon: '',
          },
          schedule: Schedule.forSpan(daySpan),
        };
      });

      this.events = [
        ...this.events.filter(event => get(event.data, 'pbehavior._id') !== pbehavior._id),
        ...events,
      ];
    },

    changeCalendarHandler() {
      this.fetchEvents();
    },

    async changedEventHandler(event) {
      const { pbehavior, color } = event.calendarEvent.data;

      if (pbehavior) {
        if (this.pbehaviorsById[pbehavior._id] || this.changedPbehaviorsById[pbehavior._id]) {
          this.$set(this.changedPbehaviorsById, pbehavior._id, pbehavior);
        } else {
          this.$set(this.addedPbehaviorsById, pbehavior._id, pbehavior);
        }

        await this.fetchEventsForPbehavior(pbehavior, color);
      }

      event.clearPlaceholder();
    },
  },
};
</script>
