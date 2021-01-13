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
    @added="applyEventChangesHandler",
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
        :filter="filter",
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
        :filter="filter",
        @close="close",
        @submit="add",
        @remove="removePbehavior"
      )
</template>

<script>
import moment from 'moment';
import { get, omit } from 'lodash';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Op } from 'dayspan';

import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES, PBEHAVIOR_TYPE_TYPES } from '@/constants';

import uid from '@/helpers/uid';
import { getScheduleForSpan, getSpanForTimestamps } from '@/helpers/dayspan';
import { pbehaviorToTimespanRequest } from '@/helpers/forms/timespans-pbehavior';
import { convertDateToTimestampByTimezone } from '@/helpers/date';

import entitiesInfoMixin from '@/mixins/entities/info';
import entitiesPbehaviorTimespansMixin from '@/mixins/entities/pbehavior/timespans';

import PbehaviorCreateEvent from './partials/pbehavior-create-event.vue';

const { mapActions: pbehaviorTypesMapActions } = createNamespacedHelpers('pbehaviorTypes');

export default {
  inject: ['$system'],
  components: { PbehaviorCreateEvent },
  mixins: [
    entitiesInfoMixin,
    entitiesPbehaviorTimespansMixin,
  ],
  model: {
    prop: 'pbehaviors',
    event: 'input',
  },
  props: {
    pbehaviorsById: {
      type: Object,
      required: true,
    },
    addedPbehaviorsById: {
      type: Object,
      required: true,
    },
    removedPbehaviorsById: {
      type: Object,
      required: true,
    },
    changedPbehaviorsById: {
      type: Object,
      required: true,
    },
    filter: {
      type: Object,
      required: false,
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
      colorsToPbehaviors: {},
      defaultTypes: [],
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

    allPbehaviorsById() {
      return {
        ...this.pbehaviorsById,
        ...this.addedPbehaviorsById,
        ...this.changedPbehaviorsById,
      };
    },

    allPbehaviorsAvailable() {
      return Object.values(this.allPbehaviorsById)
        .filter(pbehavior => !this.removedPbehaviorsById[pbehavior._id]);
    },
  },
  watch: {
    allPbehaviorsById: {
      immediate: true,
      handler() {
        this.setCalendarView();
      },
    },
  },
  mounted() {
    this.fetchEvents();
    this.fetchDefaultTypes();
  },
  methods: {
    ...pbehaviorTypesMapActions({
      fetchPbehaviorTypesListWithoutStore: 'fetchListWithoutStore',
    }),

    /**
     * Set calendar view to min event date
     */
    setCalendarView() {
      const startTimestamps = Object.values(this.allPbehaviorsById).map(({ tstart }) => tstart);

      if (startTimestamps.length) {
        const startTimestamp = Math.min.apply(null, startTimestamps);
        const calendarStart = moment.unix(startTimestamp);

        this.calendar.set({ around: calendarStart });
      }
    },

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

      const promises = this.allPbehaviorsAvailable.map(pbehavior =>
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
      const from = convertDateToTimestampByTimezone(this.calendar.filled.start.date, this.$system.timezone);
      const to = convertDateToTimestampByTimezone(this.calendar.filled.end.date, this.$system.timezone);

      const timespan = pbehaviorToTimespanRequest({
        pbehavior,
        from,
        to,
      });

      return this.fetchTimespansListWithoutStore({ data: timespan });
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
        const daySpan = getSpanForTimestamps({
          start: timespan.from,
          end: timespan.to,
          timezone: this.$system.timezone,
        });

        return {
          id: `${pbehavior._id}-${index}`,
          data: {
            ...this.$dayspan.getDefaultEventDetails(),

            color,
            pbehavior,
            title: pbehavior.name,
            withoutResize: !pbehavior.tstop,
          },
          schedule: getScheduleForSpan(daySpan),
        };
      });
    },

    /**
     * Fetch calendar event for pbehavior
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
        this.$emit('update:addedPbehaviorsById', omit(this.addedPbehaviorsById, [pbehavior._id]));
      } else {
        this.$emit('update:removedPbehaviorsById', {
          ...this.removedPbehaviorsById,
          [pbehavior._id]: pbehavior,
        });

        if (this.changedPbehaviorsById[pbehavior._id]) {
          this.$emit('update:changedPbehaviorsById', omit(this.changedPbehaviorsById, [pbehavior._id]));
        }
      }

      this.events = this.events.filter(event => get(event.data, 'pbehavior._id') !== pbehavior._id);
    },

    /**
     * Apply changes for the pbehavior
     *
     * @param {Object} event
     */
    async applyEventChangesHandler(event) {
      const { pbehavior, color } = event.calendarEvent.data;

      if (pbehavior) {
        await this.updatePbehavior(pbehavior, color);
      }

      if (event.closePopover) {
        event.closePopover(event);
      }

      event.clearPlaceholder();
    },

    /**
     * Apply changes for selected event only (on pbehavior)
     *
     * @param {Object} event
     */
    applyEventChangesForSelectedHandler({ target, calendarEvent }) {
      const pbehavior = get(calendarEvent, 'data.pbehavior');
      const originalPbehavior = this.allPbehaviorsById[pbehavior._id];
      const tstart = convertDateToTimestampByTimezone(target.start.date, this.$system.timezone);
      const tstop = convertDateToTimestampByTimezone(target.end.date, this.$system.timezone);
      const exdate = {
        begin: convertDateToTimestampByTimezone(calendarEvent.start.date, this.$system.timezone),
        end: convertDateToTimestampByTimezone(calendarEvent.end.date, this.$system.timezone),
        type: this.getOppositePbehaviorType(pbehavior.type),
      };

      const mainPbehavior = {
        ...originalPbehavior,

        exdates: [
          ...(originalPbehavior.exdates || []),
          exdate,
        ],
      };

      const newPbehavior = {
        ...pbehavior,

        _id: uid('pbehavior'),
        rrule: '',
        tstart,
        tstop,
      };

      return Promise.all([
        this.updatePbehavior(mainPbehavior, calendarEvent.data.color),
        this.updatePbehavior(newPbehavior),
      ]);
    },

    /**
     * Apply event changes for all events (on pbehavior)
     *
     * @param {Object} event
     * @returns {Promise<void>}
     */
    applyEventChangesForAllHandler({ target, calendarEvent }) {
      const pbehavior = get(calendarEvent, 'data.pbehavior');
      const originalPbehavior = this.allPbehaviorsById[pbehavior._id];
      const startDiff = target.start.secondsBetween(calendarEvent.start, Op.FLOOR, false);
      const endDiff = target.end.secondsBetween(calendarEvent.end, Op.FLOOR, false);
      const tstart = originalPbehavior.tstart + startDiff;
      const tstop = originalPbehavior.tstop + endDiff;
      const newPbehavior = {
        ...pbehavior,

        tstart,
        tstop,
      };

      return this.updatePbehavior(newPbehavior, calendarEvent.data.color);
    },

    /**
     * Close popover and clear placeholder for event
     *
     * @param {Object} [event = {}]
     */
    closePopoverForEvent(event = {}) { // TODO: move that
      if (event.clearPlaceholder) {
        event.clearPlaceholder();
      }

      if (event.closePopover) {
        event.closePopover(event);
      }
    },

    /**
     * Show modal window for recurrent changes confirmation
     *
     * @param {Object} event
     */
    showPbehaviorRecurrentChangesConfirmationModal(event) {
      this.$modals.show({
        name: MODALS.pbehaviorRecurrentChangesConfirmation,
        config: {
          action: async (type) => {
            try {
              if (type === PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES.selected) {
                await this.applyEventChangesForSelectedHandler(event);
              } else {
                await this.applyEventChangesForAllHandler(event);
              }

              this.closePopoverForEvent(event);
            } catch (err) {
              this.$popups.error({ text: err.description || err.message || this.$t('errors.default') });
              this.closePopoverForEvent(event);
            }
          },
          cancel: () => this.closePopoverForEvent(event),
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
      const {
        target = getSpanForTimestamps({
          start: pbehavior.tstart,
          end: pbehavior.tstop,
          timezone: this.$system.timezone,
        }),
      } = event;

      if (!pbehavior) {
        event.openPopover();

        return;
      }

      if (!pbehavior.rrule) {
        const tstart = convertDateToTimestampByTimezone(target.start.date, this.$system.timezone);
        const tstop = pbehavior.tstop
          ? convertDateToTimestampByTimezone(target.end.date, this.$system.timezone)
          : null;

        await this.updatePbehavior({
          ...pbehavior,

          tstart,
          tstop,
        }, event.calendarEvent.data.color);

        this.closePopoverForEvent(event);

        return;
      }

      this.showPbehaviorRecurrentChangesConfirmationModal({ ...event, target });
    },

    /**
     * Update pbehavior in the data and fetch timespans for that
     *
     * @param {Object} pbehavior
     * @param {string} [color = this.$dayspan.getDefaultEventColor()]
     * @returns {Promise<void>}
     */
    updatePbehavior(pbehavior, color = this.$dayspan.getDefaultEventColor()) {
      const hasPbehaviorInList = this.pbehaviorsById[pbehavior._id] || this.changedPbehaviorsById[pbehavior._id];

      if (hasPbehaviorInList) {
        this.$emit('update:changedPbehaviorsById', { ...this.changedPbehaviorsById, [pbehavior._id]: pbehavior });
      } else {
        this.$emit('update:addedPbehaviorsById', {
          ...this.addedPbehaviorsById,
          [pbehavior._id]: pbehavior,
        });
      }

      return this.fetchEventsForPbehavior(pbehavior, this.getColorForPbehavior(pbehavior, color));
    },

    /**
     * Fetch default pbehavior types
     *
     * @return {Promise<void>}
     */
    async fetchDefaultTypes() {
      const { data } = await this.fetchPbehaviorTypesListWithoutStore({
        params: { default: true },
      });

      this.defaultTypes = data;
    },

    /**
     * Get opposite pbehavior type for default exdate
     *
     * @param {Object} [pbehaviorType = {}]
     * @return {Object|undefined}
     */
    getOppositePbehaviorType(pbehaviorType = {}) {
      const TYPES_MAP = {
        [PBEHAVIOR_TYPE_TYPES.active]: PBEHAVIOR_TYPE_TYPES.inactive,
        [PBEHAVIOR_TYPE_TYPES.maintenance]: PBEHAVIOR_TYPE_TYPES.active,
        [PBEHAVIOR_TYPE_TYPES.pause]: PBEHAVIOR_TYPE_TYPES.active,
        [PBEHAVIOR_TYPE_TYPES.inactive]: PBEHAVIOR_TYPE_TYPES.active,
      };

      const targetType = TYPES_MAP[pbehaviorType.type];

      return this.defaultTypes.find(defaultType => defaultType.type === targetType);
    },
  },
};
</script>
