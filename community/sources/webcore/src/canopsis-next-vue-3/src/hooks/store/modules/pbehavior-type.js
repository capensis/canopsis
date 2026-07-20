import { storeToRefs } from 'pinia'

import { usePbehaviorTypeStore } from '@/stores/modules/pbehavior-type'

/**
 * Pbehavior type store hook.
 *
 * @returns {Object}
 */
export const usePbehaviorType = () => {
  const pbehaviorTypeStore = usePbehaviorTypeStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(pbehaviorTypeStore)

  return {
    pbehaviorTypes: items,
    pbehaviorTypesMeta: meta,
    pbehaviorTypesPending: pending,
    pbehaviorTypesFetchingParams: fetchingParams,

    fetchPbehaviorTypesList: pbehaviorTypeStore.fetchList,
    fetchPbehaviorTypesListWithPreviousParams: pbehaviorTypeStore.fetchListWithPreviousParams,
    fetchPbehaviorTypesListWithoutStore: pbehaviorTypeStore.fetchListWithoutStore,
    fetchPbehaviorTypesFieldList: pbehaviorTypeStore.fetchListWithoutStore,
    fetchNextPbehaviorTypePriority: pbehaviorTypeStore.fetchNextPriority,
    createPbehaviorType: pbehaviorTypeStore.create,
    updatePbehaviorType: pbehaviorTypeStore.update,
    removePbehaviorType: pbehaviorTypeStore.remove,
    bulkHidePbehaviorTypes: pbehaviorTypeStore.bulkHide,
    bulkUnhidePbehaviorTypes: pbehaviorTypeStore.bulkUnhide,
    bulkRemovePbehaviorTypes: pbehaviorTypeStore.bulkRemove,
  }
}
