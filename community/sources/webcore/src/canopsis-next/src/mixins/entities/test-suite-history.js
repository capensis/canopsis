import { mapActions, mapState } from 'pinia';

import { useTestSuiteStore } from '@/stores/modules/test-suite';

export const entitiesTestSuiteHistoryMixin = {
  computed: {
    ...mapState(useTestSuiteStore, {
      getTestSuiteHistoryListByTestSuiteId: 'getHistoryListByTestSuiteId',
      getTestSuiteHistoryPendingByTestSuiteId: 'getHistoryPendingByTestSuiteId',
    }),

    testSuiteHistory() {
      return this.getTestSuiteHistoryListByTestSuiteId(this.testSuite.test_suite_id);
    },

    testSuiteHistoryPending() {
      return this.getTestSuiteHistoryPendingByTestSuiteId(this.testSuite.test_suite_id);
    },
  },
  methods: {
    ...mapActions(useTestSuiteStore, {
      fetchTestSuiteHistorySummaryList: 'fetchHistoryList',
    }),
  },
};
