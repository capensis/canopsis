import { REMEDIATION_INSTRUCTION_EXECUTION_STATUSES } from '@/constants';

import { uid } from '@/helpers/uid';

/**
 * Add delimiters with instruction_modified_on to executions array for alarm timeline
 *
 * @param {Object[]} executions
 * @returns {Object[]}
 */
export const prepareRemediationInstructionExecutionsForAlarmTimeline = executions => (
  executions.reduce((acc, execution, index) => {
    const prevExecution = executions[index - 1];

    if (index > 0 && prevExecution.instruction_modified_on !== execution.instruction_modified_on) {
      acc.push({ _id: uid(), instruction_modified_on: prevExecution.instruction_modified_on });
    }

    acc.push(execution);

    return acc;
  }, [])
);

/**
 * Check if remediation instruction execution is currently running
 *
 * @param {Object} execution - The remediation instruction execution object
 * @returns {boolean} True if execution status is running, false otherwise
 */
export const remediationInstructionExecutionIsRunning = execution => (
  execution?.status === REMEDIATION_INSTRUCTION_EXECUTION_STATUSES.running
);
