import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('groupMetrics');

export const entitiesGroupMetricsMixin = {
  computed: {
    ...mapGetters({
      getGroupMetricsListByWidgetId: 'getListByWidgetId',
      getGroupMetricsPendingByWidgetId: 'getPendingByWidgetId',
      getGroupMetricsMetaByWidgetId: 'getMetaByWidgetId',
    }),

    groupMetrics() {
      return this.getGroupMetricsListByWidgetId(this.widget._id);
    },

    groupMetricsPending() {
      return this.getGroupMetricsPendingByWidgetId(this.widget._id);
    },

    groupMetricsMeta() {
      return this.getGroupMetricsMetaByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchGroupMetricsList: 'fetchList',
      createGroupMetricsExport: 'createExport',
    }),
  },
};
