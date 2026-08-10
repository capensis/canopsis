<template>
  <div class="fill-height">
    <c-calendar
      ref="calendar"
      :events="events"
      :loading="pending"
      :timezone="timezone"
      :no-timezone="!shownUserTimezone"
      :readonly="readOnly"
      color="primary"
      @update:timezone="handleUpdateTimezone"
      @change:event="handleUpdateEvent"
      @move:event="handleUpdateEvent"
      @resize:event="handleUpdateEvent"
      @change:pagination="debouncedFetchEvents"
    >
      <template #form-event="{ event, close }">
        <pbehavior-create-event
          :event="event"
          :entity-pattern="entityPattern"
          :default-fields="defaultFields"
          :with-inherited="withInherited"
          :timezone="timezone"
          :no-timezone="!shownUserTimezone"
          @close="close"
          @submit="addEventWithClose($event, close)"
          @remove="removePbehavior"
        />
      </template>

      <template #menu-right>
        <pbehavior-planning-calendar-legend :exception-types="exceptionTypes" />
      </template>
    </c-calendar>
  </div>
</template>

<script>
import { debounce, get, omit, uniqBy } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { MODALS, PBEHAVIOR_PLANNING_EVENT_CHANGING_TYPES, PBEHAVIOR_TYPE_TYPES } from '@/constants';
import { getMostReadableTextColor } from '@/helpers/color';
import {
  convertDateToDateObject,
  convertDateToDateObjectByTimezone,
  convertDateToTimestampByTimezone,
  getDiffBetweenDates,
  getLocalTimezone,
} from '@/helpers/date/date';
import { isFullDayEvent } from '@/helpers/entities/pbehavior/form';
import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';
import { uid } from '@/helpers/uid';
import { useModals } from '@/hooks/modals';
import { usePopups } from '@/hooks/popups';
import { useInfo } from '@/hooks/store/modules/info';
import { usePbehaviorTimespan } from '@/hooks/store/modules/pbehavior-timespan';
import { usePbehaviorType } from '@/hooks/store/modules/pbehavior-type';

import PbehaviorCreateEvent from './partials/pbehavior-create-event.vue';
import PbehaviorPlanningCalendarLegend from './partials/pbehavior-planning-calendar-legend.vue';

