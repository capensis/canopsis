import { useStoreModuleHooks } from '@/hooks/store';

const useExternalDataTableStoreModule = () => useStoreModuleHooks('externalDataTable');

export const useExternalDataTable = () => {
  const { useActions } = useExternalDataTableStoreModule();

  const actions = useActions({
    createExternalDataTable: 'create',
    updateExternalDataTable: 'update',
    removeExternalDataTable: 'remove',
    fetchExternalDataTablesListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
