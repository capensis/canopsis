<template>
  <div>
    <v-layout class="calender-wrapper">
      <c-progress-overlay :pending="pending" />
      <c-alert-overlay
        :value="hasError"
        :message="serverErrorMessage"
      />
      <c-calendar
        ref="calendar"
        :type.sync="type"
        :class="['stats-calendar-app', { single: !hasMultipleFilters }]"
        :events="events"
        readonly
        @change:pagination="changeCalendar"
        @click:event="eventClick"
      >
        <template #event="{ event }">
          {{ event.name }}
        </template>
      </c-calendar>
      <stats-calendar-menu
        v-if="hasMenu"
        :activator="menuActivator"
        :calendar-event="menuCalendarEvent"
        @click:event="menuEventClick"
        @closed="closedMenu"
      />
    </v-layout>
  </div>
</template>

<script>
import { get, isEmpty, omit } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { MODALS, MAX_LIMIT, CALENDAR_TYPES } from '@/constants';

import { convertDateToTimestamp, convertDateToTimestampByTimezone } from '@/helpers/date/date';
import { convertAlarmsToEvents, convertEventsToGroupedEvents } from '@/helpers/calendar/dayspan';
import { generatePreparedDefaultAlarmListWidget } from '@/helpers/entities/widget/form';

import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

import StatsCalendarMenu from './stats-calendar-menu.vue';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

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
      const groupByValue = this.type === CALENDAR_TYPES.month ? 'day' : 'hour';

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

      if (this.type === CALENDAR_TYPES.month) {
        return convertEventsToGroupedEvents({
          events,
          getColor: this.getCalendarEventColor,
        });
      }

      return events;
    },

    hasMultipleFilters() {
      return this.query?.filters?.length > 1;
    },

    hasFilters() {
      return this.query?.filters?.length > 0;
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

    menuEventClick(event) {
      const meta = get(event, 'data.meta', {});

      this.showAlarmsListModal(meta);
    },

    eventClick({ event, nativeEvent }) {
      const meta = get(event, 'data.meta', {});

      if (nativeEvent.target && meta.events) {
        this.menuActivator = nativeEvent.target;
        this.menuCalendarEvent = event;

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
          title: meta.filter.title
            ? this.$t('modals.alarmsList.prefixTitle', { prefix: meta.filter.title })
            : this.$t('modals.alarmsList.title'),
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
        const query = this.getCommonQuery();

        query.tstart = convertDateToTimestampByTimezone(this.$refs.calendar.filled.start, this.$system.timezone);
        query.tstop = convertDateToTimestampByTimezone(this.$refs.calendar.filled.end, this.$system.timezone);
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

<style lang="scss">
  .calender-wrapper {
    position: relative;

    .v-calendar-weekly__day {

      .v-event {
        position: absolute;
        top: 0;
        bottom: 0;
        right: 0;
        left: 0;
        width: 100% !important;
        height: 100% !important;
        display: flex;
        align-items: center;
        justify-content: center;
      }
    }

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
