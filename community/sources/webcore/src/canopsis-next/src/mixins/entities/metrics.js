import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('metrics');

export const entitiesMetricsMixin = {
  computed: {
    ...mapGetters({
      externalMetrics: 'externalMetrics',
      externalMetricsPending: 'pending',
    }),
  },
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
      fetchExternalMetricsList: 'fetchExternalMetricsList',
    }),
  },
};
