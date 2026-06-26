import { useEntityCategoryStore } from '@/stores/modules/entity-category';

/**
 * Entity category store facade.
 *
 * @returns {Object}
 */
export const useEntityCategory = () => {
  const store = useEntityCategoryStore();

  return {
    fetchCategoriesListWithoutStore: store.fetchCategoriesListWithoutStore,
  };
};
