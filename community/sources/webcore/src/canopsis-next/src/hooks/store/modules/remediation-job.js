import { storeToRefs } from 'pinia';

import { useRemediationJobStore } from '@/stores/modules/remediation-job';

/**
 * Remediation job store facade.
 *
 * @returns {Object}
 */
export const useRemediationJob = () => {
  const remediationJobStore = useRemediationJobStore();
  const { items, meta, pending, fetchingParams } = storeToRefs(remediationJobStore);

  // TODO: Ask backend to add a dedicated bulk route for remediation jobs.
  const bulkRemoveRemediationJobs = async ({ data = [] } = {}) => {
    await Promise.all(data.map(({ _id: id }) => remediationJobStore.remove({ id })));

    return [];
  };

  return {
    remediationJobs: items,
    remediationJobsMeta: meta,
    remediationJobsPending: pending,
    remediationJobsFetchingParams: fetchingParams,

    fetchRemediationJobsList: remediationJobStore.fetchList,
    fetchRemediationJobsListWithPreviousParams: remediationJobStore.fetchListWithPreviousParams,
    fetchRemediationJobsListWithoutStore: remediationJobStore.fetchListWithoutStore,
    createRemediationJob: remediationJobStore.create,
    updateRemediationJob: remediationJobStore.update,
    removeRemediationJob: remediationJobStore.remove,
    bulkRemoveRemediationJobs,
  };
};
