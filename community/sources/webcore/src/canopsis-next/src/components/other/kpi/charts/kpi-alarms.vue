<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-alarms-filters(v-model="pagination")
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

import KpiAlarmsFilters from './partials/kpi-alarms-filters.vue';

const KpiAlarmsChart = () => import(/* webpackChunkName: "Charts" */'./partials/kpi-alarms-chart.vue');

export default {
  components: { KpiAlarmsFilters, KpiAlarmsChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      alarmsMetrics: [],
      pending: false,
      query: {
        sampling: SAMPLINGS.day,
        parameters: [ALARM_METRIC_PARAMETERS.totalAlarms],
        filter: null,
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
    getQuery() {
      return {
        from: convertStartDateIntervalToTimestamp(this.query.interval.from),
        to: convertStopDateIntervalToTimestamp(this.query.interval.to),
        parameters: this.query.parameters,
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    async fetchList() {
      this.pending = true;

      this.alarmsMetrics = await this.fetchAlarmsMetricsWithoutStore({
        params: this.getQuery(),
      });

      this.pending = false;
    },
  },
};
</script>
