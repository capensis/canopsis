import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('testSuite');

export const entitiesTestSuitesMixin = {
  computed: {
    ...mapGetters({
      getTestSuitesListByWidgetId: 'getListByWidgetId',
      getTestSuitesPendingByWidgetId: 'getPendingByWidgetId',
      getTest: 'getItem',
    }),

    testSuites() {
      return this.getTestSuitesListByWidgetId(this.widget._id);
    },

    testSuitesPending() {
      return this.getTestSuitesPendingByWidgetId(this.widget._id);
    },
  },
  methods: {
    ...mapActions({
      fetchTestSuitesList: 'fetchList',
      validateTestSuitesDirectory: 'validateDirectory',
      fetchTestSuiteItemSummaryWithoutStore: 'fetchSummaryWithoutStore',
      fetchTestSuiteDetailsWithoutStore: 'fetchItemDetailsWithoutStore',
    }),
  },
};
