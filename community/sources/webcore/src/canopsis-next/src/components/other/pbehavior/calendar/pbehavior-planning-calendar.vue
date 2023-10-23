<template>
  <div class="fill-height">
    <c-calendar
      ref="calendar"
      :events="events"
      :config="calendarConfig"
      :read-only="readOnly"
      fluid
      fill-height
      current-time-for-today
      remove-events-before-move
      @changed="changedEventHandler"
      @added="applyEventChangesHandler"
      @moved="changedEventHandler"
      @resized="changedEventHandler"
      @change="fetchEvents"
      @click:event="showEventDetails"
    >
      <template #loader="">
        <c-progress-overlay
          class="calendar-progress"
          :pending="pending"
        />
      </template>
      <template #eventPopover="props">
        <ds-calendar-event-popover v-bind="props">
          <template #default="{ calendarEvent, close, edit }">
            <pbehavior-create-event
              :calendar-event="calendarEvent"
              :entity-pattern="entityPattern"
              @close="close"
              @submit="edit"
              @remove="removePbehavior"
            />
          </template>
        </ds-calendar-event-popover>
      </template>
      <template #menu-right="">
        <pbehavior-planning-calendar-legend :exception-types="exceptionTypes" />
      </template>
    </c-calendar>

    <v-menu
      v-if="selectedOpen"
      :value="true"
      :close-on-content-click="false"
      :activator="selectedElement"
      ignore-click-outside
      content-class="pbehavior-calendar-event-popover__wrapper"
    >
      <v-card>
        <v-card-text>
          <pbehavior-create-event
            :event="selectedEvent"
            :entity-pattern="entityPattern"
            @close="closeEventDetails"
            @submit="applyEventChangesHandler"
            @remove="removePbehavior"
          />
        </v-card-text>
      </v-card>
    </v-menu>
  </div>
</template>

<script>
import { get, omit, uniqBy } from 'lodash';
import { createNamespacedHelpers } from 'vuex';
import { Op } from 'dayspan';

import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES, PBEHAVIOR_TYPE_TYPES } from '@/constants';
import { CSS_COLORS_VARS } from '@/config';

import { uid } from '@/helpers/uid';
import { getMostReadableTextColor } from '@/helpers/color';
import { getSpanForTimestamps } from '@/helpers/calendar/dayspan';
import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';
import {
  convertDateToTimestampByTimezone,
  convertDateToDateObjectByTimezone,
  convertDateToDateObject,
} from '@/helpers/date/date';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesPbehaviorTimespansMixin } from '@/mixins/entities/pbehavior/timespans';

import PbehaviorCreateEvent from './partials/pbehavior-create-event.vue';
import PbehaviorPlanningCalendarLegend from './partials/pbehavior-planning-calendar-legend.vue';

const { mapActions: pbehaviorTypesMapActions } = createNamespacedHelpers('pbehaviorTypes');

