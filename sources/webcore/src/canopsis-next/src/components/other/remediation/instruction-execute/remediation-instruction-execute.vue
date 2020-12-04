<template lang="pug">
  div
    v-text-field(
      :value="executionInstruction.description",
      :label="$t('common.description')",
      readonly,
      box
    )
    remediation-instruction-execute-steps(
      :steps="executionInstruction.steps",
      :execution-id="executionInstruction._id",
      @next-step="nextStep"
    )
</template>

<script>
import entitiesRemediationInstructionExecutionMixin from '@/mixins/entities/remediation/executions';

import RemediationInstructionExecuteSteps from './remediation-instruction-execute-steps.vue';

export default {
  components: { RemediationInstructionExecuteSteps },
  mixins: [
    entitiesRemediationInstructionExecutionMixin,
  ],
  props: {
    executionInstruction: {
      type: Object,
      required: true,
    },
  },
  methods: {
    nextStep(success) {
      this.nextStepRemediationInstructionExecution({
        id: this.executionInstruction._id,
        data: { failed: !success },
      });
    },
  },
};
</script>
