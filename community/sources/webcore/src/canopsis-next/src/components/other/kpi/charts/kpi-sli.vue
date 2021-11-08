<template lang="pug">
  div
    v-layout.ml-4.mb-4(align-center)
      c-quick-date-interval-field(
        :interval="pagination.interval",
        @input="updateInterval"
      )
    div
      kpi-sli-chart(
        :metrics="sliMetrics",
        :data-type="pagination.type",
        :sampling="pagination.sampling",
        responsive
      )
</template>

<script>
import {
  QUICK_RANGES,
  SAMPLINGS,
  KPI_SLI_GRAPH_DATA_TYPE,
} from '@/constants';

import { convertStartDateIntervalToTimestamp, convertStopDateIntervalToTimestamp } from '@/helpers/date/date-intervals';

import { entitiesMetricsMixin } from '@/mixins/entities/metrics';
import { localQueryMixin } from '@/mixins/query-local/query';

import KpiSliChart from './partials/kpi-sli-chart.vue';

export default {
  components: { KpiSliChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      sliMetrics: [],
      query: {
        sampling: SAMPLINGS.month,
        type: KPI_SLI_GRAPH_DATA_TYPE.time,
        interval: {
          from: QUICK_RANGES.last30Days.start,
          to: QUICK_RANGES.last30Days.stop,
        },
      },
    };
  },
  computed: {
    interval() {
      return {
        from: convertStartDateIntervalToTimestamp(this.pagination.interval.from),
        to: convertStopDateIntervalToTimestamp(this.pagination.interval.to),
      };
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    updateInterval(interval) {
      this.query = {
        ...this.query,
        interval,
      };
    },

    async fetchList() {
      this.sliMetrics = [
        {
          timestamp: 1631145600,
          uptime: 846720,
          downtime: 362880,
          maintenance: 241920,
        },
        {
          timestamp: 1633737600,
          uptime: 864864,
          downtime: 825552,
          maintenance: 326592,
        },
      ];
    },
  },
};
</script>
