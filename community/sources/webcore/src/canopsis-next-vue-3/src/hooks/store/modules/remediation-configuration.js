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

  // TODO: Ask backend to add a dedicated bulk route for remediation configurations.
  const bulkRemoveRemediationConfigurations = async ({ data = [] } = {}) => {
    await Promise.all(data.map(({ _id: id }) => remediationConfigurationStore.remove({ id })));

    return [];
  };

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
    bulkRemoveRemediationConfigurations,
  };
};
