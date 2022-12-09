<template lang="pug">
  v-sheet.px-3.py-2(color="grey lighten-3")
    div {{ $t('remediationInstructionExecute.jobs.failedReason') }}:&nbsp;
      span.pre-wrap(v-html="job.fail_reason")
    div {{ $t('remediationInstructionExecute.jobs.output') }}:&nbsp;
      span.pre-line(v-html="output")
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { isJobExecutionCancelled, isJobExecutionRunning } from '@/helpers/forms/remediation-job';

import ProgressCell from '@/components/common/table/progress-cell.vue';

const { mapGetters } = createNamespacedHelpers('remediationJobExecution');

export default {
  components: { ProgressCell },
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    ...mapGetters(['getOutputById']),

    output() {
      return this.getOutputById(this.job._id);
    },

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
  },
};
</script>
