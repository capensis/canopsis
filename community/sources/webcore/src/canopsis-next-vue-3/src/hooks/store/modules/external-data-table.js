import { useExternalDataTableStore } from '@/stores/modules/external-data-table';

/**
 * External data table store hook.
 *
 * @returns {Object}
 */
export const useExternalDataTable = () => {
  const externalDataTableStore = useExternalDataTableStore();

  return {
    fetchExternalDataTablesListWithoutStore: externalDataTableStore.fetchListWithoutStore,
    fetchExternalDataTableWithoutStore: externalDataTableStore.fetchWithoutStore,
    fetchExternalDataTableSchema: externalDataTableStore.fetchSchema,
  };
};
