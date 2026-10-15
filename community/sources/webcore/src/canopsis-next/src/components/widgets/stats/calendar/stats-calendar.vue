<template>
  <div class="stats-calendar pa-2 w-100">
    <c-layout class="stats-calendar__wrapper w-100">
      <c-alert-overlay
        :value="hasError"
        :message="serverErrorMessage"
      />

      <div class="stats-calendar__calendar w-100">
        <c-calendar
          v-model:type="type"
          ref="calendar"
          :class="['w-100', { 'stats-calendar__app--single-month': isSingleMonthView }]"
          :events="events"
          :loading="pending"
          readonly
          @change:pagination="changeCalendar"
          @click:event="eventClick"
        >
          <template #event="{ event }">
            <div
              :style="{ backgroundColor: getCalendarEventSource(event).color }"
              class="stats-calendar__event"
              role="button"
              tabindex="0"
              @click.stop="calendarEventClick(event, $event)"
              @keydown.enter.stop="calendarEventClick(event, $event)"
            >
              {{ getCalendarEventSource(event).name }}
            </div>
          </template>
        </c-calendar>
      </div>

      <stats-calendar-menu
        v-if="hasMenu"
        :attach="calendarElement"
        :activator="menuActivator"
        :calendar-event="menuCalendarEvent"
        @click:event="menuEventClick"
        @closed="closedMenu"
      />
    </c-layout>
  </div>
</template>

<script>
import { isEmpty, omit } from 'lodash';
import { mapActions } from 'pinia';

import {
  CALENDAR_TYPES,
  MAX_LIMIT,
  MODALS,
  PATTERN_CONDITIONS,
  PATTERNS_FIELDS,
  PBEHAVIOR_CANONICAL_TYPES,
  PBEHAVIOR_PATTERN_FIELDS,
} from '@/constants';

import { useAlarmStore } from '@/stores/modules/alarm';

