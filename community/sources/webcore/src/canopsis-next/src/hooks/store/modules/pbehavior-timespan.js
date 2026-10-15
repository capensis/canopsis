import { usePbehaviorTimespanStore } from '@/stores/modules/pbehavior-timespan';

/**
 * Pbehavior timespan store facade.
 *
 * @returns {Object}
 */
export const usePbehaviorTimespan = () => {
  const pbehaviorTimespanStore = usePbehaviorTimespanStore();

  return {
    fetchPbehaviorTimespansListWithoutStore: pbehaviorTimespanStore.fetchListWithoutStore,
  };
};
