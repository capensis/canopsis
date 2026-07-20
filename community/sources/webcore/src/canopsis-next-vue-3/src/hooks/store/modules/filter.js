import { storeToRefs } from 'pinia'

import { useFilterStore } from '@/stores/modules/filter'

/**
 * Filter store facade.
 *
 * @returns {Object}
 */
export const useFilter = () => {
  const filterStore = useFilterStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(filterStore)

  return {
    filters: items,
    filtersMeta: meta,
    filtersPending: pending,
    filtersFetchingParams: fetchingParams,

    fetchFiltersList: filterStore.fetchList,
    fetchFiltersListWithPreviousParams: filterStore.fetchListWithPreviousParams,
    fetchFiltersListWithoutStore: filterStore.fetchListWithoutStore,
    createFilter: filterStore.create,
    updateFilter: filterStore.update,
    removeFilter: filterStore.remove,
  }
}
