<template lang="pug">
  div
    v-layout.white.calender-wrapper
      progress-overlay(:pending="pending")
      ds-calendar.single(
      ref="calendar",
      :events="events",
      :calendar="calendar",
      @change="changeCalendar",
      @edit="editEvent"
      )
</template>

<script>
import { Calendar, Units } from 'dayspan';

import { STATS_DURATION_UNITS } from '@/constants';

import widgetQueryMixin from '@/mixins/widget/query';
import entitiesStatsMixin from '@/mixins/entities/stats';

import { convertWatcherValuesToEvents } from '@/helpers/dayspan';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

import DsCalendar from './day-span/calendar.vue';

const CALENDAR_UNIT_TO_MOMENT_UNIT_MAP = {
  [Units.DAY]: 'day',
  [Units.WEEK]: 'week',
  [Units.MONTH]: 'month',
};

const CALENDAR_UNIT_TO_STATS_DURATION_UNIT_MAP = {
  [Units.DAY]: STATS_DURATION_UNITS.day,
  [Units.WEEK]: STATS_DURATION_UNITS.week,
  [Units.MONTH]: STATS_DURATION_UNITS.month,
};

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
      pending: false,
      values: [],
      calendar: Calendar.months(),
    };
  },
  computed: {
    events() {
      const groupByValue = this.calendar.type === Units.MONTH ? 'day' : 'hour';

      return convertWatcherValuesToEvents({
        groupByValue,
        values: this.values,
      });
    },
  },
  methods: {
    editEvent(event) {
      this.$refs.calendar.viewDay(event.day);
    },

    changeCalendar({ calendar }) {
      const momentUnit = CALENDAR_UNIT_TO_MOMENT_UNIT_MAP[calendar.type] || 'month';
      const durationUnit = CALENDAR_UNIT_TO_STATS_DURATION_UNIT_MAP[calendar.type] || STATS_DURATION_UNITS.month;

      this.calendar = calendar;
      this.query = {
        ...this.query,

        duration: `2${durationUnit}`,
        tstop: calendar.end.date.clone()
          .add(1, momentUnit)
          .utc()
          .startOf(momentUnit)
          .unix(),
      };
    },

    async fetchList() {
      this.pending = true;

      const [watcherData = { values: [] }] = await this.fetchSpecialStatsListsWithoutStore({
        params: {
          ...this.query,

          mfilter: {
            type: 'resource',
            _id: 'feeder2_9/feeder2_11',
          },
        },
      });

      this.values = watcherData.values;
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
