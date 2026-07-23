<template>
  <div class="position-relative">
    <c-progress-overlay :pending="pending" />
    <kpi-rating-filters :query="query" :min-date="minDate" @update:query="query = $event" />
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

import { EXPORT_FETCHING_INTERVAL, KPI_RATING_METRICS_FILENAME_PREFIX } from '@/config';
import {
  ALARM_METRIC_PARAMETERS,
  DATETIME_FORMATS,
  EXPORT_STATUSES,
  QUICK_RANGES,
  USER_METRIC_PARAMETERS,
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
import { localQueryMixin } from '@/mixins/query/query';
import { queryIntervalFilterMixin } from '@/mixins/query/interval';

import KpiErrorOverlay from './partials/kpi-error-overlay.vue';
import KpiRatingChart from './partials/kpi-rating-chart.vue';
import KpiRatingFilters from './partials/kpi-rating-filters.vue';

export default {
  inject: ['$system'],

  components: {
    KpiErrorOverlay,
    KpiRatingChart,
    KpiRatingFilters,
  },

  mixins: [localQueryMixin, queryIntervalFilterMixin],

  props: {
    unavailable: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { createKpiRatingExport, fetchMetricExport, fetchRatingMetricsWithoutStore } =
      useMetrics();

    const { poll: generateFile } = useFilePolling({
      createHandler: async (...args) => {
        const { _id: id } = await createKpiRatingExport(...args);

        return { id };
      },
      fetchHandler: fetchMetricExport,
      completedStatus: EXPORT_STATUSES.completed,
      failedStatus: EXPORT_STATUSES.failed,
      interval: EXPORT_FETCHING_INTERVAL,
    });

    return {
      fetchRatingMetricsWithoutStore,
      generateFile,
    };
  },

  data() {
    return {
      ratingMetrics: [],
      pending: false,
      fetchError: false,
      downloading: false,
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
    getFileName() {
      const { from, to } = this.getPreparedIntervalQuery();

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
        ...this.getPreparedIntervalQuery(),
        criteria: this.query.criteria?.id,
        filter:
          this.query.metric !== USER_METRIC_PARAMETERS.totalUserActivity
            ? this.query.filter
            : undefined,
        metric: this.query.metric,
        limit: this.query.itemsPerPage,
      };
    },

    getPreparedIntervalQuery() {
      return convertMetricsIntervalToTimestamp(this.query.interval, this.$system.timezone);
    },

    async fetchList() {
      if (isUndefined(this.query.criteria)) {
        return;
      }

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
