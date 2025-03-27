import { useStoreModuleHooks } from '@/hooks/store';

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
