<template lang="pug">
  div
    v-layout.white.calender-wrapper
      progress-overlay(:pending="pending")
      alert-overlay(
        :value="hasError",
        :message="serverErrorMessage"
      )
      ds-calendar.stats-calendar-app(
        :calendar="calendar",
        :class="{ multiple: hasMultipleFilters }",
        :events="events",
        @change="changeCalendar",
        @edit="editEvent"
      )
        v-card(slot="eventPopover", slot-scope="{ calendarEvent, details }")
          v-card-text(v-if="calendarEvent.data.events")
            v-layout(
              v-for="(event, index) in calendarEvent.data.events",
              :key="`popover-event-${index}`",
              row,
              wrap
            )
              v-flex(xs12)
                div.ds-calendar-event(
                  :style="{ backgroundColor: getCalendarEventColor(event.data.total) }",
                  @click="editEvent(event)"
                )
                  strong {{ event.data.title }}
                  p {{ event.data.description }}
</template>

<script>
import { get, isEmpty, omit } from 'lodash';
import moment from 'moment-timezone';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Units } from 'dayspan';

import { DATETIME_FORMATS, MODALS, WIDGET_TYPES, COUNTER_GROUPING_TYPES } from '@/constants';

import { generateWidgetByType } from '@/helpers/entities';
import { convertCounterGroupToEvent, convertEventsToGroupedEvents } from '@/helpers/dayspan';

import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import AlertOverlay from '@/components/layout/alert/alert-overlay.vue';

import DsCalendar from './day-span/calendar.vue';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');
const { mapActions: counterMapActions } = createNamespacedHelpers('counter');

export default {
  components: { ProgressOverlay, DsCalendar, AlertOverlay },
  mixins: [widgetFetchQueryMixin, widgetStatsWrapperMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      counters: [],
      calendar: Calendar.months(),
    };
  },
  computed: {
    grouping() {
      return this.calendar.type === Units.MONTH ? COUNTER_GROUPING_TYPES.day : COUNTER_GROUPING_TYPES.hour;
    },

    events() {
      const events = this.counters.reduce((acc, counter, index) => {
        const filter = this.query.filters[index] || {};
        const filterEvents = Object.entries(counter.group)
          .map(([timestamp, counterGroup]) => convertCounterGroupToEvent({
            timestamp,
            counterGroup,
            filter,
            grouping: this.grouping,
            getColor: this.getCalendarEventColor,
          }));

        acc.push(...filterEvents);

        return acc;
      }, []);

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
  },
  methods: {
    ...alarmMapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    ...counterMapActions({
      fetchCountersListWithoutStore: 'fetchListWithoutStore',
    }),

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

    editEvent(event) {
      const { data, schedule } = event;
      const { filter } = data;
      const { start, end } = schedule;
      const widget = generateWidgetByType(WIDGET_TYPES.alarmList);
      const widgetParameters = {
        ...this.widget.parameters.alarmsList,

        alarmsStateFilter: this.widget.parameters.alarmsStateFilter,
        liveReporting: {
          tstart: start.date.format(DATETIME_FORMATS.dateTimePicker),
          tstop: end.date.clone()
            .subtract(1, 'seconds')
            .format(DATETIME_FORMATS.dateTimePicker),
        },
      };

      if (!isEmpty(filter)) {
        widgetParameters.viewFilters = [filter];
        widgetParameters.mainFilter = filter;
      }

      this.$modals.show({
        name: MODALS.alarmsList,
        config: {
          widget: {
            ...widget,

            parameters: {
              ...widget.parameters,
              ...widgetParameters,
            },
          },
        },
      });
    },

    changeCalendar() {
      this.fetchList();
    },

    getQuery() {
      const query = omit(this.query, ['filters']);

      query.tstart = this.calendar.filled.start.date.unix();
      query.tstop = this.calendar.filled.end.date.unix();
      query.grouping = this.grouping;
      query.local_timezone = moment.tz.guess();

      return query;
    },

    async fetchList() {
      try {
        const query = this.getQuery();

        this.pending = true;
        this.serverErrorMessage = null;

        if (isEmpty(this.query.filters)) {
          const { data: [counter] } = await this.fetchCountersListWithoutStore({
            params: query,
          });

          this.counters = [counter];
        } else {
          const results = await Promise.all(this.query.filters.map(({ filter }) => this.fetchCountersListWithoutStore({
            params: {
              ...query,
              filter,
            },
          })));

          this.counters = results.map(({ data: [counter] }) => counter);
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

    & /deep/ .ds-calendar-event {
      font-size: 14px;
    }

    .stats-calendar-app {
      &:not(.multiple) {
        & /deep/ .ds-calendar-event-menu {
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

        & /deep/ .ds-week {
          .ds-ev-description {
            display: none;
          }
        }
      }

      &.multiple {
        & /deep/ .ds-calendar-event-menu {
          position: relative;
          height: 20px;

          .ds-ev-title {
            margin-right: 10px;
          }
        }
      }
    }

    & /deep/ .ds-week-view {
      .ds-ev-title {
        display: block;
      }
    }

    & /deep/ .ds-day {
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
      }

      .ds-day-header {
        z-index: 10;
      }
    }
  }

  .ds-calendar-event {
    color: white;
    margin: 1px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    padding-left: 0.5em;
    cursor: pointer;
    border-radius: 2px;
  }
</style>
