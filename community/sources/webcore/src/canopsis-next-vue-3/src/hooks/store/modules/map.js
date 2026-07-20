import { storeToRefs } from 'pinia'

import { useMapStore } from '@/stores/modules/map'

/**
 * Map store facade.
 *
 * @returns {Object}
 */
export const useMap = () => {
  const mapStore = useMapStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(mapStore)

  return {
    maps: items,
    mapsMeta: meta,
    mapsPending: pending,
    mapsFetchingParams: fetchingParams,

    fetchMapsList: mapStore.fetchList,
    fetchMapsListWithPreviousParams: mapStore.fetchListWithPreviousParams,
    fetchMapsListWithoutStore: mapStore.fetchListWithoutStore,
    fetchMapWithoutStore: mapStore.fetchItemWithoutStore,
    fetchMapStateWithoutStore: mapStore.fetchItemStateWithoutStore,
    createMap: mapStore.create,
    updateMap: mapStore.update,
    removeMap: mapStore.remove,
    bulkRemoveMaps: mapStore.bulkRemove,
  }
}
