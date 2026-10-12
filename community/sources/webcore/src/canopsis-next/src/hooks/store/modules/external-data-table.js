import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Provides hooks for accessing actions related to external data tables.
 *
 * @returns {Object} An object containing functions mapped to Vuex actions for managing external data tables.
 */
const useExternalDataTableStoreModule = () => useStoreModuleHooks('externalDataTable');

/**
 * Provides actions for interacting with the external data table module.
 *
 * This hook maps specific action names to more descriptive function names for easier use within components.
 *
 * @returns {Object} An object containing the following actions:
 * - `createExternalDataTable`: Function to create a new external data table.
 * - `updateExternalDataTable`: Function to update an existing external data table.
 * - `removeExternalDataTable`: Function to remove a specific external data table.
 * - `fetchExternalDataTablesListWithoutStore`: Function to fetch a list of external data tables without storing them in
 *    Vuex.
 * - `fetchExternalDataTableWithoutStore`: Function to fetch a specific external data table without storing it in Vuex.
 * - `fetchExternalDataTableDataWithoutStore`: Function to fetch data for a specific external data table without storing
 *    it in Vuex.
 * - `fetchExternalDataTableSchema`: Function to fetch the schema of a specific external data table.
 * - `createExternalDataTableExport`: Function to create an export for a specific external data table.
 * - `fetchExternalDataTableExportStatus`: Function to fetch the status of an export for a specific external data table.
 */
export const useExternalDataTable = () => {
  const { useActions } = useExternalDataTableStoreModule();

  const actions = useActions({
    createExternalDataTable: 'create',
    updateExternalDataTable: 'update',
    removeExternalDataTable: 'remove',
    fetchExternalDataTablesListWithoutStore: 'fetchListWithoutStore',
    fetchExternalDataTableWithoutStore: 'fetchItemWithoutStore',
    fetchExternalDataTableDataWithoutStore: 'fetchDataWithoutStore',
    fetchExternalDataTableSchema: 'fetchSchema',

    createExternalDataTableExport: 'createExport',
    fetchExternalDataTableExportStatus: 'fetchExportStatus',
  });

  return {
    ...actions,
  };
};
