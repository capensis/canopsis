import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook for accessing the entity Vuex store module
 * Creates hooks for accessing entity module's getters and actions
 *
 * @returns {Object} An object containing store module utilities
 * @property {Object} store - The Vuex store instance
 * @property {Object} module - The entity Vuex module
 * @property {Function} useGetters - Function to access entity module getters
 * @property {Function} useActions - Function to access entity module actions
 * @example
 * const { useGetters, useActions } = useEntityStoreModule();
 * const { entities } = useGetters(['entities']);
 * const { fetchList } = useActions(['fetchList']);
 */
const useEntityStoreModule = () => useStoreModuleHooks('entity');

/**
 * Hook for managing entity-related operations and state
 *
 * @returns {Object} An object containing entity-related actions
 * @property {Function} fetchContextEntitiesListWithoutStore - Fetches the list of context entities
 *                                                             without storing in the store
 * @property {Function} fetchEntityInfosLogsListWithoutStore - Fetches the list of entity infos logs
 *                                                             without storing in the store
 * @property {Function} archiveDisabledEntitiesData - Archives disabled entities data
 * @property {Function} archiveUnlinkedEntitiesData - Archives unlinked entities data
 * @property {Function} cleanArchivedEntitiesData - Cleans archived entities data
 */
export const useEntity = () => {
  const { useActions } = useEntityStoreModule();

  const actions = useActions({
    fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore',
    fetchEntityInfosLogsListWithoutStore: 'fetchEntityInfosLogsListWithoutStore',
    archiveDisabledEntitiesData: 'archiveDisabledEntitiesData',
    archiveUnlinkedEntitiesData: 'archiveUnlinkedEntitiesData',
    cleanArchivedEntitiesData: 'cleanArchivedEntitiesData',
  }); // TODO: add another actions and getters

  return {
    ...actions,
  };
};
