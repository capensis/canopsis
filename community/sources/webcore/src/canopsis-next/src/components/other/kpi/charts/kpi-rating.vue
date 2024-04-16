<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <kpi-rating-filters
      v-model="options"
      :min-date="minDate"
    />
    <kpi-rating-chart
      :metrics="ratingMetrics"
      :metric="query.metric"
      :downloading="downloading"
      :min-date="minDate"
      responsive
      @export:csv="exportMetricsAsCsv"
      @export:png="exportMetricsAsPng"
    />
    <kpi-error-overlay v-if="unavailable || fetchError" />
  </div>
</template>

<script>
import { isUndefined } from 'lodash';

import { KPI_RATING_METRICS_FILENAME_PREFIX } from '@/config';
import { QUICK_RANGES, ALARM_METRIC_PARAMETERS, DATETIME_FORMATS, USER_METRIC_PARAMETERS } from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone, convertDateToString } from '@/helpers/date/date';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query/query';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { metricsExportMixinCreator } from '@/mixins/widget/metrics/export';

import KpiRatingFilters from './partials/kpi-rating-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiRatingChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-rating-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiErrorOverlay, KpiRatingFilters, KpiRatingChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    queryIntervalFilterMixin,
    metricsExportMixinCreator({
      createExport: 'createKpiRatingExport',
      fetchExport: 'fetchMetricExport',
    }),
  ],
  props: {
    unavailable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      ratingMetrics: [],
      pending: false,
      fetchError: false,
      minDate: null,
      query: {
        criteria: undefined,
        filter: undefined,
        metric: ALARM_METRIC_PARAMETERS.ticketActiveAlarms,
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
    };
  },
  watch: {
    unavailable(unavailable) {
      if (!unavailable) {
        this.fetchList();
      }
    },
  },
  methods: {
    getFileName() {
      const { from, to } = this.getIntervalQuery();

      const fromTime = convertDateToString(from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(to, DATETIME_FORMATS.short);

      return [
        KPI_RATING_METRICS_FILENAME_PREFIX,
        fromTime,
        toTime,
        this.query.metric,
        this.query.criteria?.id,
      ].join('-');
    },

    customQueryCondition(query, oldQuery) {
      return !isUndefined(query.criteria) && isMetricsQueryChanged(query, oldQuery, this.minDate);
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),

        criteria: this.query.criteria?.id,
        filter: this.query.metric !== USER_METRIC_PARAMETERS.totalUserActivity ? this.query.filter : undefined,
        metric: this.query.metric,
        limit: this.query.itemsPerPage,
      };
    },

    async fetchList() {
      try {
        this.pending = true;

        const params = this.getQuery();

        const {
          data: ratingMetrics,
          meta: { min_date: minDate },
        } = await this.fetchRatingMetricsWithoutStore({ params });

        this.ratingMetrics = convertMetricsToTimezone(ratingMetrics, this.$system.timezone);
        this.minDate = convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone);
        this.fetchError = false;

        if (params.from < this.minDate) {
          this.updateQueryField('interval', { ...this.query.interval, from: this.minDate });
        }
      } catch (err) {
        this.fetchError = true;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
