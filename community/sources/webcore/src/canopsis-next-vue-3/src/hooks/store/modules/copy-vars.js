import { useCopyVarsStore } from '@/stores/modules/copy-vars';

/**
 * Hook for accessing copy variables actions.
 *
 * @returns {Object}
 */
export const useCopyVars = () => {
  const copyVarsStore = useCopyVarsStore();

  return {
    fetchEventFiltersCopyVarsWithoutStore: copyVarsStore.fetchEventFiltersCopyVarsWithoutStore,
    fetchDynamicInfosCopyVarsWithoutStore: copyVarsStore.fetchDynamicInfosCopyVarsWithoutStore,
  };
};
