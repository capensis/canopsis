import { useStoreModuleHooks } from '@/hooks/store';

const usePbehaviorTypeStoreModule = () => useStoreModuleHooks('pbehaviorTypes');

// TODO: add comments
export const usePbehaviorType = () => {
  const { useActions } = usePbehaviorTypeStoreModule();

  const actions = useActions({
    fetchPbehaviorTypesListWithoutStore: 'fetchListWithoutStore',
    // TODO: finish add another actions and getters
  });

  return {
    ...actions,
  };
};
