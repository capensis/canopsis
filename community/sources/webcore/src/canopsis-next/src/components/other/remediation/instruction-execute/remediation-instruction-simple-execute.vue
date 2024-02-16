<template>
  <v-layout column>
    <v-layout
      class="mb-4"
      align-center
    >
      <span class="text-subtitle-1 mr-5">{{ $t('remediation.instructionExecute.jobs.title') }}</span>
      <v-btn
        class="primary ma-0"
        v-if="!isJobsFinished"
        :loading="executed"
        @click="$emit('run:jobs')"
      >
        <span>{{ $t('remediation.instructionExecute.runJobs') }}</span>
        <v-icon right>
          arrow_right
        </v-icon>
      </v-btn>
      <template v-if="isJobsFinished">
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
import { isJobExecutionSucceeded, isJobFinished } from '@/helpers/entities/remediation/job/form';

import RemediationInstructionExecuteJobsTable from './remediation-instruction-assigned-jobs-table.vue';

export default {
  components: { RemediationInstructionExecuteJobsTable },
  props: {
    jobs: {
      type: Array,
      required: true,
    },
    executed: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    statusIcon() {
      if (this.isJobsSucceeded) {
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

    isJobsFinished() {
      return this.jobs.every(isJobFinished);
    },

    isJobsSucceeded() {
      return this.jobs.every(isJobExecutionSucceeded);
    },
  },
};
</script>
