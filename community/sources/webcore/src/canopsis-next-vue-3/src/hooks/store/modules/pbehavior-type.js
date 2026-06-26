import { usePbehaviorStore } from '@/stores/modules/pbehavior';

/**
 * Pbehavior type store hook.
 *
 * @returns {Object}
 */
export const usePbehaviorType = () => {
  const pbehaviorStore = usePbehaviorStore();

  return {
    fetchPbehaviorTypesListWithoutStore: pbehaviorStore.fetchPbehaviorTypesListWithoutStore,
    fetchPbehaviorTypesFieldList: pbehaviorStore.fetchPbehaviorTypesFieldList,
    fetchNextPbehaviorTypePriority: pbehaviorStore.fetchNextPbehaviorTypePriority,
  };
};
