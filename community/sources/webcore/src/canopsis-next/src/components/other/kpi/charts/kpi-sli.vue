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
import { API_HOST, API_ROUTES, KPI_SLI_METRICS_FILENAME_PREFIX } from '@/config';
import {
  QUICK_RANGES,
  SAMPLINGS,
  KPI_SLI_GRAPH_DATA_TYPE,
  DATETIME_FORMATS,
} from '@/constants';

import { saveFile } from '@/helpers/file/files';
import {
  convertStartDateIntervalToTimestampByTimezone,
  convertStopDateIntervalToTimestampByTimezone,
} from '@/helpers/date/date-intervals';
import {
  convertDateToStartOfDayTimestampByTimezone,
  convertDateToString,
} from '@/helpers/date/date';
import { convertMetricsToTimezone, isMetricsQueryChanged } from '@/helpers/metrics';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { exportMixinCreator } from '@/mixins/widget/export';

import KpiSliFilters from './partials/kpi-sli-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiSliChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-sli-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiErrorOverlay, KpiSliFilters, KpiSliChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportMixinCreator({
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
  computed: {
    interval() {
      return {
        from: convertStartDateIntervalToTimestampByTimezone(
          this.query.interval.from,
          DATETIME_FORMATS.datePicker,
          SAMPLINGS.day,
          this.$system.timezone,
        ),
        to: convertStopDateIntervalToTimestampByTimezone(
          this.query.interval.to,
          DATETIME_FORMATS.datePicker,
          SAMPLINGS.day,
          this.$system.timezone,
        ),
      };
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
      const fromTime = convertDateToString(this.interval.from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(this.interval.to, DATETIME_FORMATS.short);

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
        ...this.interval,

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

    async exportSliMetricsAsCsv() {
      this.downloading = true;

      try {
        const fileData = await this.generateFile({
          data: this.getQuery(),
        });

        this.downloadFile(`${API_HOST}${API_ROUTES.metrics.exportMetric}/${fileData._id}/download`);
      } catch (err) {
        this.$popups.error({ text: err?.error ?? this.$t('errors.default') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>
