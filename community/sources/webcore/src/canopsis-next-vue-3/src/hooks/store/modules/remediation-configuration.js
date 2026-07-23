import { storeToRefs } from 'pinia';

import { useRemediationConfigurationStore } from '@/stores/modules/remediation-configuration';

/**
 * Remediation configuration store facade.
 *
 * @returns {Object}
 */
export const useRemediationConfiguration = () => {
  const remediationConfigurationStore = useRemediationConfigurationStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(remediationConfigurationStore);

  return {
    remediationConfigurations: items,
    remediationConfigurationsMeta: meta,
    remediationConfigurationsPending: pending,
    remediationConfigurationsFetchingParams: fetchingParams,

    fetchRemediationConfigurationsList: remediationConfigurationStore.fetchList,
    fetchRemediationConfigurationsListWithPreviousParams:
      remediationConfigurationStore.fetchListWithPreviousParams,
    fetchRemediationConfigurationsListWithoutStore:
      remediationConfigurationStore.fetchListWithoutStore,
    createRemediationConfiguration: remediationConfigurationStore.create,
    updateRemediationConfiguration: remediationConfigurationStore.update,
    removeRemediationConfiguration: remediationConfigurationStore.remove,
  };
};
