<template lang="pug">
  div
    v-layout.ml-4.mb-4(align-center)
      c-quick-date-interval-field(
        :interval="pagination.interval",
        @input="updateInterval"
      )
    div
      kpi-alarms-chart(:metrics="alarmsMetrics", :sampling="pagination.sampling", responsive)
</template>

<script>
import {
  QUICK_RANGES,
  ALARM_METRIC_PARAMETERS,
  SAMPLINGS,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';

const KpiAlarmsChart = () => import(/* webpackChunkName: "Charts" */'./partials/kpi-alarms-chart.vue');

export default {
  components: { KpiAlarmsChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      alarmsMetrics: [],
      query: {
        sampling: SAMPLINGS.day,
        parameters: [ALARM_METRIC_PARAMETERS.totalAlarms],
        interval: {
          from: QUICK_RANGES.last30Days.start,
          to: QUICK_RANGES.last30Days.stop,
        },
      },
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    updateInterval(interval) {
      this.updateQueryField('interval', interval);
    },

    getQuery() {
      return {
        from: convertStartDateIntervalToTimestamp(this.pagination.interval.from),
        to: convertStopDateIntervalToTimestamp(this.pagination.interval.to),
        parameters: this.pagination.parameters,
        sampling: this.pagination.sampling,
      };
    },

    async fetchList() {
      this.alarmsMetrics = await this.fetchAlarmsMetricsWithoutStore({
        params: this.getQuery(),
      });
    },
  },
};
</script>
