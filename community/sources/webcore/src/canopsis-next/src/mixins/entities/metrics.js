import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metrics');

export const entitiesMetricsMixin = {
  methods: {
    ...mapActions({
      fetchSliMetricsWithoutStore: 'fetchSliMetricsWithoutStore',
      fetchRatingMetricsWithoutStore: 'fetchRatingMetricsWithoutStore',
      fetchAlarmsMetricsWithoutStore: 'fetchAlarmsMetricsWithoutStore',
      fetchAggregateMetricsWithoutStore: 'fetchAggregateMetricsWithoutStore',
      createKpiAlarmExport: 'createKpiAlarmExport',
      createKpiRatingExport: 'createKpiRatingExport',
      createKpiSliExport: 'createKpiSliExport',
      fetchMetricExport: 'fetchMetricExport',
    }),
  },
};
