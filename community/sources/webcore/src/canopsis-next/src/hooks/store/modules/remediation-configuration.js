import { useStoreModuleHooks } from '@/hooks/store';

const useRemediationConfigurationStoreModule = () => useStoreModuleHooks('remediationConfiguration');

/**
 * Hook to access remediation configuration store.
 *
 * @returns {Object} An object containing getters and actions for remediation configurations.
 */
export const useRemediationConfiguration = () => {
  const { useGetters, useActions } = useRemediationConfigurationStoreModule();

  const getters = useGetters({
    remediationConfigurations: 'items',
    remediationConfigurationsPending: 'pending',
    remediationConfigurationsMeta: 'meta',
  });

  const actions = useActions({
    fetchRemediationConfigurationsList: 'fetchList',
    fetchRemediationConfigurationsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchRemediationConfigurationsListWithoutStore: 'fetchListWithoutStore',
    createRemediationConfiguration: 'create',
    updateRemediationConfiguration: 'update',
    removeRemediationConfiguration: 'remove',
    bulkRemoveRemediationConfigurations: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
