import { storeToRefs } from 'pinia'

import { usePbehaviorExceptionStore } from '@/stores/modules/pbehavior-exception'

/**
 * Pbehavior exception store hook.
 *
 * @returns {Object}
 */
export const usePbehaviorException = () => {
  const pbehaviorExceptionStore = usePbehaviorExceptionStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(pbehaviorExceptionStore)

  return {
    pbehaviorExceptions: items,
    pbehaviorExceptionsMeta: meta,
    pbehaviorExceptionsPending: pending,
    pbehaviorExceptionsFetchingParams: fetchingParams,

    fetchPbehaviorExceptionsList: pbehaviorExceptionStore.fetchList,
    fetchPbehaviorExceptionsListWithPreviousParams:
      pbehaviorExceptionStore.fetchListWithPreviousParams,
    fetchPbehaviorExceptionsListWithoutStore: pbehaviorExceptionStore.fetchListWithoutStore,
    createPbehaviorException: pbehaviorExceptionStore.create,
    updatePbehaviorException: pbehaviorExceptionStore.update,
    removePbehaviorException: pbehaviorExceptionStore.remove,
    importPbehaviorException: pbehaviorExceptionStore.import,
    bulkHidePbehaviorExceptions: pbehaviorExceptionStore.bulkHide,
    bulkUnhidePbehaviorExceptions: pbehaviorExceptionStore.bulkUnhide,
    bulkRemovePbehaviorExceptions: pbehaviorExceptionStore.bulkRemove,
  }
}
