import { useRemediationInstructionExecutionStore } from '@/stores/modules/remediation-instruction-execution';

/**
 * Remediation instruction execution store facade.
 *
 * @returns {Object}
 */
export const useRemediationInstructionExecution = () => {
  const store = useRemediationInstructionExecutionStore();

  return {
    fetchRemediationInstructionExecutionWithoutStore: store.fetchItemWithoutStore,
    fetchPausedRemediationInstructionExecutionsListWithoutStore: store.fetchPausedListWithoutStore,
    fetchRemediationInstructionExecutionStatusesListWithoutStore:
      store.fetchStatusesListWithoutStore,
    createRemediationInstructionExecution: store.create,
    updateRemediationInstructionExecution: store.update,
    cancelRemediationInstructionExecution: store.cancel,
    nextRemediationInstructionExecutionOperation: store.nextOperation,
    nextRemediationInstructionExecutionStep: store.nextStep,
    pauseRemediationInstructionExecution: store.pause,
    previousRemediationInstructionExecutionOperation: store.previousOperation,
    resumeRemediationInstructionExecution: store.resume,
    fetchAlarmRemediationInstructionExecutionsWithoutStore: store.fetchAlarmExecutionsWithoutStore,
    readRemediationInstructionExecution: store.read,
  };
};
