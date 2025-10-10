import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Remediation Instruction Stats Store Module.
 *
 * @returns {Object} An object containing getters and actions for the remediation instruction stats.
 */
const useRemdeitionInstructionStatsStoreModule = () => useStoreModuleHooks('remediationInstructionStats');

/**
 * Hook to access remediation instruction stats.
 *
 * @returns {Object} An object containing:
 * - Getters for remediation instruction stats, pending status, and metadata.
 * - Actions to fetch lists and summaries without using the store.
 */
export const useRemdeitionInstructionStats = () => {
  const { useGetters, useActions } = useRemdeitionInstructionStatsStoreModule();

  const getters = useGetters({
    remediationInstructionStats: 'items',
    remediationInstructionStatsPending: 'pending',
    remediationInstructionStatsMeta: 'meta',
  });

  const actions = useActions({
    fetchRemediationInstructionStatsList: 'fetchList',
    fetchRemediationInstructionStatsListWithoutStore: 'fetchListWithoutStore',
    fetchRemediationInstructionStatsSummaryWithoutStore: 'fetchSummaryWithoutStore',
    fetchRemediationInstructionStatsCommentsListWithoutStore: 'fetchCommentsWithoutStore',
    fetchRemediationInstructionStatsChangesListWithoutStore: 'fetchChangesWithoutStore',
    fetchRemediationInstructionStatsExecutionsListWithoutStore: 'fetchExecutionsWithoutStore',
  });

  return {
    ...getters,
    ...actions,
  };
};
