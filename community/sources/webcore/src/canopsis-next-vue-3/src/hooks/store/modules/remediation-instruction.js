import { storeToRefs } from 'pinia'

import { useRemediationInstructionStore } from '@/stores/modules/remediation-instruction'

/**
 * Remediation instruction store facade.
 *
 * @returns {Object}
 */
export const useRemediationInstruction = () => {
  const remediationInstructionStore = useRemediationInstructionStore()
  const { items, meta, pending, fetchingParams } = storeToRefs(remediationInstructionStore)

  return {
    remediationInstructions: items,
    remediationInstructionsMeta: meta,
    remediationInstructionsPending: pending,
    remediationInstructionsFetchingParams: fetchingParams,

    fetchRemediationInstructionsList: remediationInstructionStore.fetchList,
    fetchRemediationInstructionsListWithPreviousParams:
      remediationInstructionStore.fetchListWithPreviousParams,
    fetchRemediationInstructionsListWithoutStore: remediationInstructionStore.fetchListWithoutStore,
    fetchRemediationInstructionWithoutStore: remediationInstructionStore.fetchItemWithoutStore,
    fetchRemediationInstructionApprovalWithoutStore:
      remediationInstructionStore.fetchItemApprovalWithoutStore,
    createRemediationInstruction: remediationInstructionStore.create,
    updateRemediationInstruction: remediationInstructionStore.update,
    removeRemediationInstruction: remediationInstructionStore.remove,
    updateRemediationInstructionApproval: remediationInstructionStore.updateApproval,
    rateRemediationInstruction: remediationInstructionStore.rateInstruction,
    bulkRemoveRemediationInstructions: remediationInstructionStore.bulkRemove,
    bulkEnableRemediationInstructions: remediationInstructionStore.bulkEnable,
    bulkDisableRemediationInstructions: remediationInstructionStore.bulkDisable,
  }
}
