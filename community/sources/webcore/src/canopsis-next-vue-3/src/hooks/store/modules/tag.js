import { storeToRefs } from 'pinia'

import { useTagStore } from '@/stores/modules/tag'

/**
 * Tag store facade.
 *
 * @returns {Object}
 */
export const useTag = () => {
  const tagStore = useTagStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(tagStore)

  return {
    tags: items,
    tagsMeta: meta,
    tagsPending: pending,
    tagsFetchingParams: fetchingParams,

    fetchTagsList: tagStore.fetchList,
    fetchTagsListWithPreviousParams: tagStore.fetchListWithPreviousParams,
    fetchTagsListWithoutStore: tagStore.fetchListWithoutStore,
    createTag: tagStore.create,
    updateTag: tagStore.update,
    removeTag: tagStore.remove,
  }
}
