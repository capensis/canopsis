import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('metricsSettings');

export const entitiesMetricsSettingsMixin = {
  methods: {
    ...mapActions({
      fetchMetricsSettingsWithoutStore: 'fetchItemWithoutStore',
      updateMetricsSettings: 'update',
    }),
  },
};
