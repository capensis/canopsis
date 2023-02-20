<template lang="pug">
  div
    v-layout.calender-wrapper
      c-progress-overlay(:pending="pending")
      c-alert-overlay(
        :value="hasError",
        :message="serverErrorMessage"
      )
      ds-calendar-app.stats-calendar-app(
        :class="{ single: !hasMultipleFilters }",
        :calendar="calendar",
        :events="events",
        fluid,
        read-only,
        @change="changeCalendar",
        @edit="eventClick"
      )
      stats-calendar-menu(
        v-if="hasMenu",
        :activator="menuActivator",
        :calendarEvent="menuCalendarEvent",
        @event-click="menuEventClick",
        @closed="closedMenu"
      )
</template>

<script>
import { get, isEmpty, omit } from 'lodash';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Units } from 'dayspan';

import { MODALS, MAX_LIMIT } from '@/constants';

import { convertDateToTimestamp } from '@/helpers/date/date';
import { convertAlarmsToEvents, convertEventsToGroupedEvents } from '@/helpers/calendar/dayspan';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import StatsCalendarMenu from './stats-calendar-menu.vue';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

export default {
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
      menuActivator: null,
      menuCalendarEvent: null,
      pending: false,
      alarms: [],
      alarmsCollections: [],
      calendar: Calendar.months(),
      serverErrorMessage: null,
    };
  },
  computed: {
    hasError() {
      return !!this.serverErrorMessage;
    },

    hasMenu() {
      return this.menuActivator && this.menuCalendarEvent;
    },

    getCalendarEventColor() {
      return (count) => {
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
      };
    },

    events() {
      const groupByValue = this.calendar.type === Units.MONTH ? 'day' : 'hour';

      if (!this.hasFilters) {
        return convertAlarmsToEvents({
          groupByValue,
          alarms: this.alarms,
          getColor: this.getCalendarEventColor,
        });
      }

      const events = this.alarmsCollections.reduce((acc, alarms, index) => acc.concat(convertAlarmsToEvents({
        alarms,
        groupByValue,
        filter: this.query.filters[index],
        getColor: this.getCalendarEventColor,
      })), []);

      if (this.calendar.type !== Units.MONTH) {
        return convertEventsToGroupedEvents({
          events,
          getColor: this.getCalendarEventColor,
        });
      }

      return events;
    },

    hasMultipleFilters() {
      return get(this.query, 'filters.length', 0) > 1;
    },

    hasFilters() {
      return get(this.query, 'filters.length') > 0;
    },
  },
  methods: {
    ...alarmMapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    closedMenu() {
      this.menuActivator = null;
      this.menuCalendarEvent = null;
    },

    menuEventClick(calendarEvent) {
      const meta = get(calendarEvent, 'data.meta', {});

      this.showAlarmsListModal(meta);
    },

    eventClick({ $element, calendarEvent }) {
      const meta = get(calendarEvent, 'data.meta', {});

      if ($element && meta.events) {
        this.menuActivator = $element;
        this.menuCalendarEvent = calendarEvent;

        return;
      }

      this.showAlarmsListModal(meta);
    },

    getCommonQuery() {
      return omit(this.query, ['filters', 'considerPbehaviors']);
    },

    showAlarmsListModal(meta) {
      const widget = generatePreparedDefaultAlarmListWidget();

      widget.parameters = {
        ...widget.parameters,
        ...this.widget.parameters.alarmsList,
      };

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget,
          title: this.$t('modals.alarmsList.prefixTitle', {
            prefix: meta.filter.title,
          }),
          fetchList: (params) => {
            const newParams = {
              ...this.getCommonQuery(),
              ...params,

              tstart: convertDateToTimestamp(meta.tstart),
              tstop: convertDateToTimestamp(meta.tstop),
            };

            if (meta.filter?._id) {
              newParams.filters = [meta.filter._id];
            }

            return this.fetchAlarmsListWithoutStore({
              params: newParams,
            });
          },
        },
      });
    },

    changeCalendar() {
      this.fetchList();
    },

    async fetchList() {
      try {
        const { start, end } = this.calendar.filled;
        const query = this.getCommonQuery();

        query.tstart = start.date.unix();
        query.tstop = end.date.unix();
        query.limit = MAX_LIMIT;

        this.pending = true;
        this.serverErrorMessage = null;

        if (isEmpty(this.query.filters)) {
          let { data: alarms } = await this.fetchAlarmsListWithoutStore({
            params: query,
          });

          if (this.query.considerPbehaviors) {
            alarms = alarms.filter(alarm => isEmpty(alarm.pbehaviors));
          }

          this.alarms = alarms;
          this.alarmsCollections = [];
        } else {
          const results = await Promise.all(this.query.filters.map(({ _id: id }) => this.fetchAlarmsListWithoutStore({
            params: {
              ...query,

              filters: [id],
            },
          })));

          this.alarmsCollections = results.map(({ data: alarms }) => {
            if (this.query.considerPbehaviors) {
              return alarms.filter(alarm => isEmpty(alarm.pbehaviors));
            }

            return alarms;
          });

          this.alarms = [];
        }
      } catch (err) {
        this.serverErrorMessage = err.description || this.$t('errors.statsRequestProblem');
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .calender-wrapper {
    position: relative;

    & ::v-deep .ds-calendar-event {
      font-size: 14px;
    }

    & ::v-deep .ds-calendar-app.stats-calendar-app {
      .ds-calendar-event {
        cursor: pointer !important;
      }

      &.single {
        .ds-calendar-event-menu {
          position: absolute;
          left: 0;
          top: 0;
          width: 100%;
          height: 100% !important;
          padding: 4px;

          .v-menu__activator {
            width: 100%;
            height: 100%;
          }

          .ds-calendar-event {
            padding-left: 0;
            display: flex;
            height: 100%;
            width: 100%;

            & > span {
              margin: auto;
              text-align: center;
            }

            .ds-ev-description {
              display: none;
            }
          }
        }

        .ds-week {
          .ds-ev-description {
            display: none;
          }
        }
      }

      &:not(.single) {
        & ::v-deep .ds-calendar-event-menu {
          position: relative;
          height: 20px;

          .ds-calendar-event {
            top: 0 !important;
          }

          .ds-ev-title {
            margin-right: 10px;
          }
        }
      }
    }

    & ::v-deep .ds-week-view {
      .ds-ev-title {
        display: block;
      }
    }

    & ::v-deep .ds-day {
      position: relative;

      .ds-dom {
        border-radius: 12px;
        background-color: white;
        display: inline-block;
        position: relative;
        z-index: 1;

        &.ds-today-dom {
          background-color: #4285f4;
        }

        .theme--dark & {
          background-color: black;
        }
      }

      .ds-day-header {
        z-index: 10;
      }

      .ds-calendar-event > .v-menu__activator {
        height: 100%;
      }
    }
  }
</style>
