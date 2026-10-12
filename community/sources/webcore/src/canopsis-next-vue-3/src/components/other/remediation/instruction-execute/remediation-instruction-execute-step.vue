<template>
  <v-row class="flex-grow-0" no-gutters>
    <v-col class="mt-3" cols="1">
      <c-layout class="fill-height" align="center" column>
        <v-avatar
          class="text-white"
          color="primary"
          size="32"
        >
          {{ stepNumber }}
        </v-avatar>
        <span
          v-if="!isLastStep"
          class="step-line mt-3"
        />
      </c-layout>
    </v-col>
    <v-col cols="11">
      <c-text-field
        :model-value="step.name"
        :label="t('common.name')"
        variant="filled"
        readonly
        hide-details
      />
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
        @finish="showEndpointModal"
        @next-operation="$emit('next-operation')"
        @previous-operation="$emit('previous-operation')"
        @execute-job="$emit('execute-job', $event)"
        @cancel-job-execution="$emit('cancel-job-execution', $event)"
      />
    </v-col>
  </v-row>
</template>

<script>
import { MODALS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';

import RemediationInstructionStatus from './partials/remediation-instruction-status.vue';
import RemediationInstructionExecuteStepOperations from './remediation-instruction-execute-step-operations.vue';

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

  emits: ['next-step', 'next-operation', 'previous-operation', 'execute-job', 'cancel-job-execution'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const modals = useModals();

    const showEndpointModal = () => {
      modals.show({
        name: MODALS.confirmation,
        dialogProps: {
          persistent: true,
        },
        config: {
          hideTitle: true,
          text: props.step.endpoint,
          action: () => emit('next-step', true),
          cancel: () => emit('next-step', false),
        },
      });
    };

    return {
      t,
      showEndpointModal,
    };
  },
};
</script>

<style lang="scss" scoped>
.step-line {
  flex: 1;
  width: 2px;
  background: rgb(var(--v-theme-primary));
}
</style>
