import {
  createEntityCategory,
  useEntityCategoryStore,
} from '@/stores/modules/entity-category';

/**
 * Entity category store facade.
 *
 * @returns {Object}
 */
export const useEntityCategory = () => {
  const store = useEntityCategoryStore();

  return {
    createEntityCategory,
    fetchCategoriesListWithoutStore: store.fetchCategoriesListWithoutStore,
  };
};
