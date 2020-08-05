<template lang="pug">
  ds-calendar-app(
    :calendar="calendar",
    :config="calendarConfig",
    :events="events",
    :readOnly="readOnly",
    fluid,
    fillHeight,
    @change="changeCalendarHandler",
    @changed="changedEventHandler",
    @added="changedEventHandler",
    @moved="movedEventHandler"
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
        @submit="edit",
        @remove="removePbehavior"
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
        @submit="add",
        @remove="removePbehavior"
      )
</template>

<script>
import { get } from 'lodash';
import moment from 'moment';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Schedule, Day, DaySpan, Op, Units } from 'dayspan';

import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES } from '@/constants';

import uid from '@/helpers/uid';

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
      removedPbehaviorsById: {},
      changedPbehaviorsById: {},
      addedPbehaviorsById: {},
      colorsToPbehaviors: {},
    };
  },
  computed: {
    calendarConfig() {
      return {
        dsCalendarEventTime: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
          },
        },
        dsCalendarEvent: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
          },
        },
        dsCalendarEventPlaceholder: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
          },
        },
        dsCalendarEventTimePlaceholder: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
          },
        },
      };
    },

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
      }).filter(pbehavior => !this.removedPbehaviorsById[pbehavior._id]);
    },
  },
  mounted() {
    this.fetchEvents();
  },
  methods: {
    ...mapActions({
      fetchTimespans: 'fetchItems',
    }),

    getColorForPbehavior(pbehavior = {}, color = this.$dayspan.getDefaultEventColor()) {
      if (!this.colorsToPbehaviors[pbehavior._id]) {
        this.colorsToPbehaviors[pbehavior._id] = color;
      }

      return this.colorsToPbehaviors[pbehavior._id];
    },

    async fetchEvents() {
      this.pending = true;

      const promises = this.allPbehaviors.map(pbehavior =>
        this.fetchEventsForPbehavior(pbehavior, this.getColorForPbehavior(pbehavior)));

      await Promise.all(promises);

      this.pending = false;
    },

    async fetchEventsForPbehavior(pbehavior, color = this.$dayspan.getDefaultEventColor()) {
      const viewFrom = this.calendar.filled.start.date.unix();
      const viewTo = this.calendar.filled.end.date.unix();
      const byDate = [Units.MONTH, Units.YEAR].includes(this.calendar.type);

      const timespans = await this.fetchTimespans({
        data: {
          rrule: pbehavior.rrule,
          start_at: pbehavior.tstart,
          end_at: pbehavior.tstop,
          view_from: pbehavior.tstart < viewFrom ? pbehavior.tstart : viewFrom,
          view_to: pbehavior.tstop > viewTo ? pbehavior.tstop : viewTo,
          by_date: byDate,
        },
      });

      const events = timespans.map((timespan, index) => {
        const startMoment = moment.utc(timespan.from, 'X');
        const endMoment = moment.utc(timespan.to, 'X');

        if (byDate) {
          endMoment.endOf('day');
        }

        const startDay = new Day(startMoment);
        const endDay = new Day(endMoment);
        const daySpan = new DaySpan(startDay, endDay);

        return {
          id: `${pbehavior._id}-${index}`,
          data: {
            ...this.$dayspan.getDefaultEventDetails(),

            color,
            pbehavior,
            title: pbehavior.name,
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

    removePbehavior(pbehavior) {
      if (this.addedPbehaviorsById[pbehavior._id]) {
        this.$delete(this.addedPbehaviorsById, pbehavior._id);
      } else {
        this.$set(this.removedPbehaviorsById, pbehavior._id, pbehavior);

        if (this.changedPbehaviorsById[pbehavior._id]) {
          this.$delete(this.changedPbehaviorsById, pbehavior._id);
        }
      }

      this.events = this.events.filter(event => get(event.data, 'pbehavior._id') !== pbehavior._id);
    },

    async changedEventHandler(event) {
      const { pbehavior, color } = event.calendarEvent.data;

      if (event.closePopover) {
        event.closePopover();
      }

      if (pbehavior) {
        await this.updatePbehavior(pbehavior, color);
      }

      event.clearPlaceholder();
    },

    movedEventHandler(event) {
      const pbehavior = get(event.calendarEvent, 'data.pbehavior');

      if (pbehavior) {
        this.$modals.show({
          name: MODALS.planningEventChangingConfirmation,
          config: {
            action: async (type) => {
              if (type === PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected) {
                const tstart = event.target.start.unix();
                const tstop = event.target.end.unix();

                const newPbehavior = {
                  ...pbehavior,

                  _id: uid('pbehavior'),
                  tstart,
                  tstop,
                };

                await this.updatePbehavior(newPbehavior);
              } else {
                const startDiff = event.target.start.millisBetween(event.calendarEvent.start, Op.NONE, false);
                const endDiff = event.target.end.millisBetween(event.calendarEvent.end, Op.NONE, false);

                const tstart = moment.unix(pbehavior.tstart).add(startDiff, 'millisecond').unix();
                const tstop = moment.unix(pbehavior.tstop).add(endDiff, 'millisecond').unix();

                const newPbehavior = {
                  ...pbehavior,

                  tstart,
                  tstop,
                };

                await this.updatePbehavior(newPbehavior, event.calendarEvent.data.color);
              }

              event.clearPlaceholder();
            },
            cancelAction: event.clearPlaceholder,
          },
        });
      } else {
        event.openPopover();
      }
    },

    updatePbehavior(pbehavior, color) {
      if (this.pbehaviorsById[pbehavior._id] || this.changedPbehaviorsById[pbehavior._id]) {
        this.$set(this.changedPbehaviorsById, pbehavior._id, pbehavior);
      } else {
        this.$set(this.addedPbehaviorsById, pbehavior._id, pbehavior);
      }

      return this.fetchEventsForPbehavior(pbehavior, this.getColorForPbehavior(pbehavior, color));
    },
  },
};
</script>
