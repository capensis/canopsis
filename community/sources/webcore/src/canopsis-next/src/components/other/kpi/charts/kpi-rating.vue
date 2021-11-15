<template lang="pug">
  div
    v-layout.ml-4.mb-4(align-center)
      c-quick-date-interval-field(
        :interval="pagination.interval",
        @input="updateInterval"
      )
    div
      kpi-rating-chart(:metrics="ratingMetrics", :metric="pagination.metric", responsive)
</template>

<script>
import {
  QUICK_RANGES,
  ALARM_METRIC_PARAMETERS,
  KPI_RATING_CRITERIA,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';

import KpiRatingChart from './partials/kpi-rating-chart.vue';

export default {
  components: { KpiRatingChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      ratingMetrics: [],
      query: {
        criteria: KPI_RATING_CRITERIA.user,
        metric: ALARM_METRIC_PARAMETERS.ticketAlarms,
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
        criteria: this.pagination.criteria,
        metric: this.pagination.metric,
        limit: this.pagination.rowsPerPage,
      };
    },

    async fetchList() {
      this.ratingMetrics = await this.fetchRatingMetricsWithoutStore({
        params: this.getQuery(),
      });
    },
  },
};
</script>
