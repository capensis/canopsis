import { useStoreModuleHooks } from '@/hooks/store';

/**
 * Hook to use the remediation instruction execution store module.
 *
 * @returns {Object} The store module hooks for remediation instruction execution.
 */
export const useRemediationInstructionExecutionStoreModule = () => (
  useStoreModuleHooks('remediationInstructionExecution')
);

/**
 * Hook to use actions related to remediation instruction execution.
 *
 * @returns {Object} An object containing action methods for remediation instruction execution.
 * @property {Function} fetchPausedExecutionsWithoutStore - Fetches the list of paused executions without using the
 *                                                          store.
 * @property {Function} fetchRemediationInstructionExecutionWithoutStore - Fetches a remediation instruction execution
 *                                                                         item without using the store.
 * @property {Function} createRemediationInstructionExecution - Creates a new remediation instruction execution.
 * @property {Function} cancelRemediationInstructionExecution - Cancels an existing remediation instruction execution.
 * @property {Function} nextOperationRemediationInstructionExecution - Proceeds to the next operation in the remediation
 *                                                                     instruction execution.
 * @property {Function} nextStepRemediationInstructionExecution - Proceeds to the next step in the remediation
 *                                                                instruction execution.
 * @property {Function} pauseRemediationInstructionExecution - Pauses the remediation instruction execution.
 * @property {Function} previousOperationRemediationInstructionExecution - Returns to the previous operation in the
 *                                                                         remediation instruction execution.
 * @property {Function} resumeRemediationInstructionExecution - Resumes a paused remediation instruction execution.
 * @property {Function} fetchAlarmRemediationInstructionExecutionsWithoutStore - Fetches alarm remediation instruction
 *                                                                               executions without using the store.
 */
export const useRemediationInstructionExecution = () => {
  const { useActions } = useRemediationInstructionExecutionStoreModule();

  return useActions({
    fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
    fetchExecutionsStatusesWithoutStore: 'fetchStatusesListWithoutStore',
    fetchRemediationInstructionExecutionWithoutStore: 'fetchItemWithoutStore',
    createRemediationInstructionExecution: 'create',
    cancelRemediationInstructionExecution: 'cancel',
    nextOperationRemediationInstructionExecution: 'nextOperation',
    nextStepRemediationInstructionExecution: 'nextStep',
    pauseRemediationInstructionExecution: 'pause',
    previousOperationRemediationInstructionExecution: 'previousOperation',
    resumeRemediationInstructionExecution: 'resume',
    fetchAlarmRemediationInstructionExecutionsWithoutStore: 'fetchAlarmExecutionsWithoutStore',
    readRemediationInstructionExecution: 'read',
  });
};
