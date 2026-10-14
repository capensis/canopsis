import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates hooks for accessing the map Vuex store module
 *
 * @returns {Object} Store module hooks for map namespace
 * @property {import('vuex').Store} store - Vuex store instance
 * @property {import('vuex').Module} module - Map module instance
 * @property {Function} useGetters - Hook for accessing map module getters
 * @property {Function} useActions - Hook for accessing map module actions
 */
const useMapsStoreModule = () => useStoreModuleHooks('map');

/**
 * Hook for managing map-related operations and state
 * Replaces the entitiesMapMixin functionality with Composition API
 *
 * @returns {Object} An object containing map getters, actions, and methods
 * @property {import('vue').ComputedRef} items - Map items
 * @property {import('vue').ComputedRef} pending - Pending state
 * @property {import('vue').ComputedRef} meta - Metadata
 * @property {Function} fetchList - Fetches the list of maps
 * @property {Function} createMap - Creates a new map
 * @property {Function} updateMap - Updates an existing map
 * @property {Function} removeMap - Removes a map
 * @property {Function} fetchItemWithoutStore - Fetches a map without storing it
 * @property {Function} fetchItemStateWithoutStore - Fetches a map state without storing it
 * @property {Function} bulkRemoveMaps - Bulk removes maps
 */
export const useMaps = () => {
  const { useGetters, useActions } = useMapsStoreModule();

  const getters = useGetters({
    items: 'items',
    pending: 'pending',
    meta: 'meta',
  });

  const actions = useActions({
    fetchList: 'fetchList',
    fetchMapsListWithoutStore: 'fetchListWithoutStore',
    createMap: 'create',
    updateMap: 'update',
    removeMap: 'remove',
    fetchItemWithoutStore: 'fetchItemWithoutStore',
    fetchItemStateWithoutStore: 'fetchItemStateWithoutStore',
    bulkRemoveMaps: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