export default {
  components: {
    PbehaviorCreateEvent,
    PbehaviorPlanningCalendarLegend,
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
    defaultFields: {
      type: Object,
      required: false,
    },
    withInherited: {
      type: Boolean,
      default: false,
    },
    timezone: {
      type: String,
      default: getLocalTimezone(),
    },
  },

  emits: [
    'update:addedPbehaviorsById',
    'update:removedPbehaviorsById',
    'update:changedPbehaviorsById',
    'update:timezone',
  ],

  data() {
    return {
      pending: false,
      exceptionTypes: [],
      events: [],
      defaultTypes: [],
    };
  },

  setup() {
    const { shownUserTimezone } = useInfo();
    const modals = useModals();
    const popups = usePopups();
    const { fetchPbehaviorTimespansListWithoutStore } = usePbehaviorTimespan();
    const { fetchPbehaviorTypesFieldList } = usePbehaviorType();

    return {
      shownUserTimezone,
      modals,
      popups,
      fetchPbehaviorTimespansListWithoutStore,
      fetchPbehaviorTypesFieldList,
    };
  },

  computed: {
    allPbehaviorsById() {
      return {
        ...this.pbehaviorsById,
        ...this.addedPbehaviorsById,
        ...this.changedPbehaviorsById,
      };
    },

    allPbehaviorsAvailable() {
      return Object.values(this.allPbehaviorsById).filter(
        pbehavior => !this.removedPbehaviorsById[pbehavior._id],
      );
    },
  },

  watch: {
    pbehaviorsById: {
      immediate: true,
      handler() {
        this.$nextTick(this.setCalendarView);
      },
    },
    timezone: 'fetchEvents',
  },

  created() {
    this.debouncedFetchEvents = debounce(this.fetchEvents, 300);
  },

  mounted() {
    this.debouncedFetchEvents();
    this.fetchDefaultTypes();
  },

  methods: {
    /**
     * Set calendar view to min event date.
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

    /**
     * Fetch timespans and convert that into events for every pbehavior from data.
     */
    async fetchEvents() {
      this.pending = true;

      const promises = this.allPbehaviorsAvailable.map(pbehavior =>
        this.fetchEventsForPbehavior(pbehavior),
      );

      await Promise.all(promises);

      this.pending = false;
    },

    /**
     * Fetch timespans for a pbehavior.
     *
     * @param {Object} pbehavior
     * @returns {Promise<Array>}
     */
    fetchTimespansForPbehavior(pbehavior) {
      const { start, end } = this.$refs.calendar.filled;

      const from = convertDateToTimestampByTimezone(start, this.timezone, true);
      const to = convertDateToTimestampByTimezone(end, this.timezone, true);

      const timespan = pbehaviorToTimespanRequest({
        pbehavior,
        from,
        to,
      });

      return this.fetchPbehaviorTimespansListWithoutStore({ data: timespan });
    },

    /**
     * Convert pbehavior timespans to events.
     *
     * @param {Object} pbehavior
     * @param {Object[]} timespans
     * @returns {Object[]}
     */
    convertTimespansToEvents({ pbehavior, timespans }) {
      const isTimed = !isFullDayEvent(
        convertDateToDateObjectByTimezone(pbehavior.tstart, this.timezone),
        pbehavior.tstop && convertDateToDateObjectByTimezone(pbehavior.tstop, this.timezone),
      );

      return timespans.map((timespan, index) => {
        const type = timespan.type ?? pbehavior.type;
        const color =
          pbehavior.color || type.color || pbehavior.type?.color || CSS_COLORS_VARS.secondary;
        const iconColor = getMostReadableTextColor(color, { level: 'AA', size: 'large' });
        const start = convertDateToDateObjectByTimezone(timespan.from, this.timezone);
        const end = convertDateToDateObjectByTimezone(timespan.to, this.timezone);

        return {
          id: `${pbehavior._id}-${index}`,
          color,
          iconColor,
          start,
          end,
          icon: type.icon_name,
          name: pbehavior.name,
          timed: isTimed,
          data: {
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
     * Fetch calendar event for pbehavior.
     *
     * @param {Object} pbehavior
     */
    async fetchEventsForPbehavior(pbehavior) {
      const timespans = await this.fetchTimespansForPbehavior(pbehavior);

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
     * Remove pbehavior from events.
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
          this.$emit(
            'update:changedPbehaviorsById',
            omit(this.changedPbehaviorsById, [pbehavior._id]),
          );
        }
      }

      this.events = this.events.filter(
        event => get(event.data, 'pbehavior._id') !== pbehavior._id,
      );
    },

    /**
     * Apply changes for the pbehavior.
     *
     * @param {Object} event
     * @param {Function} close
     */
    async addEventWithClose(event, close) {
      const { pbehavior } = event;

      if (pbehavior) {
        await this.updatePbehavior(pbehavior);
      }

      close();
    },

    /**
     * Apply changes for selected event only.
     *
     * @param {Object} event
     * @returns {Promise<Array>}
     */
    applyEventChangesForSelectedHandler({ data, start, oldStart, end, oldEnd }) {
      const { pbehavior } = data;
      const originalPbehavior = this.allPbehaviorsById[pbehavior._id];
      const timezone = pbehavior.timezone || this.timezone;
      const tstart = convertDateToTimestampByTimezone(start, timezone);
      const tstop = convertDateToTimestampByTimezone(end, this.timezone);

      const exdate = {
        begin: convertDateToTimestampByTimezone(oldStart, timezone),
        end: convertDateToTimestampByTimezone(oldEnd, timezone),
        type: this.getOppositePbehaviorType(pbehavior.type),
      };

      const mainPbehavior = {
        ...originalPbehavior,

        exdates: [...(originalPbehavior.exdates || []), exdate],
      };

      const newPbehavior = {
        ...pbehavior,

        _id: uid('pbehavior'),
        rrule: '',
        tstart,
        tstop,
      };

      return Promise.all([this.updatePbehavior(mainPbehavior), this.updatePbehavior(newPbehavior)]);
    },

    /**
     * Apply event changes for all events.
     *
     * @param {Object} event
     * @returns {Promise<void>}
     */
    applyEventChangesForAllHandler({ data, start, oldStart, end, oldEnd }) {
      const { pbehavior } = data;
      const originalPbehavior = this.allPbehaviorsById[pbehavior._id];
      const tstart = originalPbehavior.tstart + getDiffBetweenDates(start, oldStart);
      const tstop = pbehavior.tstop
        ? originalPbehavior.tstop + getDiffBetweenDates(end, oldEnd)
        : null;

      const newPbehavior = {
        ...pbehavior,

        tstart,
        tstop,
      };

      return this.updatePbehavior(newPbehavior);
    },

    closePopoverForEvent() {
      this.$refs.calendar.clearPlaceholder();
    },

    /**
     * Show modal window for recurrent changes confirmation.
     *
     * @param {Object} event
     */
    showPbehaviorRecurrentChangesConfirmationModal(event) {
      this.modals.show({
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
              console.error(err);

              this.popups.error({
                text: err.description || err.message || this.$t('errors.default'),
              });
              this.closePopoverForEvent(event);
            }
          },
          cancel: () => this.closePopoverForEvent(event),
        },
      });
    },

    /**
     * Changed event handler.
     *
     * @param {Object} event
     */
    async handleUpdateEvent(event) {
      const { pbehavior } = event.data;

      if (!pbehavior.rrule) {
        const tstart = convertDateToTimestampByTimezone(event.start, this.timezone);
        const tstop = pbehavior.tstop
          ? convertDateToTimestampByTimezone(event.end, this.timezone)
          : null;

        await this.updatePbehavior({
          ...pbehavior,

          tstart,
          tstop,
        });

        this.closePopoverForEvent();

        return;
      }

      this.showPbehaviorRecurrentChangesConfirmationModal(event);
    },

    /**
     * Update pbehavior in the data and fetch timespans for that.
     *
     * @param {Object} pbehavior
     */
    async updatePbehavior(pbehavior) {
      this.pending = true;

      const hasPbehaviorInList =
        this.pbehaviorsById[pbehavior._id] || this.changedPbehaviorsById[pbehavior._id];

      if (hasPbehaviorInList) {
        this.$emit('update:changedPbehaviorsById', {
          ...this.changedPbehaviorsById,
          [pbehavior._id]: pbehavior,
        });
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
     * Fetch default pbehavior types.
     */
    async fetchDefaultTypes() {
      const { data } = await this.fetchPbehaviorTypesFieldList({
        params: { default: true },
      });

      this.defaultTypes = data;
    },

    /**
     * Get opposite pbehavior type for default exdate.
     *
     * @param {Object} [pbehaviorType = {}]
     * @returns {Object|undefined}
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

    handleUpdateTimezone(newTimezone) {
      this.$emit('update:timezone', newTimezone);
    },
  },
};
</script>
