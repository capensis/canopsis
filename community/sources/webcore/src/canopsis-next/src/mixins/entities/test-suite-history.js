import { createNamespacedHelpers } from 'vuex';

const { mapGetters, mapActions } = createNamespacedHelpers('testSuite/history');

export const entitiesTestSuiteHistoryMixin = {
  computed: {
    ...mapGetters({
      getTestSuiteHistoryListByTestSuiteId: 'getListByTestSuiteId',
      getTestSuiteHistoryPendingByTestSuiteId: 'getPendingByTestSuiteId',
    }),

    testSuiteHistory() {
      return this.getTestSuiteHistoryListByTestSuiteId(this.testSuite.test_suite_id);
    },

    testSuiteHistoryPending() {
      return this.getTestSuiteHistoryPendingByTestSuiteId(this.testSuite.test_suite_id);
    },
  },
  methods: {
    ...mapActions({
      fetchTestSuiteHistorySummaryList: 'fetchList',
    }),
  },
};
