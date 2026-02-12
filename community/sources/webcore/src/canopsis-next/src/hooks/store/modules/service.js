import { useStoreModuleHooks } from '@/hooks/store';

const useServiceStoreModule = () => useStoreModuleHooks('service');

/**
 * Hook for accessing service store actions
 *
 * @returns {Object} Service store actions
 */
export const useService = () => {
  const { useActions } = useServiceStoreModule();

  const actions = useActions({
    fetchEntityInfosKeysWithoutStore: 'fetchInfosKeysWithoutStore',
    fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
    fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    fetchEntityUpstreamWithoutStore: 'fetchUpstreamWithoutStore',
    fetchEntityDownstreamsWithoutStore: 'fetchDownstreamsWithoutStore',
    // TODO: finish add another actions and getters
  });

  return {
    ...actions,
  };
};
