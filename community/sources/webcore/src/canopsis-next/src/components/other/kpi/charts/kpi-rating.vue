<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-rating-filters(v-model="pagination", :min-date="minDate")
    kpi-rating-chart(
      :metrics="ratingMetrics",
      :metric="pagination.metric",
      :downloading="downloading",
      :min-date="minDate",
      responsive,
      @export:csv="exportRatingMetricsAsCsv",
      @export:png="exportRatingMetricsAsPng"
    )
    kpi-error-overlay(v-if="unavailable || fetchError")
</template>

<script>
import { isUndefined } from 'lodash';

import { API_HOST, API_ROUTES, KPI_RATING_METRICS_FILENAME_PREFIX } from '@/config';

import {
  QUICK_RANGES,
  ALARM_METRIC_PARAMETERS,
  DATETIME_FORMATS,
  USER_METRIC_PARAMETERS,
  SAMPLINGS,
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

import KpiRatingFilters from './partials/kpi-rating-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiRatingChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-rating-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiErrorOverlay, KpiRatingFilters, KpiRatingChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportMixinCreator({
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
      downloading: false,
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
    unavailable(unavailable) {
      if (!unavailable) {
        this.fetchList();
      }
    },
  },
  methods: {
    getFileName() {
      const fromTime = convertDateToString(this.interval.from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(this.interval.to, DATETIME_FORMATS.short);

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
      return !isUndefined(query.criteria) && isMetricsQueryChanged(query, oldQuery, this.minDate);
    },

    getQuery() {
      return {
        ...this.interval,

        criteria: this.query.criteria?.id,
        filter: this.query.metric !== USER_METRIC_PARAMETERS.totalUserActivity ? this.query.filter : undefined,
        metric: this.query.metric,
        limit: this.query.rowsPerPage,
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

    async exportRatingMetricsAsCsv() {
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
