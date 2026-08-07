import { useStoreModuleHooks } from '@/hooks/store';

const useRemediationJobStoreModule = () => useStoreModuleHooks('remediationJob');

/**
 * Hook to access remediation job store.
 *
 * @returns {Object} An object containing getters and actions for remediation jobs.
 */
export const useRemediationJob = () => {
  const { useGetters, useActions } = useRemediationJobStoreModule();

  const getters = useGetters({
    remediationJobs: 'items',
    remediationJobsPending: 'pending',
    remediationJobsMeta: 'meta',
  });

  const actions = useActions({
    fetchRemediationJobsList: 'fetchList',
    fetchRemediationJobsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchRemediationJobsListWithoutStore: 'fetchListWithoutStore',
    createRemediationJob: 'create',
    updateRemediationJob: 'update',
    removeRemediationJob: 'remove',
    bulkRemoveRemediationJobs: 'bulkRemove',
  });

  return {
    ...getters,
    ...actions,
  };
};
