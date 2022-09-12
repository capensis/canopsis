<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-sli-filters(v-model="pagination", :min-date="minDate")
    kpi-sli-chart(
      :metrics="sliMetrics",
      :data-type="pagination.type",
      :sampling="pagination.sampling",
      :downloading="downloading",
      :min-date="minDate",
      responsive,
      @export:csv="exportSliMetricsAsCsv",
      @export:png="exportSliMetricsAsPng"
    )
    kpi-error-overlay(v-if="unavailable || fetchError")
</template>

<script>
import { KPI_SLI_METRICS_FILENAME_PREFIX } from '@/config';
import {
  QUICK_RANGES,
  SAMPLINGS,
  KPI_SLI_GRAPH_DATA_TYPE,
  DATETIME_FORMATS,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { convertDateToStartOfDayTimestamp, convertDateToString } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';
import { isMetricsQueryChanged } from '@/helpers/metrics';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { exportCsvMixinCreator } from '@/mixins/widget/export';

import KpiSliFilters from './partials/kpi-sli-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiSliChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-sli-chart.vue');

export default {
  components: { KpiErrorOverlay, KpiSliFilters, KpiSliChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportCsvMixinCreator({
      createExport: 'createKpiSliExport',
      fetchExport: 'fetchMetricExport',
      fetchExportFile: 'fetchMetricCsvFile',
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
      downloading: false,
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
      const fromTime = convertDateToString(
        convertStartDateIntervalToTimestamp(this.query.interval.from),
        DATETIME_FORMATS.short,
      );
      const toTime = convertDateToString(
        convertStopDateIntervalToTimestamp(this.query.interval.to),
        DATETIME_FORMATS.short,
      );

      return [KPI_SLI_METRICS_FILENAME_PREFIX, fromTime, toTime, this.query.sampling].join('-');
    },

    async exportSliMetricsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    getQuery() {
      return {
        from: convertStartDateIntervalToTimestamp(this.query.interval.from),
        to: convertStopDateIntervalToTimestamp(this.query.interval.to),
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

        this.sliMetrics = sliMetrics;
        this.minDate = convertDateToStartOfDayTimestamp(minDate);

        if (params.from < this.minDate) {
          this.updateQueryField('interval', { ...this.query.interval, from: this.minDate });
        }
      } catch (err) {
        this.fetchError = true;
      } finally {
        this.pending = false;
      }
    },

    async exportSliMetricsAsCsv() {
      this.downloading = true;

      await this.exportAsCsv({
        name: this.getFileName(),
        data: this.getQuery(),
      });

      this.downloading = false;
    },
  },
};
</script>
