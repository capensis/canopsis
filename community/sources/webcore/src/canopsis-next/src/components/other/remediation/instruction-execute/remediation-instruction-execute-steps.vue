<template lang="pug">
  v-layout(column)
    remediation-instruction-execute-step(
      v-for="(step, index) in steps",
      :key="step.step_id",
      :step="step",
      :step-number="index + 1",
      :is-last-step="index === steps.length - 1",
      :is-first-step="index === 0",
      :execution-id="executionId",
      v-on="$listeners"
    )
</template>

<script>
import RemediationInstructionExecuteStep from './remediation-instruction-execute-step.vue';

export default {
  components: { RemediationInstructionExecuteStep },
  props: {
    executionId: {
      type: String,
      required: true,
    },
    steps: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
