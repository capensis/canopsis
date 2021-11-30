<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="pending")
    kpi-alarms-filters(v-model="pagination")
    kpi-alarms-chart(
      :metrics="alarmsMetrics",
      :sampling="pagination.sampling",
      :downloading="downloading",
      responsive,
      @export:csv="exportAlarmMetricsAsCsv",
      @export:png="exportAlarmMetricsAsPng",
      @zoom="updateQueryField('interval', $event)"
    )
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
import { convertDateToString } from '@/helpers/date/date';
import { saveFile } from '@/helpers/file/files';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';
import { exportCsvMixinCreator } from '@/mixins/widget/export';

import KpiAlarmsFilters from './partials/kpi-alarms-filters.vue';

const KpiAlarmsChart = () => import(/* webpackChunkName: "Charts" */'./partials/kpi-alarms-chart.vue');

export default {
  components: { KpiAlarmsFilters, KpiAlarmsChart },
  mixins: [
    entitiesMetricsMixin,
    localQueryMixin,
    exportCsvMixinCreator({
      createExport: 'createKpiAlarmExport',
      fetchExport: 'fetchMetricExport',
      fetchExportFile: 'fetchMetricCsvFile',
    }),
  ],
  data() {
    return {
      alarmsMetrics: [],
      downloading: false,
      pending: false,
      query: {
        sampling: SAMPLINGS.day,
        parameters: [ALARM_METRIC_PARAMETERS.totalAlarms],
        filter: null,
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
    };
  },
  mounted() {
    this.fetchList();
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
      this.pending = true;

      this.alarmsMetrics = await this.fetchAlarmsMetricsWithoutStore({
        params: this.getQuery(),
      });

      this.pending = false;
    },

    async exportAlarmMetricsAsCsv() {
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
