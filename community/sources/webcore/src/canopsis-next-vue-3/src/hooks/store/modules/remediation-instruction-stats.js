import { storeToRefs } from 'pinia'

import { useRemediationInstructionStatsStore } from '@/stores/modules/remediation-instruction-stats'

/**
 * Remediation instruction stats store facade.
 *
 * @returns {Object}
 */
export const useRemediationInstructionStats = () => {
  const remediationInstructionStatsStore = useRemediationInstructionStatsStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(remediationInstructionStatsStore)

  return {
    remediationInstructionStats: items,
    remediationInstructionStatsMeta: meta,
    remediationInstructionStatsPending: pending,
    remediationInstructionStatsFetchingParams: fetchingParams,

    fetchRemediationInstructionStatsList: remediationInstructionStatsStore.fetchList,
    fetchRemediationInstructionStatsListWithPreviousParams:
      remediationInstructionStatsStore.fetchListWithPreviousParams,
    fetchRemediationInstructionStatsListWithoutStore:
      remediationInstructionStatsStore.fetchListWithoutStore,
    fetchRemediationInstructionStatsSummaryWithoutStore:
      remediationInstructionStatsStore.fetchSummaryWithoutStore,
    fetchRemediationInstructionStatsChangesWithoutStore:
      remediationInstructionStatsStore.fetchChangesWithoutStore,
    fetchRemediationInstructionStatsExecutionsWithoutStore:
      remediationInstructionStatsStore.fetchExecutionsWithoutStore,
    fetchRemediationInstructionStatsCommentsWithoutStore:
      remediationInstructionStatsStore.fetchCommentsWithoutStore,
  }
}
