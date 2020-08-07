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
import moment from 'moment-timezone';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Schedule, Day, DaySpan, Op, Units, Constants } from 'dayspan';

import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES } from '@/constants';

import uid from '@/helpers/uid';

import PbehaviorCreateEvent from './partials/pbehavior-create-event.vue';

const { mapActions } = createNamespacedHelpers('pbehaviorTimespan');

function momentFromTimestampByTimezone(timestamp, sourceTimezone, localTimezone = moment.tz.guess()) {
  const result = moment.unix(timestamp);

  if (sourceTimezone === localTimezone) {
    return result;
  }

  return result.tz(sourceTimezone).tz(localTimezone, true);
}

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
      const viewFrom = moment(this.calendar.filled.start.date).tz('Europe/Paris', true).unix();
      const viewTo = moment(this.calendar.filled.end.date).tz('Europe/Paris', true).unix();
      const byDate = [Units.MONTH, Units.YEAR].includes(this.calendar.type);

      const timespans = await this.fetchTimespans({
        data: {
          rrule: pbehavior.rrule,
          start_at: pbehavior.tstart,
          end_at: pbehavior.tstop,
          view_from: (pbehavior.tstart < viewFrom && pbehavior.tstop > viewFrom) ? pbehavior.tstart : viewFrom,
          view_to: (pbehavior.tstop > viewTo && pbehavior.tstart < viewTo) ? pbehavior.tstop : viewTo,
          exdates: pbehavior.exdates,
          exceptions: pbehavior.exceptions,
          by_date: byDate,
        },
      });

      const events = timespans.map((timespan, index) => {
        const startMoment = momentFromTimestampByTimezone(timespan.from, 'Europe/Paris');
        const endMoment = momentFromTimestampByTimezone(timespan.to, 'Europe/Paris');
        if (byDate) {
          endMoment.endOf('day');
        }
        const startDay = new Day(startMoment);
        const endDay = new Day(endMoment);
        const daySpan = new DaySpan(startDay, endDay);

        function getSchedule(span) {
          const { start } = span;
          const minutes = span.minutes();
          const isDay = (minutes % Constants.MINUTES_IN_DAY === 0) || ((minutes + 1) % Constants.MINUTES_IN_DAY === 0);

          if (isDay) {
            const duration = Math.ceil(minutes / Constants.MINUTES_IN_DAY);

            return Schedule.forDay(start, duration);
          }

          const isHour = minutes % Constants.MINUTES_IN_HOUR === 0;
          const duration = isHour ? minutes / Constants.MINUTES_IN_HOUR : minutes;
          const durationUnit = isHour ? 'hours' : 'minutes';

          return Schedule.forTime(start, start.asTime(), duration, durationUnit);
        }

        return {
          id: `${pbehavior._id}-${index}`,
          data: {
            ...this.$dayspan.getDefaultEventDetails(),

            color,
            pbehavior,
            title: pbehavior.name,
          },
          schedule: getSchedule(daySpan),
        };
      });

      this.events = [
        ...this.events.filter(event => get(event.data, 'pbehavior._id') !== pbehavior._id),
        ...events,
      ];
    },

    changeCalendarHandler() {
      this.query = {
        start: moment(this.calendar.filled.start.date).tz('Europe/Paris', true).unix(),
        end: moment(this.calendar.filled.end.date).tz('Europe/Paris', true).unix(),
      };
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

    async movedEventHandler(event) {
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

      this.$modals.show({
        name: MODALS.planningEventChangingConfirmation,
        config: {
          action: async (type) => {
            if (type === PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected) {
              const tstart = moment(event.target.start.date).tz('Europe/Paris', true).unix();
              const tstop = moment(event.target.end.date).tz('Europe/Paris', true).unix();

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

              await Promise.all([
                this.updatePbehavior(mainPbehavior, event.calendarEvent.data.color),
                this.updatePbehavior(newPbehavior),
              ]);
            } else {
              const startDiff = event.target.start.minutesBetween(event.calendarEvent.start, Op.FLOOR, false);
              const endDiff = event.target.end.minutesBetween(event.calendarEvent.end, Op.FLOOR, false);
              const tstart = pbehavior.tstart + (startDiff * 60);
              const tstop = pbehavior.tstop + (endDiff * 60);
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
