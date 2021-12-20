import { createNamespacedHelpers } from 'vuex';

const { mapActions } = createNamespacedHelpers('remediationInstructionExecution');

/**
 * @mixin
 */
export const entitiesRemediationInstructionExecutionMixin = {
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
      fetchRemediationInstructionExecutionWithoutStore: 'fetchItemWithoutStore',
      createRemediationInstructionExecution: 'create',
      cancelRemediationInstructionExecution: 'cancel',
      nextOperationRemediationInstructionExecution: 'nextOperation',
      nextStepRemediationInstructionExecution: 'nextStep',
      pauseRemediationInstructionExecution: 'pause',
      previousOperationRemediationInstructionExecution: 'previousOperation',
      resumeRemediationInstructionExecution: 'resume',
      rateRemediationInstructionExecution: 'rate',
    }),
  },
};