export default {
  inject: ['$system'],
  components: {
    PbehaviorCreateEvent,
    PbehaviorPlanningCalendarLegend,
  },
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
    entityPattern: {
      type: Array,
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
      exceptionTypes: [],
      events: [],
      defaultTypes: [],

      selectedElement: null,
      selectedOpen: false,
      selectedEvent: null,
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
            contentClass: 'pbehavior-calendar-event-popover__wrapper',
          },
        },
        dsCalendarEvent: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
            contentClass: 'pbehavior-calendar-event-popover__wrapper',
          },
        },
        dsCalendarEventPlaceholder: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
            contentClass: 'pbehavior-calendar-event-popover__wrapper',
          },
        },
        dsCalendarEventTimePlaceholder: {
          popoverProps: {
            openOnHover: false,
            ignoreClickOutside: true,
            ignoreClickUpperOutside: true,
            contentClass: 'pbehavior-calendar-event-popover__wrapper',
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
    pbehaviorsById: {
      immediate: true,
      handler() {
        this.$nextTick(this.setCalendarView);
      },
    },
  },
  mounted() {
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
        const startDate = convertDateToDateObject(startTimestamp);

        if (this.$refs.calendar.focus.getTime() > startDate.getTime()) {
          this.$refs.calendar.setFocusDate(startDate);
        }
      }
    },

    showEventDetails({ nativeEvent, event }) {
      if (this.selectedOpen) {
        return;
      }

      this.selectedEvent = event;
      this.selectedElement = nativeEvent.target;

      this.selectedOpen = true;

      nativeEvent.stopPropagation();
    },

    closeEventDetails() {
      this.selectedOpen = false;
      this.selectedElement = null;
      this.selectedEvent = null;
    },

    /**
     * Fetch timespans and convert that into events for every pbehavior from data
     *
     * @returns {Promise<void>}
     */
    async fetchEvents({ start, end }) {
      this.pending = true;

      const promises = this.allPbehaviorsAvailable
        .map(pbehavior => this.fetchEventsForPbehavior(pbehavior, { start, end }));

      await Promise.all(promises);

      this.pending = false;
    },

    /**
     * Fetch timespans for a pbehavior
     *
     * @param {Object} pbehavior
     * @returns {AxiosPromise<any>}
     */
    fetchTimespansForPbehavior(pbehavior, { start, end }) {
      const from = convertDateToTimestampByTimezone(start, this.$system.timezone);
      const to = convertDateToTimestampByTimezone(end, this.$system.timezone);

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
     * @returns {Object[]}
     */
    convertTimespansToEvents({
      pbehavior,
      timespans,
    }) {
      return timespans.map((timespan, index) => {
        const type = timespan.type ?? pbehavior.type;

        /**
         * If there is `type` field in timespan it means that timespan is exception date with a `type`
         */
        const color = pbehavior.color || type.color || pbehavior.type?.color || CSS_COLORS_VARS.secondary;
        const iconColor = getMostReadableTextColor(color, { level: 'AA', size: 'large' });

        const start = convertDateToDateObjectByTimezone(timespan.from, this.$system.timezone);
        const end = convertDateToDateObjectByTimezone(timespan.to, this.$system.timezone);

        return {
          id: `${pbehavior._id}-${index}`,
          pbehavior,
          color,
          iconColor,
          start,
          end,
          icon: type.icon_name,
          name: pbehavior.name,
          data: {
            ...this.$dayspan.getDefaultEventDetails(),

            pbehavior,
            color,
            icon: type.icon_name,
            title: pbehavior.name,
            withoutResize: !pbehavior.tstop,
          },
        };
      });
    },

    /**
     * Fetch calendar event for pbehavior
     *
     * @param {Object} pbehavior
     * @returns {Promise<void>}
     */
    async fetchEventsForPbehavior(pbehavior, { start, end }) {
      const timespans = await this.fetchTimespansForPbehavior(pbehavior, { start, end });

      const events = this.convertTimespansToEvents({
        pbehavior,
        timespans,
      });

      const exceptionTypes = timespans.filter(({ type }) => type).map(({ type }) => type);

      this.exceptionTypes = uniqBy(exceptionTypes, '_id');
      this.events = [
        ...this.events.filter(event => get(event.data, 'pbehavior._id') !== pbehavior._id),
        ...events,
      ];
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
     * @returns {Promise<void>}
     */
    async updatePbehavior(pbehavior) {
      this.pending = true;

      const hasPbehaviorInList = this.pbehaviorsById[pbehavior._id] || this.changedPbehaviorsById[pbehavior._id];

      if (hasPbehaviorInList) {
        this.$emit('update:changedPbehaviorsById', { ...this.changedPbehaviorsById, [pbehavior._id]: pbehavior });
      } else {
        this.$emit('update:addedPbehaviorsById', {
          ...this.addedPbehaviorsById,
          [pbehavior._id]: pbehavior,
        });
      }

      await this.fetchEventsForPbehavior(pbehavior);

      this.pending = false;
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

<style lang="scss">
// We've turned off animation because of render freezes on css animation on the render process

.calendar-progress .v-progress-circular--indeterminate .v-progress-circular__overlay {
  animation: none;
}

.pbehavior-calendar-event-popover__wrapper {
  max-height: 95%;
  max-width: 95% !important;
  width: 980px !important;
  top: 50% !important;
  transform: translate3d(0, -50%, 0);
}
</style>
