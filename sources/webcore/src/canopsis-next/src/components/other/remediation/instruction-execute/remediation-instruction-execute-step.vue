<template lang="pug">
  v-layout
    v-flex.mt-3(xs1)
      v-layout.fill-height(align-center, column)
        v-avatar.white--text(color="primary", size="32") {{ stepNumber }}
        span.step-line.primary.mt-3(v-if="!isLast")
    v-flex(xs11)
      v-layout
        v-text-field(
          :value="step.name",
          :label="$t('common.description')",
          readonly,
          hide-details,
          box
        )
      remediation-instruction-status(
        :completed-at="step.completed_at",
        :time-to-complete="step.time_to_complete",
        :failed-at="step.failed_at"
      )
      remediation-instruction-execute-operations(
        :operations="step.operations",
        :step-number="stepNumber"
      )
</template>

<script>
import RemediationInstructionExecuteOperations from './remediation-instruction-execute-operations.vue';
import RemediationInstructionStatus from './partials/remediation-instruction-status.vue';

export default {
  components: { RemediationInstructionExecuteOperations, RemediationInstructionStatus },
  props: {
    step: {
      type: Object,
      required: true,
    },
    stepNumber: {
      type: [Number, String],
      required: true,
    },
    isLast: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.step-line {
  flex: 1;
  background: red;
  width: 2px;
}
</style>
