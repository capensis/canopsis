import { usePbehaviorStore } from '@/stores/modules/pbehavior';

/**
 * Pbehavior store facade.
 *
 * @returns {Object}
 */
export const usePbehavior = () => {
  const store = usePbehaviorStore();

  return {
    fetchPbehaviorsListWithoutStore: store.fetchPbehaviorsListWithoutStore,
    fetchPbehaviorReasonsListWithoutStore: store.fetchPbehaviorReasonsListWithoutStore,
    fetchPbehaviorTypesListWithoutStore: store.fetchPbehaviorTypesListWithoutStore,
    fetchPbehaviorTypesFieldList: store.fetchPbehaviorTypesFieldList,
    fetchNextPbehaviorTypePriority: store.fetchNextPbehaviorTypePriority,
  };
};
