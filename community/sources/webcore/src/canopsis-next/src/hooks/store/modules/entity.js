import { useStoreModuleHooks } from '@/hooks/store';

const useEntityStoreModule = () => useStoreModuleHooks('entity');

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
