import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('techMetrics');

export const entitiesTechMetricsMixin = {
  methods: {
    ...mapActions({
      createTechMetricsExport: 'createTechMetricsExport',
      fetchTechMetricsExport: 'fetchTechMetricsExport',
    }),
  },
};
