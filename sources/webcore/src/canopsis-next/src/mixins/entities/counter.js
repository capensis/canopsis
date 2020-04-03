import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('counter');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      getCountersListByWidgetId: 'getListByWidgetId',
      getCountersPendingByWidgetId: 'getPendingByWidgetId',
    }),

    counters() {
      return this.getCountersListByWidgetId(this.widget._id);
    },

    countersPending() {
      return this.getCountersPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchCountersList: 'fetchList',
    }),
  },
};
