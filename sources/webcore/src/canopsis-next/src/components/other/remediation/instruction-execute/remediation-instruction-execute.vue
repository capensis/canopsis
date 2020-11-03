<template lang="pug">
  div
    v-text-field(
      :value="instruction.description",
      :label="$t('common.description')",
      readonly,
      box
    )
    remediation-instruction-execute-steps(:steps="instruction.steps")
</template>

<script>
import RemediationInstructionExecuteSteps from './remediation-instruction-execute-steps.vue';

export default {
  components: { RemediationInstructionExecuteSteps },
  props: {
    instruction: {
      type: Object,
      required: true,
    },
  },
};
</script>
