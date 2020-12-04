<template lang="pug">
  tr
    td.pa-0
      v-tooltip(:disabled="!isFailedJob", right)
        v-btn.primary(
          :disabled="isRunningJob || isFailedJob",
          slot="activator",
          round,
          small,
          block,
          @click="$emit('execute-job', job)"
        ) {{ job.name }}
        span {{ job.fail_reason }}
    td.text-xs-center
      span(v-if="!isCancelledJob") {{ job.started_at | date('long', true, '-') }}
      span(v-else) -
    progress-cell.text-xs-center(:pending="shownLaunchedPendingJob")
      span {{ job.launched_at | date('long', true, '-') }}
    progress-cell.text-xs-center(:pending="shownCompletedPendingJob")
      span {{ job.completed_at | date('long', true, '-') }}
</template>

<script>
import { REMEDIATION_JOB_EXECUTION_STATUSES } from '@/constants';

import ProgressCell from '@/components/other/table/progress-cell.vue';

export default {
  components: { ProgressCell },
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isRunningJob() {
      return this.job.status === REMEDIATION_JOB_EXECUTION_STATUSES.running;
    },

    isCancelledJob() {
      return this.job.status === REMEDIATION_JOB_EXECUTION_STATUSES.canceled;
    },

    isStartedJob() {
      return !!this.job.started_at;
    },

    isLaunchedJob() {
      return !!this.job.launched_at;
    },

    isCompletedJob() {
      return !!this.job.completed_at;
    },

    isFailedJob() {
      return !!this.job.fail_reason;
    },

    shownLaunchedPendingJob() {
      return !this.isCancelledJob
        && !this.isFailedJob
        && !this.isLaunchedJob
        && this.isStartedJob;
    },

    shownCompletedPendingJob() {
      return !this.isCancelledJob
        && !this.isFailedJob
        && !this.isCompletedJob
        && this.isStartedJob
        && this.isLaunchedJob;
    },
  },
};
</script>
