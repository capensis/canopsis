import { storeToRefs } from 'pinia';

import { usePbehaviorReasonStore } from '@/stores/modules/pbehavior-reason';

/**
 * Pbehavior reason store hook.
 *
 * @returns {Object}
 */
export const usePbehaviorReason = () => {
  const pbehaviorReasonStore = usePbehaviorReasonStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(pbehaviorReasonStore);

  return {
    pbehaviorReasons: items,
    pbehaviorReasonsMeta: meta,
    pbehaviorReasonsPending: pending,
    pbehaviorReasonsFetchingParams: fetchingParams,

    fetchPbehaviorReasonsList: pbehaviorReasonStore.fetchList,
    fetchPbehaviorReasonsListWithPreviousParams: pbehaviorReasonStore.fetchListWithPreviousParams,
    fetchPbehaviorReasonsListWithoutStore: pbehaviorReasonStore.fetchListWithoutStore,
    createPbehaviorReason: pbehaviorReasonStore.create,
    updatePbehaviorReason: pbehaviorReasonStore.update,
    removePbehaviorReason: pbehaviorReasonStore.remove,
    bulkHidePbehaviorReasons: pbehaviorReasonStore.bulkHide,
    bulkUnhidePbehaviorReasons: pbehaviorReasonStore.bulkUnhide,
    bulkRemovePbehaviorReasons: pbehaviorReasonStore.bulkRemove,
  };
};
