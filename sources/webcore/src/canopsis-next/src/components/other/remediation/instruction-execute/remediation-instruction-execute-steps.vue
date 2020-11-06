<template lang="pug">
  v-layout(column)
    remediation-instruction-execute-step(
      v-for="(step, index) in steps",
      :key="step.step_id",
      :step="step",
      :step-number="index + 1",
      :is-last="index === steps.length - 1",
      :is-first="index === 0"
    )
</template>

<script>
import RemediationInstructionExecuteStep from './remediation-instruction-execute-step.vue';

export default {
  components: { RemediationInstructionExecuteStep },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
