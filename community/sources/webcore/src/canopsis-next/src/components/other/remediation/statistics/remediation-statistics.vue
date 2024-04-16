<template>
  <div class="position-relative">
    <c-progress-overlay :pending="remediationMetricsPending" />
    <remediation-statistics-filters
      v-model="options"
      :min-date="minDate"
    />
    <remediation-statistics-chart
      :metrics="metrics"
      :data-type="query.type"
      :sampling="query.sampling"
      :min-date="minDate"
      :downloading="downloading"
      responsive
      @export:csv="exportRemediationStatisticsAsCsv"
      @export:png="exportRemediationStatisticsAsPng"
    />
  </div>
</template>

<script>
import { REMEDIATION_STATISTICS_FILENAME_PREFIX } from '@/config';
import {
  QUICK_RANGES,
  SAMPLINGS,
  REMEDIATION_STATISTICS_CHART_DATA_TYPE,
  DATETIME_FORMATS,
  REMEDIATION_INSTRUCTION_TYPES,
} from '@/constants';

import { convertDateToStartOfDayTimestampByTimezone, convertDateToString } from '@/helpers/date/date';
import {
  convertStartDateIntervalToTimestampByTimezone,
  convertStopDateIntervalToTimestampByTimezone,
} from '@/helpers/date/date-intervals';
import { convertMetricsToTimezone } from '@/helpers/entities/metric/list';
import { isMetricsQueryChanged } from '@/helpers/entities/metric/query';
import { getExportMetricDownloadFileUrl } from '@/helpers/entities/metric/url';
import { saveFile } from '@/helpers/file/files';

import { localQueryMixin } from '@/mixins/query/query';
import { entitiesRemediationStatisticMixin } from '@/mixins/entities/remediation/statistic';
import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { exportMixinCreator } from '@/mixins/widget/export';

import RemediationStatisticsFilters from './partials/remediation-statistics-filters.vue';

const RemediationStatisticsChart = () => import(/* webpackChunkName: "Remediation" */ './partials/remediation-statistics-chart.vue');

export default {
  inject: ['$system'],
  components: {
    RemediationStatisticsFilters,
    RemediationStatisticsChart,
  },
  mixins: [
    localQueryMixin,
    entitiesRemediationStatisticMixin,
    entitiesMetricsMixin,
    exportMixinCreator({
      createExport: 'createRemediationExport',
      fetchExport: 'fetchMetricExport',
    }),
  ],
  data() {
    return {
      downloading: false,
      query: {
        sampling: SAMPLINGS.day,
        type: REMEDIATION_STATISTICS_CHART_DATA_TYPE.percent,
        instruction: '',
        interval: {
          from: QUICK_RANGES.last7Days.start,
          to: QUICK_RANGES.last7Days.stop,
        },
      },
    };
  },
  computed: {
    minDate() {
      const { minDate } = this.remediationMetricsMeta ?? {};

      return minDate ? convertDateToStartOfDayTimestampByTimezone(minDate, this.$system.timezone) : null;
    },

    metrics() {
      return convertMetricsToTimezone(
        this.remediationMetrics.filter(({ timestamp }) => timestamp >= this.minDate),
        this.$system.timezone,
      );
    },

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
    minDate() {
      if (this.interval.from < this.minDate) {
        this.updateQueryField('interval', { ...this.query.interval, from: this.minDate });
      }
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    customQueryCondition(query, oldQuery) {
      return isMetricsQueryChanged(query, oldQuery, this.minDate);
    },

    getInstructionString(instruction) {
      if (!instruction) {
        return this.$t('remediation.statistic.allInstructions');
      }

      return this.isInstructionType(instruction)
        ? this.$t(`remediation.instruction.types.${instruction}`)
        : instruction;
    },

    getFileName() {
      const fromTime = convertDateToString(this.interval.from, DATETIME_FORMATS.short);
      const toTime = convertDateToString(this.interval.to, DATETIME_FORMATS.short);

      const instruction = this.getInstructionString(this.query.instruction);

      return [
        REMEDIATION_STATISTICS_FILENAME_PREFIX,
        fromTime,
        toTime,
        this.query.sampling,
        this.query.type,
        instruction,
      ].join('-');
    },

    async exportRemediationStatisticsAsPng(blob) {
      try {
        await saveFile(blob, this.getFileName());
      } catch (err) {
        console.error(err);

        this.$popups.error({ text: err.message || this.$t('errors.default') });
      }
    },

    async exportRemediationStatisticsAsCsv() {
      this.downloading = true;

      try {
        const fileData = await this.generateFile({
          data: this.getQuery(),
        });

        this.downloadFile(getExportMetricDownloadFileUrl(fileData._id));
      } catch (err) {
        this.$popups.error({ text: this.$t('kpi.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },

    isInstructionType(instruction) {
      return [REMEDIATION_INSTRUCTION_TYPES.manual, REMEDIATION_INSTRUCTION_TYPES.auto].includes(instruction);
    },

    getQuery() {
      const { instruction } = this.query;
      const query = {
        ...this.interval,

        sampling: this.query.sampling,
      };

      if (this.isInstructionType(instruction)) {
        query.instruction_type = instruction;
      } else if (instruction) {
        query.instruction = instruction;
      }

      return query;
    },

    fetchList() {
      this.fetchRemediationMetricsList({ params: this.getQuery() });
    },
  },
};
</script>
