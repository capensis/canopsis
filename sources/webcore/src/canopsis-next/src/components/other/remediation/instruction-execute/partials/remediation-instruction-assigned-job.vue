<template lang="pug">
  tr
    td.pa-0
      v-layout(row, align-center)
        v-tooltip(v-if="isFailedJob || isCompletedJob", :disabled="!isFailedJob", right)
          v-icon.mr-1(slot="activator", :color="isFailedJob ? 'error' : 'success'") {{ statusIcon }}
          span {{ job.fail_reason }}
        v-btn.primary(
          :disabled="isRunningJob || isFailedJob",
          round,
          small,
          block,
          @click="$emit('execute-job', job)"
        ) {{ job.name }}
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

import ProgressCell from '@/components/common/table/progress-cell.vue';

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

    statusIcon() {
      return this.isFailedJob ? 'cancel' : 'check';
    },
  },
};
</script>
