<template>
  <v-layout>
    <v-flex
      class="mt-3"
      xs1
    >
      <v-layout
        class="fill-height"
        align-center
        column
      >
        <v-avatar
          class="white--text"
          color="primary"
          size="32"
        >
          {{ stepNumber }}
        </v-avatar>
        <span
          v-if="!isLastStep"
          class="step-line primary mt-3"
        />
      </v-layout>
    </v-flex>
    <v-flex xs11>
      <v-layout>
        <v-text-field
          :value="step.name"
          :label="$t('common.name')"
          readonly
          hide-details
          filled
        />
      </v-layout>
      <remediation-instruction-status
        :completed-at="step.completed_at"
        :time-to-complete="step.time_to_complete"
        :failed-at="step.failed_at"
      />
      <remediation-instruction-execute-step-operations
        :operations="step.operations"
        :step-number="stepNumber"
        :is-first-step="isFirstStep"
        :next-pending="nextPending"
        :previous-pending="previousPending"
        v-on="$listeners"
        @finish="showEndpointModal"
      />
    </v-flex>
  </v-layout>
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
    previousPending: {
      type: Boolean,
      default: false,
    },
    nextPending: {
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
