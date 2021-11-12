<template lang="pug">
  div
    kpi-sli-filters(v-model="pagination")
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

import KpiSliFilters from './partials/kpi-sli-filters.vue';

const KpiSliChart = () => import(/* webpackChunkName: "Charts" */ './partials/kpi-sli-chart.vue');

export default {
  components: { KpiSliFilters, KpiSliChart },
  mixins: [entitiesMetricsMixin, localQueryMixin],
  data() {
    return {
      sliMetrics: [],
      query: {
        sampling: SAMPLINGS.day,
        type: KPI_SLI_GRAPH_DATA_TYPE.percent,
        filter: null,
        interval: {
          from: QUICK_RANGES.last30Days.start,
          to: QUICK_RANGES.last30Days.stop,
        },
      },
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    getQuery() {
      return {
        from: convertStartDateIntervalToTimestamp(this.query.interval.from),
        to: convertStopDateIntervalToTimestamp(this.query.interval.to),
        in_percents: this.query.type === KPI_SLI_GRAPH_DATA_TYPE.percent,
        sampling: this.query.sampling,
        filter: this.query.filter,
      };
    },

    async fetchList() {
      this.sliMetrics = await this.fetchSliMetricsWithoutStore({
        params: this.getQuery(),
      });
    },
  },
};
</script>

<style scoped lang="scss">
.kpi-sli-filters {
  &__sampling {
    max-width: 200px;
  }
}
</style>
