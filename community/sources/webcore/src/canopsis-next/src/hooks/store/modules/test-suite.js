import { storeToRefs } from 'pinia';

import { useTestSuiteStore } from '@/stores/modules/test-suite';

/**
 * Test suite store facade.
 *
 * @returns {Object}
 */
export const useTestSuite = () => {
  const testSuiteStore = useTestSuiteStore();
  const refs = storeToRefs(testSuiteStore);

  return {
    ...refs,
    validateTestSuiteDirectory: testSuiteStore.validateDirectory,
    fetchTestSuitesList: testSuiteStore.fetchList,
    fetchTestSuiteSummaryWithoutStore: testSuiteStore.fetchSummaryWithoutStore,
    fetchTestSuiteGanttIntervalsWithoutStore: testSuiteStore.fetchItemGanttIntervalsWithoutStore,
    fetchTestSuiteDetailsWithoutStore: testSuiteStore.fetchItemDetailsWithoutStore,
    fetchTestSuiteEntityGanttIntervalsWithoutStore:
      testSuiteStore.fetchEntityGanttIntervalsWithoutStore,
    fetchTestSuiteHistoryList: testSuiteStore.fetchHistoryList,
  };
};
