import { storeToRefs } from 'pinia'

import { useShareTokenStore } from '@/stores/modules/share-token'

/**
 * Share token store facade.
 *
 * @returns {Object}
 */
export const useShareToken = () => {
  const shareTokenStore = useShareTokenStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(shareTokenStore)

  return {
    shareTokens: items,
    shareTokensMeta: meta,
    shareTokensPending: pending,
    shareTokensFetchingParams: fetchingParams,

    fetchShareTokensList: shareTokenStore.fetchList,
    fetchShareTokensListWithPreviousParams: shareTokenStore.fetchListWithPreviousParams,
    createShareToken: shareTokenStore.create,
    updateShareToken: shareTokenStore.update,
    removeShareToken: shareTokenStore.remove,
  }
}
