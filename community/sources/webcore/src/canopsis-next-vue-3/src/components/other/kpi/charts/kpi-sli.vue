<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <kpi-sli-filters :query="query" :min-date="minDate" @update:query="query = $event" />
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
import { EXPORT_FETCHING_INTERVAL, KPI_SLI_METRICS_FILENAME_PREFIX } from '@/config';
import {
  DATETIME_FORMATS,
  EXPORT_STATUSES,
  KPI_SLI_GRAPH_DATA_TYPE,
  QUICK_RANGES,
  SAMPLINGS,
} from '@/constants';
import {
  convertDateToStartOfDayTimestampByTimezone,
  convertDateToString,
} from '@/helpers/date/date';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';
import {
  convertMetricsIntervalToTimestamp,
  isMetricsQueryChanged,
} from '@/helpers/entities/metric/query';
import { getExportMetricDownloadFileUrl } from '@/helpers/entities/metric/url';
import { saveFile } from '@/helpers/file/files';
import { openUrlInNewTab } from '@/helpers/url';
import { useFilePolling } from '@/hooks/polling';
import { useMetrics } from '@/hooks/store/modules/metrics';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';
import { localQueryMixin } from '@/mixins/query/query';

import KpiErrorOverlay from './partials/kpi-error-overlay.vue';
import KpiSliChart from './partials/kpi-sli-chart.vue';
import KpiSliFilters from './partials/kpi-sli-filters.vue';

export default {
  inject: ['$system'],

  components: {
    KpiErrorOverlay,
    KpiSliChart,
    KpiSliFilters,
  },

  mixins: [localQueryMixin, queryIntervalFilterMixin],

  props: {
    unavailable: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { createKpiSliExport, fetchMetricExport, fetchSliMetricsWithoutStore } = useMetrics();

    const { poll: generateFile } = useFilePolling({
      createHandler: async (...args) => {
        const { _id: id } = await createKpiSliExport(...args);

        return { id };
      },
      fetchHandler: fetchMetricExport,
      completedStatus: EXPORT_STATUSES.completed,
      failedStatus: EXPORT_STATUSES.failed,
      interval: EXPORT_FETCHING_INTERVAL,
    });

    return {
      fetchSliMetricsWithoutStore,
      generateFile,
    };
  },

  data() {
    return {
      sliMetrics: [],
      pending: false,
      fetchError: false,
      downloading: false,
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
      const { from, to } = this.getPreparedIntervalQuery();

      const fromTime = convertDateToString(from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(to, DATETIME_FORMATS.short);

      return [KPI_SLI_METRICS_FILENAME_PREFIX, fromTime, toTime, this.query.sampling].join('-');
    },

    getQuery() {
      return {
        ...this.getPreparedIntervalQuery(),
        in_percents: this.query.type === KPI_SLI_GRAPH_DATA_TYPE.percent,
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    getPreparedIntervalQuery() {
      return convertMetricsIntervalToTimestamp(this.query.interval, this.$system.timezone);
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
        console.error(err);

        this.fetchError = true;
      } finally {
        this.pending = false;
      }
    },

    async exportMetricsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        console.error(err);
        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    async exportMetricsAsCsv() {
      try {
        this.downloading = true;

        const fileData = await this.generateFile({
          data: this.getQuery(),
        });

        openUrlInNewTab(getExportMetricDownloadFileUrl(fileData._id ?? fileData.id));
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: this.$t('kpi.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>
