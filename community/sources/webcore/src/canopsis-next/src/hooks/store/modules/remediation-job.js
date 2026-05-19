import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Remediation Job Store Module.
 *
 * @returns {Object} An object containing getters and actions for the remediation job.
 */
const useRemediationJobStore = () => useStoreModuleHooks('remediationJob');

/**
 * Hook to access remediation job store.
 *
 * @returns {Object} An object containing:
 * - Getters for remediation jobs, pending status, and metadata.
 * - Actions to fetch lists and manage remediation jobs.
 */
export const useRemediationJob = () => {
  const { useGetters, useActions } = useRemediationJobStore();

  const getters = useGetters({
    remediationJobs: 'items',
    remediationJobsMeta: 'meta',
    remediationJobsPending: 'pending',
  });

  const actions = useActions({
    fetchRemediationJobsList: 'fetchList',
    fetchRemediationJobsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchRemediationJobsListWithoutStore: 'fetchListWithoutStore',
    createRemediationJob: 'create',
    updateRemediationJob: 'update',
    removeRemediationJob: 'remove',
  });

  return {
    ...getters,
    ...actions,
  };
};
