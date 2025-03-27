import { useStoreModuleHooks } from '@/hooks/store';

const useExternalDataTableImportStoreModule = () => useStoreModuleHooks('externalDataTable/import');

export const useExternalDataTableImport = () => {
  const { useActions } = useExternalDataTableImportStoreModule();

  const actions = useActions({
    createExternalDataTableImport: 'create',
    fetchExternalDataTableImportData: 'fetchData',
    fetchExternalDataTableImportStatus: 'fetchStatus',
    completeExternalDataTableImport: 'complete',
  });

  return {
    ...actions,
  };
};
