<template>
  <c-layout column>
    <c-layout class="mb-4" align="center">
      <span class="text-subtitle-1 mr-5">
        {{ t('remediation.instructionExecute.jobs.title') }}
      </span>
      <v-btn
        v-if="!isFinished"
        :loading="executing"
        color="primary"
        class="ma-0"
        @click="$emit('run:jobs')"
      >
        <span>{{ t('remediation.instructionExecute.runJobs') }}</span>
        <v-icon class="ml-1">
          arrow_right
        </v-icon>
      </v-btn>
      <template v-else>
        <v-icon :color="statusIcon.color">
          {{ statusIcon.name }}
        </v-icon>
        <span class="ml-2">
          {{ statusIcon.text }}
        </span>
      </template>
    </c-layout>
    <remediation-instruction-execute-jobs-table :jobs="jobs" />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  isInstructionExecutionCompleted,
  isInstructionExecutionFinished,
} from '@/helpers/entities/remediation/instruction-execution/form';

import { useI18n } from '@/hooks/i18n';

import RemediationInstructionExecuteJobsTable from './remediation-instruction-assigned-jobs-table.vue';

export default {
  components: { RemediationInstructionExecuteJobsTable },

  props: {
    jobs: {
      type: Array,
      required: true,
    },
    instructionExecution: {
      type: Object,
      default: () => ({}),
    },
    executing: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['run:jobs'],

  setup(props) {
    const { t } = useI18n();

    const isFinished = computed(() => isInstructionExecutionFinished(props.instructionExecution));
    const isInstructionExecutionSucceeded = computed(() =>
      isInstructionExecutionCompleted(props.instructionExecution),
    );
    const statusIcon = computed(() => {
      if (isInstructionExecutionSucceeded.value) {
        return {
          name: 'check_circle',
          color: 'primary',
          text: t('remediation.instructionExecute.jobs.instructionComplete'),
        };
      }

      return {
        name: 'error',
        color: 'error',
        text: t('remediation.instructionExecute.jobs.instructionFailed'),
      };
    });

    return {
      t,
      isFinished,
      statusIcon,
    };
  },
};
</script>
