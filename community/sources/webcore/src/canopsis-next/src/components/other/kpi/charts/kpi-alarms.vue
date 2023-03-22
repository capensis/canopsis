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
import { KPI_ALARM_METRICS_FILENAME_PREFIX } from '@/config';
import {
  QUICK_RANGES,
  ALARM_METRIC_PARAMETERS,
  SAMPLINGS,
  DATETIME_FORMATS,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';
import { convertDateToStartOfDayTimestamp, convertDateToString } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';
import { isMetricsQueryChanged } from '@/helpers/metrics';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { exportCsvMixinCreator } from '@/mixins/widget/export';

import KpiAlarmsFilters from './partials/kpi-alarms-filters.vue';
import KpiErrorOverlay from './partials/kpi-error-overlay.vue';

const KpiAlarmsChart = () => import(/* webpackChunkName: "Charts" */'./partials/kpi-alarms-chart.vue');

export default {
  components: { KpiErrorOverlay, KpiAlarmsFilters, KpiAlarmsChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportCsvMixinCreator({
      createExport: 'createKpiAlarmExport',
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
        from: convertStartDateIntervalToTimestamp(this.query.interval.from),
        to: convertStopDateIntervalToTimestamp(this.query.interval.to),
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
        from: convertStartDateIntervalToTimestamp(this.query.interval.from),
        to: convertStopDateIntervalToTimestamp(this.query.interval.to),
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

        this.alarmsMetrics = alarmsMetrics;
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

    async exportAlarmMetricsAsCsv() {
      this.downloading = true;

      try {
        await this.exportAsCsv({
          name: this.getFileName(),
          data: this.getQuery(),
        });
      } catch (err) {
        this.$popups.error({ text: this.$t('kpi.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>
