import { useStoreModuleHooks } from '@/hooks/store';

const usePbehaviorReasonStoreModule = () => useStoreModuleHooks('pbehaviorReasons');

// TODO: add comments
export const usePbehaviorReason = () => {
  const { useActions } = usePbehaviorReasonStoreModule();

  const actions = useActions({
    fetchPbehaviorReasonsList: 'fetchList',
    fetchPbehaviorReasonsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchPbehaviorReasonsListWithoutStore: 'fetchListWithoutStore',
    createPbehaviorReason: 'create',
    updatePbehaviorReason: 'update',
    removePbehaviorReason: 'remove',
    bulkHidePbehaviorReasons: 'bulkHide',
    bulkUnhidePbehaviorReasons: 'bulkUnhide',
    bulkRemovePbehaviorReasons: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
