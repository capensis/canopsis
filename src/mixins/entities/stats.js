import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('stats');

export default {
  computed: {
    ...mapGetters({
      getStatsListByWidgetId: 'getListByWidgetId',
      getStatsPendingByWidgetId: 'getPendingByWidgetId',
      getStatsErrorByWidgetId: 'getErrorByWidgetId',
    }),

    stats() {
      return this.getStatsListByWidgetId(this.widget._id);
    },

    statsPending() {
      return this.getStatsPendingByWidgetId(this.widget._id);
    },

    statsError() {
      return this.getStatsErrorByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchStats: 'fetchStats',
    }),
  },
};
