<template lang="pug">
  div
    remediation-instruction-execute-simple-steps(
      v-if="isManualSimpleType",
      :jobs="instructionExecution.steps[0].operations[0].jobs"
    )
    template(v-else)
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
import RemediationInstructionExecuteSimpleSteps from './remediation-instruction-execute-simple-steps.vue';

export default {
  components: { RemediationInstructionExecuteSimpleSteps, RemediationInstructionExecuteSteps },
  props: {
    instructionExecution: {
      type: Object,
      required: true,
    },
    simple: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isManualSimpleType() {
      return this.instructionExecution;
    },
  },
};
</script>
