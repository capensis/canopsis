<template lang="pug">
  div
    v-layout.white.calender-wrapper
      progress-overlay(:pending="pending")
      ds-calendar.single(
      ref="calendar",
      :events="events",
      :calendar="calendar",
      @change="changeCalendar",
      @edit="viewDay"
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
    viewDay(event) {
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

      if (this.query.mfilter) {
        const [watcherData = { values: [] }] = await this.fetchSpecialStatsListsWithoutStore({
          params: this.query,
        });

        this.values = watcherData.values;
      } else {
        this.values = [];
      }

      this.pending = false;
    },
  },
};
</script>
