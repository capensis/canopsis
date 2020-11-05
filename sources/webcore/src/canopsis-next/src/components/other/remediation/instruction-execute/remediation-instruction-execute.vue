<template lang="pug">
  div
    v-text-field(
      :value="executionInstruction.description",
      :label="$t('common.description')",
      readonly,
      box
    )
    remediation-instruction-execute-steps(:steps="executionInstruction.steps")
</template>

<script>
import entitiesRemediationInstructionExecutionMixin from '@/mixins/entities/remediation/executions';

import RemediationInstructionExecuteSteps from './remediation-instruction-execute-steps.vue';

export default {
  components: { RemediationInstructionExecuteSteps },
  mixins: [entitiesRemediationInstructionExecutionMixin],
  props: {
    executionInstructionId: {
      type: [String, Number],
      required: true,
    },
  },
  computed: {
    executionInstruction() {
      return this.getRemediationInstructionExecution(this.executionInstructionId);
    },
  },
};
</script>