import { convertAlarmsToEvents, convertEventsToGroupedEvents } from '@/helpers/calendar/calendar';
import {
  convertDateToTimestamp,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import StatsCalendarMenu from './stats-calendar-menu.vue';

const getCalendarEventSource = event => event?.input ?? event;
const getCalendarEventMeta = event => getCalendarEventSource(event)?.data?.meta ?? {};

export default {
  inject: ['$system'],

  components: {
    StatsCalendarMenu,
  },

  mixins: [widgetFetchQueryMixin],

  props: {
    widget: {
      type: Object,
      required: true,
    },
  },

  data() {
    return {
      type: CALENDAR_TYPES.month,
      menuActivator: null,
      menuCalendarEvent: null,
      pending: false,
      alarms: [],
      alarmsCollections: [],
      serverErrorMessage: null,
    };
  },

  computed: {
    hasError() {
      return !!this.serverErrorMessage;
    },

    calendarElement() {
      return this.$refs.calendar?.$el;
    },

    hasMenu() {
      return this.menuActivator && this.menuCalendarEvent;
    },

    hasMultipleFilters() {
      return this.query?.filters?.length > 1;
    },

    hasFilters() {
      return this.query?.filters?.length > 0;
    },

    isSingleMonthView() {
      return this.type === CALENDAR_TYPES.month && !this.hasMultipleFilters;
    },

    noPbehaviorPattern() {
      return JSON.stringify([[
        {
          field: PBEHAVIOR_PATTERN_FIELDS.canonicalType,
          cond: {
            type: PATTERN_CONDITIONS.equal,
            value: PBEHAVIOR_CANONICAL_TYPES.active,
          },
        },
      ]]);
    },

    events() {
      const groupByValue = this.type === CALENDAR_TYPES.month ? 'day' : 'hour';

      if (!this.hasFilters) {
        return convertAlarmsToEvents({
          groupByValue,
          alarms: this.alarms,
          getColor: this.getCalendarEventColor,
        });
      }

      const events = this.alarmsCollections.reduce((acc, collection, index) => acc.concat(
        convertAlarmsToEvents({
          alarms: collection,
          groupByValue,
          filter: this.query.filters[index],
          getColor: this.getCalendarEventColor,
        }),
      ), []);

      if (this.type === CALENDAR_TYPES.month) {
        return convertEventsToGroupedEvents({
          events,
          getColor: this.getCalendarEventColor,
        });
      }

      return events;
    },
  },

  methods: {
    ...mapActions(useAlarmStore, {
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    getCalendarEventSource,

    getCalendarEventMeta(event) {
      return getCalendarEventMeta(event);
    },

    getCalendarEventColor(count) {
      const { criticityLevels, criticityLevelsColors } = this.widget.parameters;

      if (count >= criticityLevels.critical) {
        return criticityLevelsColors.critical;
      }

      if (count >= criticityLevels.major) {
        return criticityLevelsColors.major;
      }

      if (count >= criticityLevels.minor) {
        return criticityLevelsColors.minor;
      }

      return criticityLevelsColors.ok;
    },

    getCommonQuery() {
      const query = omit(this.query, ['filters', 'considerPbehaviors']);

      if (this.query.considerPbehaviors) {
        query[PATTERNS_FIELDS.pbehavior] = this.noPbehaviorPattern;
      }

      return query;
    },

    closedMenu() {
      this.menuActivator = null;
      this.menuCalendarEvent = null;
    },

    showAlarmsListModal(meta) {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget._id = this.widget._id;
      widget.parameters = {
        ...widget.parameters,
        ...this.widget.parameters.alarmsList,
      };

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: meta.filter?.title
            ? this.$t('modals.alarmsList.prefixTitle', { prefix: meta.filter.title })
            : this.$t('modals.alarmsList.title'),
          fetchList: (params) => {
            const nextParams = {
              ...params,
              ...this.getCommonQuery(),

              tstart: convertDateToTimestamp(meta.tstart),
              tstop: convertDateToTimestamp(meta.tstop),
            };

            if (meta.filter?._id) {
              nextParams.filters = [meta.filter._id];
            }

            return this.fetchAlarmsListWithoutStore({ params: nextParams });
          },
        },
      });
    },

    menuEventClick(event) {
      this.showAlarmsListModal(this.getCalendarEventMeta(event));
    },

    calendarEventClick(event, nativeEvent) {
      this.eventClick({ event, nativeEvent });
    },

    eventClick({ event, nativeEvent }) {
      const meta = this.getCalendarEventMeta(event);

      if (nativeEvent.target && meta.events) {
        this.menuActivator = nativeEvent.target;
        this.menuCalendarEvent = event;

        return;
      }

      this.showAlarmsListModal(meta);
    },

    async fetchList() {
      const filled = this.$refs.calendar?.filled;

      if (!filled?.start || !filled?.end || isEmpty(this.query)) {
        return;
      }

      try {
        const commonQuery = {
          ...this.getCommonQuery(),
          tstart: convertDateToTimestampByTimezone(filled.start, this.$system.timezone),
          tstop: convertDateToTimestampByTimezone(filled.end, this.$system.timezone),
          limit: MAX_LIMIT,
        };

        this.pending = true;
        this.serverErrorMessage = null;

        if (isEmpty(this.query.filters)) {
          const { data = [] } = await this.fetchAlarmsListWithoutStore({ params: commonQuery });

          this.alarms = data;
          this.alarmsCollections = [];

          return;
        }

        const results = await Promise.all(this.query.filters.map(({ _id: id }) => this.fetchAlarmsListWithoutStore({
            params: {
              ...commonQuery,

              filters: [id],
            },
          })));

        this.alarmsCollections = results.map(({ data = [] }) => data);
        this.alarms = [];
      } catch (err) {
        this.serverErrorMessage = err.description || this.$t('errors.statsRequestProblem');
      } finally {
        this.pending = false;
      }
    },

    changeCalendar() {
      if (!isEmpty(this.query)) {
        this.fetchList();
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.stats-calendar__wrapper {
  position: relative;
}

.stats-calendar__calendar {
  min-width: 0;
}

.stats-calendar__event {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: white;
  font-size: 14px;
}

:deep(.stats-calendar__app--single-month .v-calendar-weekly__day .v-event) {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100% !important;
  height: 100% !important;
  font-size: 14px;
}

:deep(.stats-calendar__app--single-month .v-calendar-weekly__day-label) {
  position: relative;
  z-index: 2;
  pointer-events: auto;
}

:deep(.c-calendar) {
  width: 100%;
}

:deep(.c-calendar__toolbar) {
  display: grid !important;
  grid-template-columns: 1fr auto 1fr;
  align-items: center;
  width: 100%;
  z-index: 3;
}

:deep(.c-calendar__toolbar-column) {
  width: auto;
  max-width: none;
  flex: initial;
}

:deep(.c-calendar__pagination) {
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

:deep(.c-timezone-field) {
  max-width: 400px;
}

:deep(.c-calendar__menu-right) {
  justify-content: flex-end !important;
}
</style>
