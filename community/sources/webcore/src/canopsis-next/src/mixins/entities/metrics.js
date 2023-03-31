import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metrics');

export const entitiesMetricsMixin = {
  methods: {
    ...mapActions({
      fetchSliMetricsWithoutStore: 'fetchSliMetricsWithoutStore',
      fetchRatingMetricsWithoutStore: 'fetchRatingMetricsWithoutStore',
      fetchAlarmsMetricsWithoutStore: 'fetchAlarmsMetricsWithoutStore',
      fetchAggregateMetrics: 'fetchAggregateMetrics',
      createKpiAlarmExport: 'createKpiAlarmExport',
      createKpiAlarmAggregateExport: 'createKpiAlarmAggregateExport',
      createKpiRatingExport: 'createKpiRatingExport',
      createKpiSliExport: 'createKpiSliExport',
      fetchMetricExport: 'fetchMetricExport',
    }),
  },
};
