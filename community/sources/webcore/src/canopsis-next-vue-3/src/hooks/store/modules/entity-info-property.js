import { storeToRefs } from 'pinia'

import { useEntityInfoPropertyStore } from '@/stores/modules/entity-info-property'

/**
 * Entity info property store facade.
 *
 * @returns {Object}
 */
export const useEntityInfoProperty = () => {
  const entityInfoPropertyStore = useEntityInfoPropertyStore()
  const { items, meta, pending, fetchingParams, itemsWithAlias, itemsWithoutAlias } =
    storeToRefs(entityInfoPropertyStore)

  return {
    entityInfoProperties: items,
    entityInfoPropertiesMeta: meta,
    entityInfoPropertiesPending: pending,
    entityInfoPropertiesFetchingParams: fetchingParams,
    entityInfoPropertiesWithAlias: itemsWithAlias,
    entityInfoPropertiesWithoutAlias: itemsWithoutAlias,

    fetchEntityInfoPropertiesList: entityInfoPropertyStore.fetchList,
    fetchEntityInfoPropertiesListWithPreviousParams:
      entityInfoPropertyStore.fetchListWithPreviousParams,
    fetchEntityInfoPropertiesListWithoutStore: entityInfoPropertyStore.fetchListWithoutStore,
    createEntityInfoProperty: entityInfoPropertyStore.create,
    updateEntityInfoProperty: entityInfoPropertyStore.update,
    removeEntityInfoProperty: entityInfoPropertyStore.remove,
    bulkRemoveEntityInfoProperties: entityInfoPropertyStore.bulkRemove,
  }
}
