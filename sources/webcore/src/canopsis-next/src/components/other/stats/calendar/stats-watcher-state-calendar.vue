<template lang="pug">
  div
    v-layout.white.calender-wrapper
      progress-overlay(:pending="pending")
      ds-calendar.single(
      :events="events",
      :calendar="calendar",
      @change="changeCalendar",
      @edit="editEvent"
      )
</template>

<script>
import moment from 'moment';
import { Calendar, Schedule, Day } from 'dayspan';

import { WATCHER_STATES_COLORS } from '@/constants';

import widgetQueryMixin from '@/mixins/widget/query';
import entitiesStatsMixin from '@/mixins/entities/stats';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import DsCalendar from './day-span/calendar.vue';

export default {
  components: { ProgressOverlay, DsCalendar },
  mixins: [widgetQueryMixin, entitiesStatsMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      values: [],
      events: [],
      alarmsCollections: [],
      calendar: Calendar.months(),
    };
  },
  methods: {
    editEvent() {},

    changeCalendar({ calendar }) {
      this.calendar = calendar;
      this.query = {
        ...this.query,
        tstart: calendar.start.date.unix(),
        tstop: calendar.end.date.unix(),
      };
    },

    async fetchList() {
      this.pending = true;

      const result = await this.fetchSpecialStatsListsWithoutStore({
        params: {
          mfilter: {
            type: 'resource',
            _id: 'feeder2_9/feeder2_11',
          },
          tstop: 1556668800,
          duration: '10m',
        },
      });

      this.events = result.reduce((acc, { values }) => {
        const events = values.map(({ duration, start, state }) => {
          const dateObject = moment.unix(start);
          const startDay = new Day(dateObject);

          return {
            data: {
              color: WATCHER_STATES_COLORS[state],
            },
            schedule: new Schedule({
              duration,
              on: startDay,
              times: [startDay.asTime()],
              durationUnit: 'seconds',
            }),
          };
        });

        return acc.concat(events);
      }, []);

      this.pending = false;
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

      & /deep/ .ds-month .ds-past {
        background: #eee;
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
