<template lang="pug">
  tr
    td.pa-0 {{ job.name }}
    td.text-xs-center
      span(v-if="!isCancelledJob") {{ job.started_at | date('long', '-') }}
      span(v-else) -
    progress-cell.text-xs-center(:pending="shownLaunchedPendingJob")
      span.error--text(v-if="hasJobsInQueue") {{ queueNumberText }}
      span(v-else) {{ job.launched_at | date('long', '-') }}
    progress-cell.text-xs-center(:pending="shownCompletedPendingJob")
      v-layout(v-if="isFailedJob", row, align-center, justify-center)
        span.error--text {{ $t('common.failed') }}
        c-help-icon.ml-1.cursor-pointer(:text="job.fail_reason", color="error", size="20", top)
      span(v-else) {{ job.completed_at | date('long', '-') }}
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

    hasStatusMessage() {
      return this.job.output || this.job.fail_reason;
    },

    hasJobsInQueue() {
      return this.job.queue_number > 0;
    },

    queueNumberText() {
      return this.$t('remediationInstructionExecute.queueNumber', {
        number: this.job.queue_number,
        name: this.job.name,
      });
    },
  },
};
</script>
