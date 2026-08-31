import { storeToRefs } from 'pinia';
import { onMounted } from 'vue';

import { useEntityInfoPropertyStore } from '@/stores/modules/entity-info-property';

/**
 * Entity info property store facade.
 *
 * @returns {Object}
 */
export const useEntityInfoProperty = () => {
  const entityInfoPropertyStore = useEntityInfoPropertyStore();
  const { items, meta, pending, fetchingParams, itemsWithAlias, itemsWithoutAlias } =
    storeToRefs(entityInfoPropertyStore);

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

    fetchAllEntityInfoPropertiesList: () =>
      entityInfoPropertyStore.fetchList({ params: { paginate: false } }),
  };
};

/**
 * Hook for automatically fetching entity info properties on component mount.
 * This hook combines entity info property functionality with automatic data fetching
 * when the component is mounted, ensuring that entity info properties are always
 * available when the component initializes.
 */
export const useEntityInfoPropertyFetching = () => {
  const { fetchAllEntityInfoPropertiesList } = useEntityInfoProperty();

  onMounted(fetchAllEntityInfoPropertiesList);
};
