import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Remediation Instruction Store Module.
 *
 * @returns {Object} An object containing getters and actions for the remediation instruction.
 */
const useRemdeitionInstructionStore = () => useStoreModuleHooks('remediationInstruction');

/**
 * Hook to access remediation instruction store.
 *
 * @returns {Object} An object containing:
 * - Getters for remediation instruction stats, pending status, and metadata.
 * - Actions to fetch lists and summaries without using the store.
 */
export const useRemdeitionInstruction = () => {
  const { useGetters, useActions } = useRemdeitionInstructionStore();

  const getters = useGetters({
    remediationInstructions: 'items',
    remediationInstructionsMeta: 'meta',
    remediationInstructionsPending: 'pending',
  });

  const actions = useActions({
    fetchRemediationInstructionsList: 'fetchList',
    fetchRemediationInstructionsListWithPreviousParams: 'fetchListWithPreviousParams',
    fetchRemediationInstructionsListWithoutStore: 'fetchListWithoutStore',
    createRemediationInstruction: 'create',
    updateRemediationInstruction: 'update',
    removeRemediationInstruction: 'remove',
    rateRemediationInstruction: 'rateInstruction',
    updateRemediationInstructionApproval: 'updateApproval',
    bulkEnableRemediationInstructions: 'bulkEnable',
    bulkDisableRemediationInstructions: 'bulkDisable',
  });

  return {
    ...getters,
    ...actions,
  };
};
