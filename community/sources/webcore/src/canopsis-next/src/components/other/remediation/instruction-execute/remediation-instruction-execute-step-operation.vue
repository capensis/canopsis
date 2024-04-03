<template>
  <v-layout>
    <v-flex
      class="mt-3"
      xs1
    >
      <v-layout justify-center>
        <v-avatar
          class="white--text"
          color="primary"
          size="32"
        >
          {{ operationNumber }}
        </v-avatar>
      </v-layout>
    </v-flex>
    <v-flex xs11>
      <v-layout>
        <v-text-field
          :value="operation.name"
          :label="$t('common.name')"
          readonly
          hide-details
          filled
        />
      </v-layout>
      <remediation-instruction-status
        :completed-at="operation.completed_at"
        :time-to-complete="operation.time_to_complete"
        :failed-at="operation.failed_at"
        :started-at="operation.started_at"
      />
      <v-expand-transition>
        <v-layout
          v-if="isShownDetails"
          column
        >
          <text-editor-blurred
            :value="operation.description"
            :label="$t('common.description')"
            hide-details
          />
          <v-layout class="mt-4">
            <span class="text-subtitle-1">{{ $t('remediation.instructionExecute.jobs.title') }}</span>
          </v-layout>
          <v-layout column>
            <remediation-instruction-execute-assigned-jobs-table
              v-if="operation.jobs.length"
              :jobs="operation.jobs"
              class="mt-4"
              executable
              cancelable
              @execute-job="executeJob"
              @cancel-job-execution="cancelJobExecution"
            />
          </v-layout>
          <v-layout
            class="mb-2"
            justify-end
          >
            <v-btn
              :disabled="(isFirstOperation && isFirstStep) || nextPending"
              :loading="previousPending"
              class="accent"
              @click="$listeners.previous"
            >
              {{ $t('common.previous') }}
            </v-btn>
            <v-btn
              :disabled="previousPending"
              :loading="nextPending"
              class="primary mr-0"
              @click="$listeners.next"
            >
              {{ $t('common.next') }}
            </v-btn>
          </v-layout>
        </v-layout>
      </v-expand-transition>
    </v-flex>
  </v-layout>
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
    previousPending: {
      type: Boolean,
      default: false,
    },
    nextPending: {
      type: Boolean,
      default: false,
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
