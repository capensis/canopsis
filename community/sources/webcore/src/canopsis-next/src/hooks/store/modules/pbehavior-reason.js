import { useStoreModuleHooks } from '@/hooks/store';

const usePbehaviorReasonStoreModule = () => useStoreModuleHooks('pbehaviorReasons');

// TODO: add comments
export const usePbehaviorReason = () => {
  const { useActions } = usePbehaviorReasonStoreModule();

  const actions = useActions({
    fetchPbehaviorReasonsListWithoutStore: 'fetchListWithoutStore',
    // TODO: finish add another actions and getters
  });

  return {
    ...actions,
  };
};
