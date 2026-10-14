import { useStoreModuleHooks } from '@/hooks/store';

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
const useExternalDataTableRecordStoreModule = () => useStoreModuleHooks('externalDataTable/record');

export const useExternalDataTableRecord = () => {
  const { useActions } = useExternalDataTableRecordStoreModule();

  const actions = useActions({
    fetchExternalDataTableRecordsListWithoutStore: 'fetchListWithoutStore',
    createExternalDataTableRecord: 'create',
    updateExternalDataTableRecord: 'update',
    removeExternalDataTableRecord: 'remove',
    bulkRemoveExternalDataTableRecord: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
