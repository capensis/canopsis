<template lang="pug">
  div.position-relative
    c-progress-overlay(:pending="remediationMetricsPending")
    remediation-statistics-filters(v-model="pagination", :min-date="minDate")
    remediation-statistics-chart(
      :metrics="metrics",
      :data-type="pagination.type",
      :sampling="pagination.sampling",
      :min-date="minDate",
      responsive
    )
</template>

<script>
import {
  QUICK_RANGES,
  SAMPLINGS,
  REMEDIATION_STATISTICS_CHART_DATA_TYPE,
  DATETIME_FORMATS,
} from '@/constants';

import {
  convertDateToStartOfDayTimestampByTimezone,
} from '@/helpers/date/date';
import {
  convertStartDateIntervalToTimestampByTimezone,
  convertStopDateIntervalToTimestampByTimezone,
} from '@/helpers/date/date-intervals';
import { isMetricsQueryChanged, convertMetricsToTimezone } from '@/helpers/metrics';

import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationStatisticMixin } from '@/mixins/entities/remediation/statistic';

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
  ],
  data() {
    return {
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

    getQuery() {
      return {
        ...this.interval,

        instruction: this.query.instruction,
        sampling: this.query.sampling,
      };
    },

    fetchList() {
      this.fetchRemediationMetricsList({ params: this.getQuery() });
    },
  },
};
</script>
