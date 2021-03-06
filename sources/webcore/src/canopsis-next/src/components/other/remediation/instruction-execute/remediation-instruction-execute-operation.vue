<template lang="pug">
  v-layout
    v-flex.mt-3(xs1)
      v-layout(justify-center)
        v-avatar.white--text(color="primary", size="32") {{ operationNumber }}
    v-flex(xs11)
      v-layout
        v-text-field(
          :value="operation.name",
          :label="$t('common.name')",
          readonly,
          hide-details,
          box
        )
      remediation-instruction-status(
        :completed-at="operation.completed_at",
        :time-to-complete="operation.time_to_complete",
        :failed-at="operation.failed_at",
        :started-at="operation.started_at"
      )
      v-expand-transition
        v-layout(v-if="isShownDetails", column)
          text-editor-blurred(
            :value="operation.description",
            :label="$t('common.description')",
            hide-details
          )
          remediation-instruction-execute-assigned-jobs.mt-4(
            v-if="operation.jobs.length",
            :jobs="operation.jobs",
            :operation-id="operation.operation_id",
            :execution-id="executionId"
          )
          v-layout.mb-2(row, justify-end)
            v-btn.accent(
              :disabled="isFirstOperation && isFirstStep",
              @click="$listeners.previous"
            ) {{ $t('common.previous') }}
            v-btn.primary.mr-0(@click="$listeners.next") {{ $t('common.next') }}
</template>

<script>
import TextEditorBlurred from '@/components/common/text-editor/text-editor-blurred.vue';

import RemediationInstructionStatus from './partials/remediation-instruction-status.vue';
import RemediationInstructionExecuteAssignedJobs from './remediation-instruction-assigned-jobs.vue';

export default {
  components: {
    TextEditorBlurred,
    RemediationInstructionStatus,
    RemediationInstructionExecuteAssignedJobs,
  },
  props: {
    executionId: {
      type: String,
      required: true,
    },
    isFirstStep: {
      type: Boolean,
      default: false,
    },
    isFirstOperation: {
      type: Boolean,
      default: false,
    },
    operation: {
      type: Object,
      required: true,
    },
    operationNumber: {
      type: [Number, String],
      required: true,
    },
  },
  computed: {
    isShownDetails() {
      return !this.isCompletedOperation && !this.isFailedOperation && this.isStartedOperation;
    },

    isStartedOperation() {
      return !!this.operation.started_at;
    },

    isCompletedOperation() {
      return !!this.operation.completed_at;
    },

    isFailedOperation() {
      return !!this.operation.failed_at;
    },
  },
};
</script>
