<template lang="pug">
  v-layout
    v-flex.mt-3(xs1)
      v-layout.fill-height(align-center, column)
        v-avatar.white--text(color="primary", size="32") {{ stepNumber }}
        span.step-line.primary.mt-3(v-if="!isLastStep")
    v-flex(xs11)
      v-layout
        v-text-field(
          :value="step.name",
          :label="$t('common.name')",
          readonly,
          hide-details,
          box
        )
      remediation-instruction-status(
        :completed-at="step.completed_at",
        :time-to-complete="step.time_to_complete",
        :failed-at="step.failed_at"
      )
      remediation-instruction-execute-step-operations(
        :operations="step.operations",
        :step-number="stepNumber",
        :is-first-step="isFirstStep",
        v-on="$listeners",
        @finish="showEndpointModal"
      )
</template>

<script>
import { MODALS } from '@/constants';

import RemediationInstructionExecuteStepOperations from './remediation-instruction-execute-step-operations.vue';
import RemediationInstructionStatus from './partials/remediation-instruction-status.vue';

export default {
  components: {
    RemediationInstructionExecuteStepOperations,
    RemediationInstructionStatus,
  },
  props: {
    step: {
      type: Object,
      required: true,
    },
    stepNumber: {
      type: [Number, String],
      required: true,
    },
    isFirstStep: {
      type: Boolean,
      default: false,
    },
    isLastStep: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    showEndpointModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        dialogProps: {
          persistent: true,
        },
        config: {
          hideTitle: true,
          text: this.step.endpoint,
          action: () => this.$emit('next-step', true),
          cancel: () => this.$emit('next-step', false),
        },
      });
    },
  },
};
</script>

<style lang="scss">
.step-line {
  flex: 1;
  width: 2px;
}
</style>
