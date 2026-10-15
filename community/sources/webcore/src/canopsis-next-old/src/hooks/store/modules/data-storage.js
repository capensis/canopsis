import { useStoreModuleHooks } from '../index';

/**
 * Hook for accessing data storage settings module
 *
 * @returns {Object} An object containing data storage settings actions
 */
const useDataStorageStoreModule = () => useStoreModuleHooks('dataStorage');

/**
 * Hook for accessing data storage settings actions
 * Provides access to data storage settings functionality
 *
 * @returns {Object} An object containing data storage actions
 * @property {Function} fetchDataStorageSettingsWithoutStore - Fetches data storage settings without storing
 * @property {Function} updateDataStorageSettings - Updates data storage settings
 * @property {Function} archiveDisabledEntitiesData - Archives disabled entities
 * @property {Function} archiveUnlinkedEntitiesData - Archives unlinked entities
 * @property {Function} cleanArchivedEntitiesData - Cleans archived entities
 */
export const useDataStorage = () => {
  const { useActions } = useDataStorageStoreModule();

  const actions = useActions({
    fetchDataStorageSettingsWithoutStore: 'fetchItemWithoutStore',
    updateDataStorageSettings: 'update',
    archiveDisabledEntitiesData: 'archiveDisabledEntitiesData',
    archiveUnlinkedEntitiesData: 'archiveUnlinkedEntitiesData',
    cleanArchivedEntitiesData: 'cleanArchivedEntitiesData',
  });

  return {
    ...actions,
  };
};
