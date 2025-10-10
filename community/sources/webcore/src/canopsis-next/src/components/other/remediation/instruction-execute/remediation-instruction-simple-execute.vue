<template>
  <v-layout column>
    <v-layout
      class="mb-4"
      align-center
    >
      <span class="text-subtitle-1 mr-5">{{ $t('remediation.instructionExecute.jobs.title') }}</span>
      <v-btn
        v-if="!isInstructionExecutionFinished"
        :loading="executed"
        class="primary ma-0"
        @click="$emit('run:jobs')"
      >
        <span>{{ $t('remediation.instructionExecute.runJobs') }}</span>
        <v-icon right>
          arrow_right
        </v-icon>
      </v-btn>
      <template v-else>
        <v-icon :color="statusIcon.color">
          {{ statusIcon.name }}
        </v-icon>
        <span class="ml-2">{{ statusIcon.text }}</span>
      </template>
    </v-layout>
    <remediation-instruction-execute-jobs-table :jobs="jobs" />
  </v-layout>
</template>

<script>
import {
  isInstructionExecutionCompleted,
  isInstructionExecutionFinished,
} from '@/helpers/entities/remediation/instruction-execution/form';

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
    executed: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    statusIcon() {
      if (this.isInstructionExecutionSucceeded) {
        return {
          name: 'check_circle',
          color: 'primary',
          text: this.$t('remediation.instructionExecute.jobs.instructionComplete'),
        };
      }

      return {
        name: 'error',
        color: 'error',
        text: this.$t('remediation.instructionExecute.jobs.instructionFailed'),
      };
    },

    isInstructionExecutionFinished() {
      return isInstructionExecutionFinished(this.instructionExecution);
    },

    isInstructionExecutionSucceeded() {
      return isInstructionExecutionCompleted(this.instructionExecution);
    },
  },
};
</script>
