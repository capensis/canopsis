import { useStoreModuleHooks } from '@/hooks/store';

const usePbehaviorExceptionStoreModule = () => useStoreModuleHooks('pbehaviorExceptions');

/**
 * Hook for accessing pbehavior exception store module actions
 *
 * @returns {Object} An object containing pbehavior exception actions
 */
export const usePbehaviorException = () => {
  const { useActions } = usePbehaviorExceptionStoreModule();

  const actions = useActions({
    fetchPbehaviorExceptionsList: 'fetchList',
    fetchPbehaviorExceptionsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchPbehaviorExceptionsListWithoutStore: 'fetchListWithoutStore',
    createPbehaviorException: 'create',
    updatePbehaviorException: 'update',
    removePbehaviorException: 'remove',
    importPbehaviorException: 'import',
    bulkHidePbehaviorExceptions: 'bulkHide',
    bulkUnhidePbehaviorExceptions: 'bulkUnhide',
    bulkRemovePbehaviorExceptions: 'bulkRemove',
  });

  return {
    ...actions,
  };
};
