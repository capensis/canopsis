import { unref } from 'vue';

import { MODALS, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

import { useModals } from '@/hooks/modals';

/**
 * Hook for showing execution modal based on remediation instruction type.
 *
 * @param {import('vue').Ref<Object>|Object} execution - Execution object.
 * @returns {{ showExecutionModal: Function }}
 */
export const useShowExecutionModal = (execution) => {
  const modals = useModals();

  const showExecutionModal = () => {
    const unwrappedExecution = unref(execution);

    return modals.show({
      name:
        unwrappedExecution.type === REMEDIATION_INSTRUCTION_TYPES.simpleManual
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
