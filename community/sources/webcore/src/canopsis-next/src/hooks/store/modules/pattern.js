import { storeToRefs } from 'pinia';

import { usePatternStore } from '@/stores/modules/pattern';

/**
 * Pattern store facade.
 *
 * @returns {Object}
 */
export const usePattern = () => {
  const store = usePatternStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(store);

  return {
    patterns: items,
    patternsMeta: meta,
    patternsPending: pending,
    patternsFetchingParams: fetchingParams,

    fetchPatternsList: store.fetchList,
    fetchPatternsListWithPreviousParams: store.fetchListWithPreviousParams,
    fetchPatternsListWithoutStore: store.fetchListWithoutStore,
    createPattern: store.create,
    updatePattern: store.update,
    removePattern: store.remove,
    bulkRemovePatterns: store.bulkRemove,
    checkPatternsEntitiesCount: store.checkPatternsEntitiesCount,
    checkPatternsAlarmsCount: store.checkPatternsAlarmsCount,
  };
};
