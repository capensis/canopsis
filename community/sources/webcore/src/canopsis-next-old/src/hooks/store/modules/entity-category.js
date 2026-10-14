import { useStoreModuleHooks } from '@/hooks/store';

const useEntityCategoryStoreModule = () => useStoreModuleHooks('entityCategory');

// TODO: add comments
export const useEntityCategory = () => {
  const { useActions } = useEntityCategoryStoreModule();

  const actions = useActions({
    fetchCategoriesListWithoutStore: 'fetchListWithoutStore',
    // TODO: finish add another actions and getters
  });

  return {
    ...actions,
  };
};
