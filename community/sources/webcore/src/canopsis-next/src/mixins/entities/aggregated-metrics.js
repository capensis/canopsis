import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('aggregatedMetrics');

export const entitiesAggregatedMetricsMixin = {
  computed: {
    ...mapGetters({
      getAggregatedMetricsListByWidgetId: 'getListByWidgetId',
      getAggregatedMetricsPendingByWidgetId: 'getPendingByWidgetId',
    }),

    aggregatedMetrics() {
      return this.getAggregatedMetricsListByWidgetId(this.widget._id);
    },

    aggregatedMetricsPending() {
      return this.getAggregatedMetricsPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchAggregatedMetricsList: 'fetchList',
    }),
  },
};
