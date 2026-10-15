import { storeToRefs } from 'pinia';

import { useRemediationJobExecutionStore } from '@/stores/modules/remediation-job-execution';

/**
 * Remediation job execution store facade.
 *
 * @returns {Object}
 */
export const useRemediationJobExecution = () => {
  const store = useRemediationJobExecutionStore();
  const { outputs, getOutputById } = storeToRefs(store);

  return {
    remediationJobExecutionOutputs: outputs,
    getRemediationJobExecutionOutputById: getOutputById,
    createRemediationJobExecution: store.create,
    cancelRemediationJobExecution: store.cancel,
    fetchRemediationJobExecutionOutput: store.fetchOutput,
  };
};
