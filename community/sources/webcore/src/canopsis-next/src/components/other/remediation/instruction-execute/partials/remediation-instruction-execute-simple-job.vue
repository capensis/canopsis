<template lang="pug">
  tr
    td.pa-0 {{ job.name }}
    td(v-if="rowMessage", colspan="3")
      div.error--text.text-xs-center {{ rowMessage }}
    template(v-else)
      td.text-xs-center
        span {{ job.started_at | date('long', '-') }}
      progress-cell.text-xs-center(:pending="shownLaunchedPendingJob")
        span {{ job.launched_at | date('long', '-') }}
      progress-cell.text-xs-center(:pending="shownCompletedPendingJob")
        v-layout(v-if="isFailedJob", row, align-center, justify-center)
          span.error--text {{ $t('common.failed') }}
          c-help-icon.ml-1.cursor-pointer(:text="job.fail_reason", color="error", size="20", top)
        span(v-else) {{ job.completed_at | date('long', '-') }}
</template>

<script>
import { isJobExecutionCancelled } from '@/helpers/forms/remediation-job';

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
      return isJobExecutionCancelled(this.job);
    },

    isStartedJob() {
      return !!this.job.started_at;
    },

    isLaunchedJob() {
      return !!this.job.launched_at;
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
        && !this.job.completed_at
        && this.isStartedJob
        && this.isLaunchedJob;
    },

    rowMessage() {
      if (this.job.queue_number > 0) {
        return this.$t('remediation.instructionExecute.queueNumber', {
          number: this.job.queue_number,
          name: this.job.name,
        });
      }

      if (this.isCancelledJob) {
        return this.$t('remediation.instructionExecute.jobs.stopped');
      }

      return '';
    },
  },
};
</script>
