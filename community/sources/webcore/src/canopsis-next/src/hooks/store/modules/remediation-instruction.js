import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the Remediation Instruction Store Module.
 *
 * @returns {Object} An object containing getters and actions for the remediation instruction.
 */
const useRemdeitionInstructionStoreModule = () => useStoreModuleHooks('remediationInstruction');

/**
 * Hook to access remediation instruction store.
 *
 * @returns {Object} An object containing:
 * - Getters for remediation instruction stats, pending status, and metadata.
 * - Actions to fetch lists and summaries without using the store.
 */
export const useRemdeitionInstructionStore = () => {
  const { useActions } = useRemdeitionInstructionStoreModule();

  const actions = useActions({
    fetchRemediationInstructionsListWithoutStore: 'fetchListWithoutStore',
  });

  return {
    ...actions,
  };
};
