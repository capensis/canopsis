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
    @added="applyEventChangingHandler",
    @moved="changedEventHandler",
    @resized="changedEventHandler"
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
import moment from 'moment-timezone';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Day, DaySpan, Op, Units } from 'dayspan';

import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES } from '@/constants';

import uid from '@/helpers/uid';
import { getScheduleForSpan } from '@/helpers/dayspan';
import { convertTimestampToMomentByTimezone } from '@/helpers/date';

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
      addedPbehaviorsById: {
        '123s': {
          _id: '123s',
          tstart: 1595887200,
          tstop: 1596059999,
          name: 'asd',
          rrule: 'FREQ=WEEKLY;COUNT=12',
        },
      },
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

    isCalendarTypeWeek() {
      return [Units.MONTH, Units.YEAR].includes(this.calendar.type);
    },
  },
  mounted() {
    this.fetchEvents();
  },
  methods: {
    ...mapActions({
      fetchTimespans: 'fetchItems',
    }),

    /**
     * Get color for pbehavior and save that into data for correct displaying
     *
     * @param {Object} [pbehavior = {}]
     * @param {string} [color = this.$dayspan.getDefaultEventColor()]
     * @returns {string}
     */
    getColorForPbehavior(pbehavior = {}, color = this.$dayspan.getDefaultEventColor()) {
      if (!this.colorsToPbehaviors[pbehavior._id]) {
        this.colorsToPbehaviors[pbehavior._id] = color;
      }

      return this.colorsToPbehaviors[pbehavior._id];
    },

    /**
     * Fetch timespans and convert that into events for every pbehavior from data
     *
     * @returns {Promise<void>}
     */
    async fetchEvents() {
      this.pending = true;

      const promises = this.allPbehaviors.map(pbehavior =>
        this.fetchEventsForPbehavior(pbehavior, this.getColorForPbehavior(pbehavior)));

      await Promise.all(promises);

      this.pending = false;
    },

    /**
     * Fetch timespans for a pbehavior
     *
     * @param {Object} pbehavior
     * @returns {AxiosPromise<any>}
     */
    fetchTimespansForPbehavior(pbehavior) {
      const viewFrom = moment(this.calendar.filled.start.date).tz('Europe/Paris', true).unix();
      const viewTo = moment(this.calendar.filled.end.date).tz('Europe/Paris', true).unix();

      return this.fetchTimespans({
        data: {
          rrule: pbehavior.rrule,
          start_at: pbehavior.tstart,
          end_at: pbehavior.tstop,
          view_from: (pbehavior.tstart < viewFrom && pbehavior.tstop > viewFrom) ? pbehavior.tstart : viewFrom,
          view_to: (pbehavior.tstop > viewTo && pbehavior.tstart < viewTo) ? pbehavior.tstop : viewTo,
          exdates: pbehavior.exdates,
          exceptions: pbehavior.exceptions,
          by_date: this.isCalendarTypeWeek,
        },
      });
    },

    /**
     * Convert pbehavior timespans to events
     *
     * @param {Object} pbehavior
     * @param {Object[]} timespans
     * @param {string} color
     * @returns {Object[]}
     */
    convertTimespansToEvents({
      pbehavior,
      timespans,
      color,
    }) {
      return timespans.map((timespan, index) => {
        const startMoment = convertTimestampToMomentByTimezone(timespan.from, 'Europe/Paris');
        const endMoment = convertTimestampToMomentByTimezone(timespan.to, 'Europe/Paris');

        if (this.isCalendarTypeWeek) {
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
          schedule: getScheduleForSpan(daySpan),
        };
      });
    },

    /**
     *
     * @param {Object} pbehavior
     * @param {string} color
     * @returns {Promise<void>}
     */
    async fetchEventsForPbehavior(pbehavior, color = this.$dayspan.getDefaultEventColor()) {
      const timespans = await this.fetchTimespansForPbehavior(pbehavior);

      const events = this.convertTimespansToEvents({
        pbehavior,
        timespans,
        color,
      });

      this.events = [
        ...this.events.filter(event => get(event.data, 'pbehavior._id') !== pbehavior._id),
        ...events,
      ];
    },

    /**
     * Calendar change handler
     */
    changeCalendarHandler() {
      this.fetchEvents();
    },

    /**
     * Remove pbehavior from events
     *
     * @param {Object} pbehavior
     */
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

    /**
     * Apply changes for the pbehavior
     *
     * @param {Object} event
     */
    async applyEventChangingHandler(event) {
      const { pbehavior, color } = event.calendarEvent.data;

      if (pbehavior) {
        await this.updatePbehavior(pbehavior, color);
      }

      if (event.closePopover) {
        event.closePopover();
      }

      event.clearPlaceholder();
    },

    /**
     * Apply changes for selected event only (on pbehavior)
     *
     * @param {Object} event
     */
    applyEventChangingForSelectedHandler(event) {
      const pbehavior = get(event.calendarEvent, 'data.pbehavior');
      const target = event.target || event.calendarEvent.span;
      const tstart = moment(target.start.date).tz('Europe/Paris', true).unix();
      const tstop = moment(target.end.date).tz('Europe/Paris', true).unix();

      const newPbehavior = {
        ...pbehavior,

        _id: uid('pbehavior'),
        rrule: '',
        tstart,
        tstop,
      };

      const mainPbehavior = {
        ...pbehavior,

        exdates: [
          ...(pbehavior.exdates || []),
          {
            begin: moment(event.calendarEvent.start.date).tz('Europe/Paris', true).unix(),
            end: moment(event.calendarEvent.end.date).tz('Europe/Paris', true).unix(),
          },
        ],
      };

      return Promise.all([
        this.updatePbehavior(mainPbehavior, event.calendarEvent.data.color),
        this.updatePbehavior(newPbehavior),
      ]);
    },

    /**
     * Apply event changes for all events (on pbehavior)
     *
     * @param {Object} event
     * @returns {Promise<void>}
     */
    applyEventChangingForAllHandler(event) {
      const pbehavior = get(event.calendarEvent, 'data.pbehavior');
      const startDiff = event.target.start.secondsBetween(event.calendarEvent.start, Op.FLOOR, false);
      const endDiff = event.target.end.secondsBetween(event.calendarEvent.end, Op.FLOOR, false);
      const tstart = pbehavior.tstart + startDiff;
      const tstop = pbehavior.tstop + endDiff;
      const newPbehavior = {
        ...pbehavior,

        tstart,
        tstop,
      };

      return this.updatePbehavior(newPbehavior, event.calendarEvent.data.color);
    },

    /**
     * Show planning event changes confirmation modal window
     *
     * @param {Object} event
     */
    showPlanningEventChangingConfirmationModal(event) {
      this.$modals.show({
        name: MODALS.planningEventChangingConfirmation,
        config: {
          action: async (type) => {
            if (type === PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected) {
              await this.applyEventChangingForSelectedHandler(event);
            } else {
              await this.applyEventChangingForAllHandler(event);
            }

            event.clearPlaceholder();
          },
          cancelAction: event.clearPlaceholder,
        },
      });
    },

    /**
     * Changed event handler
     *
     * @param {Object} event
     * @returns {Promise<void>}
     */
    async changedEventHandler(event) {
      const pbehavior = get(event.calendarEvent, 'data.pbehavior');

      if (!pbehavior) {
        event.openPopover();

        return;
      }

      if (!pbehavior.rrule) {
        const tstart = moment(event.target.start.date).tz('Europe/Paris', true).unix();
        const tstop = moment(event.target.end.date).tz('Europe/Paris', true).unix();

        await this.updatePbehavior({
          ...pbehavior,

          tstart,
          tstop,
        }, event.calendarEvent.data.color);

        event.clearPlaceholder();

        return;
      }

      this.showPlanningEventChangingConfirmationModal(event);
    },

    /**
     * Update pbehavior in the data and fetch timespans for that
     *
     * @param {Object} pbehavior
     * @param {string} [color = this.$dayspan.getDefaultEventColor()]
     * @returns {Promise<void>}
     */
    updatePbehavior(pbehavior, color = this.$dayspan.getDefaultEventColor()) {
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
