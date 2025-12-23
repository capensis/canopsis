import { useStoreModuleHooks } from '@/hooks/store';

const useServiceStoreModule = () => useStoreModuleHooks('service');

// TODO: add comments
export const useService = () => {
  const { useActions } = useServiceStoreModule();

  const actions = useActions({
    fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore',
    bulkRemoveServices: 'bulkRemove',
    // TODO: finish add another actions and getters
  });

  return {
    ...actions,
  };
};
