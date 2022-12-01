<template lang="pug">
  tr
    td.pa-0
      v-layout(row, align-center)
        v-tooltip(v-if="isFailedJob || isCompletedJob", :disabled="!hasStatusMessage", max-width="400", right)
          template(#activator="{ on }")
            v-icon.mr-1(v-on="on", :color="isFailedJob ? 'error' : 'success'") {{ statusIcon }}
          div(v-if="job.fail_reason")
            span {{ $t('remediationInstructionExecute.jobs.failedReason') }}:&nbsp;
            span.pre-wrap(v-html="job.fail_reason")
          div(v-if="job.output")
            span {{ $t('remediationInstructionExecute.jobs.output') }}:&nbsp;
            span.pre-wrap(v-html="job.output")
        v-btn.primary(
          :disabled="isRunningJob || isFailedJob",
          round,
          small,
          block,
          @click="$emit('execute-job', job)"
        ) {{ job.name }}
        v-tooltip(v-show="isRunningJob && hasJobsInQueue", top)
          template(#activator="{ on }")
            v-btn.error.ml-2(
              v-on="on",
              round,
              small,
              block,
              @click="$emit('cancel-job-execution', job)"
            ) {{ $t('common.cancel') }}
          span {{ queueNumberTooltip }}
    td.text-xs-center
      span(v-if="!isCancelledJob") {{ job.started_at | date('long', '-') }}
      span(v-else) -
    progress-cell.text-xs-center(:pending="shownLaunchedPendingJob")
      span {{ job.launched_at | date('long', '-') }}
    progress-cell.text-xs-center(:pending="shownCompletedPendingJob")
      span {{ job.completed_at | date('long', '-') }}
</template>

<script>
import { isJobExecutionCancelled, isJobExecutionRunning } from '@/helpers/forms/remediation-job';

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
      return isJobExecutionRunning(this.job);
    },

    isCancelledJob() {
      return isJobExecutionCancelled(this.job);
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

    statusIcon() {
      return this.isFailedJob ? 'cancel' : 'check';
    },

    queueNumberTooltip() {
      return this.$t('remediationInstructionExecute.queueNumber', {
        number: this.job.queue_number,
        name: this.job.name,
      });
    },
  },
};
</script>
