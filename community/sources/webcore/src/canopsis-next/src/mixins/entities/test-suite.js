import { mapActions, mapState } from 'pinia';

import { useTestSuiteStore } from '@/stores/modules/test-suite';

export const entitiesTestSuitesMixin = {
  computed: {
    ...mapState(useTestSuiteStore, {
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
    ...mapActions(useTestSuiteStore, {
      fetchTestSuitesList: 'fetchList',
      validateTestSuitesDirectory: 'validateDirectory',
      fetchTestSuiteItemSummaryWithoutStore: 'fetchSummaryWithoutStore',
      fetchTestSuiteDetailsWithoutStore: 'fetchItemDetailsWithoutStore',
    }),
  },
};
