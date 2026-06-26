import { useEntityStore } from '@/stores/modules/entity';

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
 * @property {Function} checkStateSetting - Checks the state setting for an entity
 */
export const useEntity = () => {
  const entityStore = useEntityStore();

  return {
    fetchContextEntitiesListWithoutStore: entityStore.fetchContextEntitiesListWithoutStore,
    fetchEntityInfosLogsListWithoutStore: entityStore.fetchEntityInfosLogsListWithoutStore,
    archiveDisabledEntitiesData: entityStore.archiveDisabledEntitiesData,
    archiveUnlinkedEntitiesData: entityStore.archiveUnlinkedEntitiesData,
    cleanArchivedEntitiesData: entityStore.cleanArchivedEntitiesData,
    checkStateSetting: entityStore.checkStateSetting,
  };
};
