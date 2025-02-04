import { useStoreModuleHooks } from '@/hooks/store';

const useUserStoreModule = () => useStoreModuleHooks('user');

// TODO: add comments
export const useUser = () => {
  const { useActions } = useUserStoreModule();

  const actions = useActions({
    fetchUsersListWithoutStore: 'fetchListWithoutStore',
    // TODO: finish add another actions and getters
  });

  return {
    ...actions,
  };
};
