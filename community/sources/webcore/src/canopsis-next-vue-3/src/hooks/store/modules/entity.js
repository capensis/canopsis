import { useEntityStore } from '@/stores/modules/entity';

/**
 * Hook for managing entity-related operations and state
 *
 * @returns {Object} An object containing entity-related actions
 */
export const useEntity = () => {
  const entityStore = useEntityStore();

  return {
    fetchContextEntitiesListWithoutStore: entityStore.fetchContextEntitiesListWithoutStore,
    fetchEntityInfosLogsListWithoutStore: entityStore.fetchEntityInfosLogsListWithoutStore,
    archiveDisabledEntitiesData: entityStore.archiveDisabledEntitiesData,
    archiveUnlinkedEntitiesData: entityStore.archiveUnlinkedEntitiesData,
    cleanArchivedEntitiesData: entityStore.cleanArchivedEntitiesData,
    checkEntityStateSetting: entityStore.checkStateSetting,
    pauseEntity: entityStore.pause,
    resumeEntity: entityStore.resume,
    bulkPauseEntities: entityStore.bulkPause,
    bulkResumeEntities: entityStore.bulkResume,
  };
};
