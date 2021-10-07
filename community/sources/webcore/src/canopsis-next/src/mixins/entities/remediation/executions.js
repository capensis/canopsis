import { createNamespacedHelpers } from 'vuex';

const { mapActions, mapGetters } = createNamespacedHelpers('remediationInstructionExecution');

/**
 * @mixin
 */
export default {
  computed: {
    ...mapGetters({
      remediationInstructionExecutions: 'items',
      getRemediationInstructionExecution: 'getItemById',
      remediationInstructionExecutionsPending: 'pending',
    }),
  },
  methods: {
    ...mapActions({
      fetchPausedExecutionsWithoutStore: 'fetchPausedListWithoutStore',
      fetchRemediationInstructionExecution: 'fetchItem',
      fetchRemediationInstructionExecutionWithoutStore: 'fetchItemWithoutStore',
      createRemediationInstructionExecution: 'create',
      cancelRemediationInstructionExecution: 'cancel',
      nextOperationRemediationInstructionExecution: 'nextOperation',
      nextStepRemediationInstructionExecution: 'nextStep',
      pauseRemediationInstructionExecution: 'pause',
      previousOperationRemediationInstructionExecution: 'previousOperation',
      resumeRemediationInstructionExecution: 'resume',
      rateRemediationInstructionExecution: 'rate',
      pingRemediationInstructionExecution: 'ping',
    }),
  },
};
