import { useStoreModuleHooks } from '@/hooks/store';

const usePbehaviorTypeStoreModule = () => useStoreModuleHooks('pbehaviorTypes');

// TODO: add comments
export const usePbehaviorType = () => {
  const { useActions } = usePbehaviorTypeStoreModule();

  const actions = useActions({
    fetchPbehaviorTypesList: 'fetchList',
    fetchPbehaviorTypesListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchPbehaviorTypesListWithoutStore: 'fetchListWithoutStore',
    createPbehaviorType: 'create',
    updatePbehaviorType: 'update',
    removePbehaviorType: 'remove',
    fetchNextPbehaviorTypePriority: 'fetchNextPriority',
    fetchPbehaviorTypesFieldList: 'fetchFieldList',
    bulkHidePbehaviorTypes: 'bulkHide',
    bulkUnhidePbehaviorTypes: 'bulkUnhide',
    bulkRemovePbehaviorTypes: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
