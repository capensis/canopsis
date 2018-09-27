import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('stats');

export default {
  computed: {
    ...mapGetters({
      getStatsListByWidgetId: 'getListByWidgetId',
      getStatsPendingByWidgetId: 'getPendingByWidgetId',
    }),

    stats() {
      return this.getStatsListByWidgetId(this.widget._id);
    },

    statsPending() {
      return this.getStatsPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchStats: 'fetchStats',
    }),
  },
};
