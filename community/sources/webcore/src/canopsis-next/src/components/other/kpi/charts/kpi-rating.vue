<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-rating-filters(v-model="pagination")
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

import KpiRatingFilters from './partials/kpi-rating-filters.vue';

const KpiRatingChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-rating-chart.vue');

export default {
  components: { KpiRatingFilters, KpiRatingChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      ratingMetrics: [],
      pending: false,
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
      this.pending = true;

      this.ratingMetrics = await this.fetchRatingMetricsWithoutStore({
        params: this.getQuery(),
      });

      this.pending = false;
    },
  },
};
</script>
