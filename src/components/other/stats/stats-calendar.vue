<template lang="pug">
  v-container
    v-layout.white(wrap, justify-space-between, align-center)
      v-flex
        v-btn(icon, @click="showSettings")
          v-icon settings
    v-layout.white.calender-wrapper
      v-fade-transition
        v-layout.white.progress(v-show="pending", column)
          v-progress-circular(indeterminate, color="primary")
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
import get from 'lodash/get';
import omit from 'lodash/omit';
import pick from 'lodash/pick';
import isEmpty from 'lodash/isEmpty';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Units } from 'dayspan';

import { SIDE_BARS } from '@/constants';
import { convertAlarmsToEvents, convertEventsToGroupedEvents } from '@/helpers/dayspan';

import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';

import DsCalendar from './day-span/calendar.vue';

const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');

export default {
  components: { DsCalendar },
  mixins: [sideBarMixin, widgetQueryMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    rowId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
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

    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsCalendarSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },

    editEvent(event) {
      const routeData = this.$router.resolve({
        name: 'alarms',
        query: {
          ...pick(event.data.meta, ['tstart', 'tstop']),
          ...pick(this.query, ['opened', 'resolved']),

          filter: JSON.stringify(event.data.meta.filter),
        },
      });

      window.open(routeData.href, '_blank');
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
      const query = omit(this.query, ['filters', 'considerPbehaviors']);

      this.pending = true;

      if (isEmpty(this.query.filters)) {
        let { alarms } = await this.fetchAlarmsListWithoutStore({
          params: query,
        });

        if (this.query.considerPbehaviors) {
          alarms = alarms.filter(alarm => isEmpty(alarm.pbehaviors));
        }

        this.alarms = alarms;
        this.alarmsCollections = [];
      } else {
        const results = await Promise.all(this.query.filters.map(({ filter }) => this.fetchAlarmsListWithoutStore({
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

      this.pending = false;
    },
  },
};
</script>

<style lang="scss" scoped>
  .calender-wrapper {
    position: relative;

    .progress {
      position: absolute;
      top: 0;
      left: 0;
      bottom: 0;
      right: 0;
      opacity: .4;
      z-index: 10;

      & /deep/ .v-progress-circular {
        top: 50%;
        left: 50%;
        margin-top: -16px;
        margin-left: -16px;
      }
    }
    .single {
      & /deep/ .ds-calendar-event-menu {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        height: 100%!important;
        padding: 4px;

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
    font-size: 12px;
    cursor: pointer;
    border-radius: 2px;
  }
</style>
