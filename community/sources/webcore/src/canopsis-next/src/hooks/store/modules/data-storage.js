import { useDataStorageStore } from '@/stores/modules/data-storage';

/**
 * Data storage store facade.
 *
 * @returns {Object}
 */
export const useDataStorage = () => {
  const dataStorageStore = useDataStorageStore();

  return {
    fetchDataStorageWithoutStore: dataStorageStore.fetchItemWithoutStore,
    updateDataStorage: dataStorageStore.update,
  };
};
