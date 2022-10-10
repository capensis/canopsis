<template lang="pug">
  div
    v-text-field(
      :value="instructionExecution.description",
      :label="$t('common.description')",
      readonly,
      box
    )
    remediation-instruction-execute-steps(
      :steps="instructionExecution.steps",
      :execution-id="instructionExecution._id",
      v-on="$listeners"
    )
</template>

<script>
import RemediationInstructionExecuteSteps from './remediation-instruction-execute-steps.vue';

export default {
  components: { RemediationInstructionExecuteSteps },
  props: {
    instructionExecution: {
      type: Object,
      required: true,
    },
  },
};
</script>
