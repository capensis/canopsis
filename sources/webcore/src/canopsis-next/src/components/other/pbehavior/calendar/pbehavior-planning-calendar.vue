<template lang="pug">
  ds-calendar-app(
    :calendar="calendar",
    :events="events",
    :readOnly="readOnly",
    fluid,
    fillHeight,
    @change="changeCalendarHandler",
    @changed="changedEventHandler",
    @added="addedEventHandler"
  )
    ds-calendar-event-popover(
      slot="eventPopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(
        slot-scope="{ placeholder, close, edit }",
        :placeholder="placeholder",
        @close="close",
        @submit="edit"
      )
    ds-calendar-event-popover(
      slot="eventCreatePopover",
      slot-scope="props",
      v-bind="props"
    )
      pbehavior-create-event(
        slot-scope="{ placeholder, close, add }",
        :placeholder="placeholder",
        @close="close",
        @submit="add"
      )
</template>

<script>
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
      changedPbehaviors: {},
      createdPbehaviors: {},
    };
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

      const promises = this.pbehaviors.map(pbehavior => this.fetchEventsForPbehavior(pbehavior));

      await Promise.all(promises);

      this.pending = false;
    },

    async fetchEventsForPbehavior(pbehavior) {
      const timespans = await this.fetchTimespans({
        data: {
          rrule: pbehavior.rrule,
          start_at: pbehavior.start_at,
          end_at: pbehavior.end_at,
          view_from: this.calendar.start.date.unix(),
          view_to: this.calendar.end.date.unix(),
        },
      });

      const events = timespans.map((timespan) => {
        const startDay = new Day(moment.unix(timespan.from));
        const endDay = new Day(moment.unix(timespan.to));
        const daySpan = new DaySpan(startDay, endDay);

        return {
          id: uuid('event'),
          data: {
            title: pbehavior.title,
            description: '',
            location: '',
            color: '#F44336',
            forecolor: '#ffffff',
            calendar: '',
            busy: true,
            icon: '',
            meta: {
              pbehavior,
            },
          },
          schedule: Schedule.forSpan(daySpan),
        };
      });

      this.events = [
        ...this.events.filter(event => event.meta.pbehavior._id !== pbehavior._id),
        ...events,
      ];
    },

    changeCalendarHandler() {
      this.fetchEvents();
    },

    changedEventHandler(event) {
      this.changedPbehaviors[event.pbehavior._id] = event.pbehavior;
      this.fetchEventsForPbehavior(event.pbehavior);

      event.clearPlaceholder();
    },
    addedEventHandler(event) {
      this.createdPbehaviors[event.pbehavior._id] = event.pbehavior;
      this.fetchEventsForPbehavior(event.pbehavior);

      event.clearPlaceholder();
    },
  },
};
</script>
