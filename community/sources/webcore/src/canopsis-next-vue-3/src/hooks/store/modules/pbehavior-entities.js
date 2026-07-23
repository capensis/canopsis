import { usePbehaviorEntitiesStore } from '@/stores/modules/pbehavior-entities';

/**
 * Pbehavior entities store facade.
 *
 * @returns {Object}
 */
export const usePbehaviorEntities = () => {
  const pbehaviorEntitiesStore = usePbehaviorEntitiesStore();

  return {
    fetchPbehaviorEntitiesListWithoutStore: pbehaviorEntitiesStore.fetchListWithoutStore,
  };
};
