import { useExternalDataTableStore } from '@/stores/modules/external-data-table';

/**
 * Provides actions for interacting with the external data table import module.
 *
 * @returns {Object} An object containing the following actions:
 * - `createExternalDataTableImport`: Function to create a new external data table import.
 * - `fetchExternalDataTableImportData`: Function to fetch data for a specific external data table import.
 * - `fetchExternalDataTableImportStatus`: Function to fetch the status of a specific external data table import.
 * - `completeExternalDataTableImport`: Function to mark a specific external data table import as complete.
 */
export const useExternalDataTableImport = () => {
  const externalDataTableStore = useExternalDataTableStore();

  return {
    createExternalDataTableImport: externalDataTableStore.createImport,
    fetchExternalDataTableImportData: externalDataTableStore.fetchImportData,
    fetchExternalDataTableImportStatus: externalDataTableStore.fetchImportStatus,
    previewExternalDataTableImport: externalDataTableStore.previewImport,
    completeExternalDataTableImport: externalDataTableStore.completeImport,
  };
};
