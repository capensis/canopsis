<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-rating-filters(v-model="pagination")
    kpi-rating-chart(
      :metrics="ratingMetrics",
      :metric="pagination.metric",
      :downloading="downloading",
      responsive,
      @export:csv="exportRatingMetricsAsCsv",
      @export:png="exportRatingMetricsAsPng"
    )
</template>

<script>
import { isUndefined, isEqual } from 'lodash';

import { KPI_RATING_METRICS_FILENAME_PREFIX } from '@/config';

import {
  QUICK_RANGES,
  ALARM_METRIC_PARAMETERS,
  DATETIME_FORMATS, USER_METRIC_PARAMETERS,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { convertDateToString } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { exportCsvMixinCreator } from '@/mixins/widget/export';

import KpiRatingFilters from './partials/kpi-rating-filters.vue';

const KpiRatingChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-rating-chart.vue');

export default {
  components: { KpiRatingFilters, KpiRatingChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportCsvMixinCreator({
      createExport: 'createKpiRatingExport',
      fetchExport: 'fetchMetricExport',
      fetchExportFile: 'fetchMetricCsvFile',
    }),
  ],
  data() {
    return {
      ratingMetrics: [],
      pending: false,
      downloading: false,
      query: {
        criteria: undefined,
        filter: undefined,
        metric: ALARM_METRIC_PARAMETERS.ticketAlarms,
        interval: {
          from: QUICK_RANGES.last30Days.start,
          to: QUICK_RANGES.last30Days.stop,
        },
      },
    };
  },
  mounted() {
    if (!isUndefined(this.query.criteria)) {
      this.fetchList();
    }
  },
  methods: {
    getFileName() {
      const fromTime = convertDateToString(
        convertStartDateIntervalToTimestamp(this.query.interval.from),
        DATETIME_FORMATS.short,
      );
      const toTime = convertDateToString(
        convertStopDateIntervalToTimestamp(this.query.interval.to),
        DATETIME_FORMATS.short,
      );

      return [
        KPI_RATING_METRICS_FILENAME_PREFIX,
        fromTime,
        toTime,
        this.query.metric,
        this.query.criteria?.id,
      ].join('-');
    },

    async exportRatingMetricsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    customQueryCondition(query, oldQuery) {
      return !isUndefined(query.criteria) && !isEqual(query, oldQuery);
    },

    getQuery() {
      return {
        from: convertStartDateIntervalToTimestamp(this.query.interval.from),
        to: convertStopDateIntervalToTimestamp(this.query.interval.to),
        criteria: this.query.criteria.id,
        filter: this.query.metric !== USER_METRIC_PARAMETERS.totalUserActivity ? this.query.filter : undefined,
        metric: this.query.metric,
        limit: this.query.rowsPerPage,
      };
    },

    async fetchList() {
      this.pending = true;

      this.ratingMetrics = await this.fetchRatingMetricsWithoutStore({
        params: this.getQuery(),
      });

      this.pending = false;
    },

    async exportRatingMetricsAsCsv() {
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
