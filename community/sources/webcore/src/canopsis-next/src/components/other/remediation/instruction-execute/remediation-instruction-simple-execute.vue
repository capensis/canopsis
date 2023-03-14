<template lang="pug">
  v-layout(column)
    v-layout.mb-4(row, align-center)
      span.subheading.mr-5 {{ $t('remediation.instructionExecute.jobs.title') }}
      v-btn.primary.ma-0(v-if="!isJobsFinished", :loading="executed", @click="$emit('run:jobs')")
        span {{ $t('remediation.instructionExecute.runJobs') }}
        v-icon(right) arrow_right
      template(v-if="isJobsFinished")
        v-icon(:color="statusIcon.color") {{ statusIcon.name }}
        span.ml-2 {{ statusIcon.text }}
    remediation-instruction-execute-jobs-table(:jobs="jobs")
</template>

<script>
import { isJobExecutionSucceeded, isJobFinished } from '@/helpers/forms/remediation-job';

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
