<template>
  <c-layout>
    <div class="mt-3 flex-grow-0" style="flex-basis: 8.33%">
      <c-layout class="fill-height" align="center" column>
        <v-avatar
          class="text-white"
          color="primary"
          size="32"
        >
          {{ operationNumber }}
        </v-avatar>
      </c-layout>
    </div>
    <div style="flex-basis: 91.67%">
      <c-text-field
        :model-value="operation.name"
        :label="t('common.name')"
        variant="filled"
        readonly
        hide-details
      />
      <remediation-instruction-status
        :completed-at="operation.completed_at"
        :time-to-complete="operation.time_to_complete"
        :failed-at="operation.failed_at"
        :started-at="operation.started_at"
      />
      <v-expand-transition>
        <c-layout
          v-if="isShownDetails"
          gap="4"
          column
        >
          <text-editor-blurred
            :value="operation.description"
            :label="t('common.description')"
            hide-details
          />
          <c-layout column>
            <span class="text-subtitle-1">
              {{ t('remediation.instructionExecute.jobs.title') }}
            </span>
            <remediation-instruction-execute-assigned-jobs-table
              v-if="operation.jobs.length"
              :jobs="operation.jobs"
              class="mt-4"
              executable
              cancelable
              @execute-job="executeJob"
              @cancel-job-execution="cancelJobExecution"
            />
          </c-layout>
          <c-layout gap="3" justify="end">
            <v-btn
              :disabled="(isFirstOperation && isFirstStep) || nextPending"
              :loading="previousPending"
              color="secondary"
              @click="$emit('previous')"
            >
              {{ t('common.previous') }}
            </v-btn>
            <v-btn
              :disabled="previousPending"
              :loading="nextPending"
              color="primary"
              @click="$emit('next')"
            >
              {{ t('common.next') }}
            </v-btn>
          </c-layout>
        </c-layout>
      </v-expand-transition>
    </div>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

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

  emits: ['next', 'previous', 'execute-job', 'cancel-job-execution'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const isStartedOperation = computed(() => !!props.operation.started_at);
    const isCompletedOperation = computed(() => !!props.operation.completed_at);
    const isFailedOperation = computed(() => !!props.operation.failed_at);
    const isShownDetails = computed(() => !isCompletedOperation.value && !isFailedOperation.value && isStartedOperation.value,
    );

    const executeJob = (job) => {
      emit('execute-job', { job, operation: props.operation });
    };

    const cancelJobExecution = (job) => {
      emit('cancel-job-execution', { job, operation: props.operation });
    };

    return {
      t,
      isShownDetails,
      executeJob,
      cancelJobExecution,
    };
  },
};
</script>
