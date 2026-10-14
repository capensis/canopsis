import { useExternalDataTableStore } from '@/stores/modules/external-data-table';

/**
 * External data table store hook.
 *
 * @returns {Object}
 */
export const useExternalDataTable = () => {
  const externalDataTableStore = useExternalDataTableStore();

  return {
    createExternalDataTable: externalDataTableStore.create,
    updateExternalDataTable: externalDataTableStore.update,
    removeExternalDataTable: externalDataTableStore.remove,
    fetchExternalDataTablesListWithoutStore: externalDataTableStore.fetchListWithoutStore,
    fetchExternalDataTableWithoutStore: externalDataTableStore.fetchItemWithoutStore,
    fetchExternalDataTableSchema: externalDataTableStore.fetchSchema,
    createExternalDataTableExport: externalDataTableStore.createExport,
    fetchExternalDataTableExportStatus: externalDataTableStore.fetchExportStatus,
    createExternalDataTableImport: externalDataTableStore.createImport,
    fetchExternalDataTableImportData: externalDataTableStore.fetchImportData,
    fetchExternalDataTableImportStatus: externalDataTableStore.fetchImportStatus,
    previewExternalDataTableImport: externalDataTableStore.previewImport,
    completeExternalDataTableImport: externalDataTableStore.completeImport,
    fetchExternalDataTableRecordsListWithoutStore:
      externalDataTableStore.fetchRecordsListWithoutStore,
    createExternalDataTableRecord: externalDataTableStore.createRecord,
    updateExternalDataTableRecord: externalDataTableStore.updateRecord,
    removeExternalDataTableRecord: externalDataTableStore.removeRecord,
    bulkRemoveExternalDataTableRecords: externalDataTableStore.bulkRemoveRecords,
  };
};
