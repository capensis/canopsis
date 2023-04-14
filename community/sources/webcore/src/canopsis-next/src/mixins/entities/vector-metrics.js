import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('vectorMetrics');

export const entitiesVectorMetricsMixin = {
  computed: {
    ...mapGetters({
      getVectorMetricsListByWidgetId: 'getListByWidgetId',
      getVectorMetricsPendingByWidgetId: 'getPendingByWidgetId',
      getVectorMetricsMetaByWidgetId: 'getMetaByWidgetId',
    }),

    vectorMetrics() {
      return this.getVectorMetricsListByWidgetId(this.widget._id);
    },

    vectorMetricsPending() {
      return this.getVectorMetricsPendingByWidgetId(this.widget._id);
    },

    vectorMetricsMeta() {
      return this.getVectorMetricsMetaByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchVectorMetricsList: 'fetchList',
    }),
  },
};
