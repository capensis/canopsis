import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Creates a store module hook for copy vars store module
 */
const useCopyVarsStoreModule = () => useStoreModuleHooks('copy/vars');

/**
 * Hook that provides access to copy vars store actions
 *
 * @returns {Object} Object containing copy vars actions
 * @returns {Function} returns.fetchEventFiltersCopyVarsWithoutStore - Fetches event filters copy vars without store
 * @returns {Function} returns.fetchDynamicInfosCopyVarsWithoutStore - Fetches dynamic infos copy vars without store
 */
export const useCopyVars = () => {
  const { useActions } = useCopyVarsStoreModule();

  const actions = useActions({
    fetchEventFiltersCopyVarsWithoutStore: 'fetchEventFiltersVarsWithoutStore',
    fetchDynamicInfosCopyVarsWithoutStore: 'fetchDynamicInfosVarsWithoutStore',
  });

  return actions;
};
