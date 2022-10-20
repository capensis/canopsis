<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-alarms-filters(v-model="pagination", :min-date="minDate")
    kpi-alarms-chart(
      :metrics="alarmsMetrics",
      :sampling="pagination.sampling",
      :downloading="downloading",
      :min-date="minDate",
      :interval="interval",
      responsive,
      @export:csv="exportAlarmMetricsAsCsv",
      @export:png="exportAlarmMetricsAsPng",
      @zoom="updateQueryField('interval', $event)"
    )
    kpi-error-overlay(v-if="unavailable || fetchError")
</template>

<script>
import { API_HOST, API_ROUTES, KPI_ALARM_METRICS_FILENAME_PREFIX } from '@/config';
import {
  QUICK_RANGES,
  ALARM_METRIC_PARAMETERS,
  SAMPLINGS,
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

import KpiAlarmsFilters from './partials/kpi-alarms-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiAlarmsChart = () => import(/* webpackChunkName: "Charts" */'./partials/kpi-alarms-chart.vue');

export default {
  inject: ['$system'],
  components: { KpiErrorOverlay, KpiAlarmsFilters, KpiAlarmsChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportMixinCreator({
      createExport: 'createKpiAlarmExport',
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
      alarmsMetrics: [],
      downloading: false,
      pending: false,
      fetchError: false,
      minDate: null,
      query: {
        sampling: SAMPLINGS.day,
        parameters: [ALARM_METRIC_PARAMETERS.createdAlarms],
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

      return [KPI_ALARM_METRICS_FILENAME_PREFIX, fromTime, toTime, this.query.sampling].join('-');
    },

    async exportAlarmMetricsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    getQuery() {
      return {
        ...this.interval,

        parameters: this.query.parameters,
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    async fetchList() {
      try {
        this.pending = true;

        const params = this.getQuery();
        const {
          data: alarmsMetrics,
          meta: { min_date: minDate },
        } = await this.fetchAlarmsMetricsWithoutStore({ params });

        this.alarmsMetrics = convertMetricsToTimezone(alarmsMetrics, this.$system.timezone);
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

    async exportAlarmMetricsAsCsv() {
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
