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
      :next-pending="nextPending",
      :previous-pending="previousPending",
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
    previousPending: {
      type: Boolean,
      default: false,
    },
    nextPending: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
