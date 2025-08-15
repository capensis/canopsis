import { onMounted } from 'vue';

import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates an instance of the entity info property store module hooks.
 *
 * @returns {Object} Store module hooks for the 'entityInfoProperty' namespace
 */
export const useEntityInfoPropertyStoreModule = () => useStoreModuleHooks('entityInfoProperty');

/**
 * Hook for accessing and managing entity info properties.
 * Provides CRUD operations for entity information properties.
 *
 * @typedef {Object} EntityInfoPropertyGetters
 * @property {Ref<Array>} items - Entity info properties items
 * @property {Ref<Object>} meta - Metadata for pagination
 * @property {Ref<boolean>} pending - Loading state
 *
 * @typedef {Object} EntityInfoPropertyActions
 * @property {Function} fetchList - Fetches entity info properties list
 * @property {Function} create - Creates a new entity info property
 * @property {Function} update - Updates an entity info property
 * @property {Function} remove - Removes an entity info property
 * @property {Function} fetchListWithoutStore - Fetches list without storing in store
 * @property {Function} fetchAllEntityInfoPropertiesList - Fetches all entity info properties list
 *
 * @returns {Object} Hook return object
 * @property {EntityInfoPropertyGetters} getters - All available getters
 * @property {EntityInfoPropertyActions} actions - All available actions
 */
export const useEntityInfoProperty = () => {
  const { useGetters, useActions } = useEntityInfoPropertyStoreModule();

  const getters = useGetters({
    entityInfoProperties: 'items',
    entityInfoPropertiesWithAlias: 'itemsWithAlias',
    entityInfoPropertiesWithoutAlias: 'itemsWithoutAlias',
    entityInfoPropertyMeta: 'meta',
    entityInfoPropertyPending: 'pending',
  });

  const actions = useActions({
    fetchEntityInfoPropertiesList: 'fetchList',
    createEntityInfoProperty: 'create',
    updateEntityInfoProperty: 'update',
    removeEntityInfoProperty: 'remove',
    fetchEntityInfoPropertiesListWithoutStore: 'fetchListWithoutStore',
  });

  /**
   * Fetches all entity info properties without pagination.
   * Retrieves the complete list of entity information properties from the server.
   */
  const fetchAllEntityInfoPropertiesList = async () => (
    actions.fetchEntityInfoPropertiesList({ params: { paginate: false } })
  );

  return {
    ...getters,
    ...actions,

    fetchAllEntityInfoPropertiesList,
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
