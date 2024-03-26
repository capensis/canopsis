<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <kpi-sli-filters
      v-model="query"
      :min-date="minDate"
    />
    <kpi-sli-chart
      :metrics="sliMetrics"
      :data-type="query.type"
      :sampling="query.sampling"
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
import { KPI_SLI_METRICS_FILENAME_PREFIX } from '@/config';
import { QUICK_RANGES, SAMPLINGS, KPI_SLI_GRAPH_DATA_TYPE, DATETIME_FORMATS } from '@/constants';

import { convertMetricIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { convertDateToStartOfDayTimestampByTimezone, convertDateToString } from '@/helpers/date/date';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query/query';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { metricsExportMixinCreator } from '@/mixins/widget/metrics/export';

import KpiSliFilters from './partials/kpi-sli-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiSliChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-sli-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiErrorOverlay, KpiSliFilters, KpiSliChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    queryIntervalFilterMixin,
    metricsExportMixinCreator({
      createExport: 'createKpiSliExport',
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
      sliMetrics: [],
      pending: false,
      fetchError: false,
      minDate: null,
      query: {
        sampling: SAMPLINGS.day,
        type: KPI_SLI_GRAPH_DATA_TYPE.percent,
        filter: null,
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
    };
  },
  computed: {
    interval() {
      return convertMetricIntervalToTimestamp({
        interval: this.query.interval,
        timezone: this.$system.timezone,
      });
    },
  },
  watch: {
    unavailable: {
      immediate: true,
      handler(unavailable) {
        if (!unavailable) {
          this.fetchList();
        }
      },
    },
  },
  methods: {
    customQueryCondition(query, oldQuery) {
      return isMetricsQueryChanged(query, oldQuery, this.minDate);
    },

    getFileName() {
      const { from, to } = this.getIntervalQuery();

      const fromTime = convertDateToString(from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(to, DATETIME_FORMATS.short);

      return [KPI_SLI_METRICS_FILENAME_PREFIX, fromTime, toTime, this.query.sampling].join('-');
    },

    getQuery() {
      return {
        ...this.getIntervalQuery(),

        in_percents: this.query.type === KPI_SLI_GRAPH_DATA_TYPE.percent,
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    async fetchList() {
      try {
        this.pending = true;

        const params = this.getQuery();

        const {
          data: sliMetrics,
          meta: { min_date: minDate },
        } = await this.fetchSliMetricsWithoutStore({ params });

        this.sliMetrics = convertMetricsToTimezone(sliMetrics, this.$system.timezone);
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
