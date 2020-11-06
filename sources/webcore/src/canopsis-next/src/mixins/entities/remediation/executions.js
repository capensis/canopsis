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
      fetchRemediationInstructionExecution: 'fetchItem',
      createRemediationInstructionExecution: 'create',
      cancelRemediationInstructionExecution: 'cancel',
      nextRemediationInstructionExecution: 'next',
      nextStepRemediationInstructionExecution: 'nextStep',
      pauseRemediationInstructionExecution: 'pause',
      previousRemediationInstructionExecution: 'previous',
      resumeRemediationInstructionExecution: 'resume',
    }),
  },
};
