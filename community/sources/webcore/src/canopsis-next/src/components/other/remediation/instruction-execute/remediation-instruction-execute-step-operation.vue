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
          v-layout.mt-4(row)
            span.subheading {{ $t('remediation.instructionExecute.jobs.title') }}
          v-layout(column)
            remediation-instruction-execute-assigned-jobs-table.mt-4(
              v-if="operation.jobs.length",
              :jobs="operation.jobs",
              executable,
              cancelable,
              @execute-job="executeJob",
              @cancel-job-execution="cancelJobExecution"
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
import RemediationInstructionExecuteAssignedJobsTable from './remediation-instruction-assigned-jobs-table.vue';

export default {
  components: {
    TextEditorBlurred,
    RemediationInstructionStatus,
    RemediationInstructionExecuteAssignedJobsTable,
  },
  props: {
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
  methods: {
    executeJob(job) {
      this.$emit('execute-job', { job, operation: this.operation });
    },

    cancelJobExecution(job) {
      this.$emit('cancel-job-execution', { job, operation: this.operation });
    },
  },
};
</script>
