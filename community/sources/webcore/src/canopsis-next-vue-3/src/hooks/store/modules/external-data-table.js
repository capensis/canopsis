import { useExternalDataTableStore } from '@/stores/modules/external-data-table'

/**
 * External data table store hook.
 *
 * @returns {Object}
 */
export const useExternalDataTable = () => {
  const externalDataTableStore = useExternalDataTableStore()

  return {
    fetchExternalDataTablesListWithoutStore: externalDataTableStore.fetchListWithoutStore,
    fetchExternalDataTableWithoutStore: externalDataTableStore.fetchWithoutStore,
    fetchExternalDataTableSchema: externalDataTableStore.fetchSchema,
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
  }
}
