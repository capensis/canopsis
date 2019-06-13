<template lang="pug">
  div
    v-layout.white.calender-wrapper
      progress-overlay(:pending="pending")
      stats-alert-overlay(:value="hasError", :message="serverErrorMessage")
      ds-calendar(
      :class="{ multiple: hasMultipleFilters, single: !hasMultipleFilters }",
      :events="events",
      @change="changeCalendar",
      @edit="editEvent"
      )
        v-card(slot="eventPopover", slot-scope="{ calendarEvent, details }")
          v-card-text
            v-layout(
            v-for="(event, index) in calendarEvent.data.meta.events",
            :key="`popover-event-${index}`",
            row,
            wrap
            )
              v-flex(xs12)
                div.ds-calendar-event(
                :style="{ backgroundColor: getStyleColor(details, event) }",
                @click="editEvent(event)"
                )
                  strong {{ event.data.title }}
                  p {{ event.data.description }}
</template>

<script>
import { get, omit, pick, isEmpty } from 'lodash';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Units } from 'dayspan';

import { MODALS, WIDGET_TYPES, LIVE_REPORTING_INTERVALS } from '@/constants';

import { convertAlarmsToEvents, convertEventsToGroupedEvents } from '@/helpers/dayspan';
import { generateWidgetByType } from '@/helpers/entities';

import modalMixin from '@/mixins/modal';
import widgetQueryMixin from '@/mixins/widget/query';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import DsCalendar from './day-span/calendar.vue';
import StatsAlertOverlay from '../partials/stats-alert-overlay.vue';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

export default {
  components: { ProgressOverlay, DsCalendar, StatsAlertOverlay },
  mixins: [modalMixin, widgetQueryMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      hasError: false,
      serverErrorMessage: null,
      alarms: [],
      alarmsCollections: [],
      calendar: Calendar.months(),
    };
  },
  computed: {
    getStyleColor() {
      return (details, calendarEvent) => {
        const past = calendarEvent.schedule.end.isBefore(new Date());

        return this.$dayspan.getStyleColor(details, calendarEvent, past);
      };
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

    editEvent(event) {
      const { meta } = event.data;
      const widget = generateWidgetByType(WIDGET_TYPES.alarmList);
      const widgetParameters = {
        ...this.widget.parameters.alarmsList,

        alarmsStateFilter: this.widget.parameters.alarmsStateFilter,
      };

      const query = { ...pick(meta, ['tstart', 'tstop']) };

      if (!isEmpty(event.data.meta.filter)) {
        widgetParameters.viewFilters = [meta.filter];
        widgetParameters.mainFilter = meta.filter;
      }

      if (query.tstart || query.tstop) {
        query.interval = LIVE_REPORTING_INTERVALS.custom;
      }

      this.showModal({
        name: MODALS.alarmsList,
        config: {
          query,
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

    changeCalendar({ calendar }) {
      this.calendar = calendar;
      this.query = {
        ...this.query,
        tstart: calendar.start.date.unix(),
        tstop: calendar.end.date.unix(),
      };
    },

    async fetchList() {
      try {
        const query = omit(this.query, ['filters', 'considerPbehaviors']);

        this.pending = true;
        this.hasError = false;
        this.serverErrorMessage = null;

        if (isEmpty(this.query.filters)) {
          let { alarms } = await this.fetchAlarmsListWithoutStore({
            withoutCatch: true,
            params: query,
          });

          if (this.query.considerPbehaviors) {
            alarms = alarms.filter(alarm => isEmpty(alarm.pbehaviors));
          }

          this.alarms = alarms;
          this.alarmsCollections = [];
        } else {
          const results = await Promise.all(this.query.filters.map(({ filter }) => this.fetchAlarmsListWithoutStore({
            withoutCatch: true,
            params: {
              ...query,
              filter,
            },
          })));


          this.alarmsCollections = results.map(({ alarms }) => {
            if (this.query.considerPbehaviors) {
              return alarms.filter(alarm => isEmpty(alarm.pbehaviors));
            }

            return alarms;
          });

          this.alarms = [];
        }
      } catch (err) {
        this.hasError = true;
        this.serverErrorMessage = err.description || null;
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

    .single {
      & /deep/ .ds-calendar-event-menu {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%!important;
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

    .multiple {
      & /deep/ .ds-calendar-event-menu {
        position: relative;
        height: 20px;

        .ds-ev-title {
          margin-right: 10px;
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
