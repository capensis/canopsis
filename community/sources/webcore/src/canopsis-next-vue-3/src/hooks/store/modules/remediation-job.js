import { storeToRefs } from 'pinia'

import { useRemediationJobStore } from '@/stores/modules/remediation-job'

/**
 * Remediation job store facade.
 *
 * @returns {Object}
 */
export const useRemediationJob = () => {
  const remediationJobStore = useRemediationJobStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(remediationJobStore)

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
  }
}
