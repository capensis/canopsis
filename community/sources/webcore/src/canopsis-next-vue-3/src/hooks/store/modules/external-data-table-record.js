import { useExternalDataTableStore } from '@/stores/modules/external-data-table';

/**
 * Provides hooks for accessing actions related to external data table records.
 *
 * @returns {Object} An object containing functions mapped to Vuex actions for managing external data table records:
 * - `fetchExternalDataTableRecordsListWithoutStore`: Fetches a list of external data table records without storing
 *    them in Vuex.
 * - `createExternalDataTableRecord`: Creates a new external data table record.
 * - `updateExternalDataTableRecord`: Updates an existing external data table record.
 * - `removeExternalDataTableRecord`: Removes a specific external data table record.
 * - `bulkRemoveExternalDataTableRecord`: Removes multiple external data table records in bulk.
 */
export const useExternalDataTableRecord = () => {
  const externalDataTableStore = useExternalDataTableStore();

  return {
    fetchExternalDataTableRecordsListWithoutStore:
      externalDataTableStore.fetchRecordsListWithoutStore,
    createExternalDataTableRecord: externalDataTableStore.createRecord,
    updateExternalDataTableRecord: externalDataTableStore.updateRecord,
    removeExternalDataTableRecord: externalDataTableStore.removeRecord,
    bulkRemoveExternalDataTableRecord: externalDataTableStore.bulkRemoveRecords,
  };
};
