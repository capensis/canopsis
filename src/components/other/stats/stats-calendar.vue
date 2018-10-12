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
      ds-calendar(:events="events", @change="changeCalendar")
        v-card(slot="eventPopover", slot-scope="{ calendarEvent }")
          v-card-text
            v-layout(
            v-for="(event, index) in calendarEvent.data.meta.events",
            :key="`popover-event-${index}`",
            row,
            wrap
            )
              v-flex(xs12)
                strong {{ event.data.title }}
                p {{ event.data.description }}
</template>

<script>
import moment from 'moment';
import omit from 'lodash/omit';
import isEmpty from 'lodash/isEmpty';
import groupBy from 'lodash/groupBy';
import { createNamespacedHelpers } from 'vuex';
import { Calendar, Day, Schedule, Units } from 'dayspan';

import { SIDE_BARS, STATS_CALENDAR_COLORS } from '@/constants';

import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';

import DsCalendar from './day-span/calendar.vue';

const { mapActions: entityMapActions } = createNamespacedHelpers('entity');
const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');
const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');

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
      events: [],
      calendar: Calendar.months(),
    };
  },
  methods: {
    ...entityMapActions({
      fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),

    ...alarmMapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    ...pbehaviorMapActions({
      fetchPbehaviorsListByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
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

    changeCalendar({ calendar }) {
      this.calendar = calendar;
      this.query = {
        ...this.query,
        tstart: calendar.start.date.unix(),
        tstop: calendar.end.date.unix(),
      };
    },

    async fetchList() {
      const query = omit(this.query, ['filters']);
      let results = [];

      this.pending = true;

      if (isEmpty(this.query.filters)) {
        results = await this.fetchAlarmsListWithoutStore({
          params: query,
        });

        this.events = this.prepareAlarms(results.alarms);
      } else {
        results = await Promise.all(this.query.filters.map(({ filter }) => this.fetchAlarmsListWithoutStore({
          params: {
            ...query,
            filter,
          },
        })));


        let events = results.reduce((acc, result, index) =>
          acc.concat(this.prepareAlarms(result.alarms, this.query.filters[index].title)), []);

        if (this.calendar.type === Units.WEEK) {
          const groupedEvents = groupBy(events, event => event.schedule.start.date.clone().startOf('hour').format());

          events = Object.keys(groupedEvents).map((dateString) => {
            const groupedEvent = groupedEvents[dateString];

            if (groupedEvent.length > 1) {
              const sum = groupedEvent.reduce((acc, event) => acc + event.data.meta.sum, 0);

              return {
                ...groupedEvent[0],

                data: {
                  title: sum,
                  color: this.getColor(sum),
                  meta: {
                    type: 'single',
                    hasPopover: true,
                    events: groupedEvent,
                  },
                },
              };
            }

            return groupedEvent[0];
          });
        }

        this.events = events;
      }

      this.pending = false;
    },

    getColor(count) {
      if (count > 50) {
        return STATS_CALENDAR_COLORS.alarm.large;
      }

      if (count > 30) {
        return STATS_CALENDAR_COLORS.alarm.medium;
      }

      return STATS_CALENDAR_COLORS.alarm.small;
    },

    prepareAlarms(alarms, prefix) {
      const startOfBy = this.calendar.type === Units.MONTH ? 'day' : 'hour';

      const groupedAlarms = groupBy(alarms, alarm => moment.unix(alarm.t).startOf(startOfBy).format());

      return Object.keys(groupedAlarms).map((dateString) => {
        const startDay = new Day(moment(dateString));
        const sum = groupedAlarms[dateString].length;

        return {
          data: {
            title: sum,
            description: prefix,
            color: this.getColor(sum),
            meta: {
              sum,
              type: prefix ? 'multiple' : 'single',
            },
          },
          schedule: new Schedule({
            on: startDay,
            times: [startDay.asTime()],
            duration: 1,
            durationUnit: 'hours',
          }),
        };
      });
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
  }
</style>
