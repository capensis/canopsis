import { unref } from 'vue';

import { REMEDIATION_INSTRUCTION_TYPES, MODALS } from '@/constants';

import { useModals } from '@/hooks/modals';

/**
 * Hook for showing execution modal based on remediation instruction type
 *
 * @param {import('vue').Ref<Object>|Object} execution - The execution object (can be reactive ref or plain object)
 * @param {string} execution.type - The type of remediation instruction
 * @param {Object} [execution.alarm] - The alarm object associated with the execution
 * @param {string} [execution.alarm._id] - The alarm ID
 * @param {string} execution.instruction_id - The instruction ID
 * @returns {Object} Object containing showExecutionModal method
 */
export const useShowExecutionModal = (execution) => {
  const modals = useModals();

  /**
   * Shows the appropriate execution modal based on instruction type
   */
  const showExecutionModal = () => {
    const unwrappedExecution = unref(execution);

    return modals.show({
      name: unwrappedExecution.type === REMEDIATION_INSTRUCTION_TYPES.simpleManual
        ? MODALS.executeRemediationSimpleInstruction
        : MODALS.executeRemediationInstruction,
      config: {
        alarmId: unwrappedExecution.alarm?._id,
        assignedInstruction: {
          _id: unwrappedExecution.instruction_id,
          execution: unwrappedExecution,
        },
      },
    });
  };

  return {
    showExecutionModal,
  };
};
