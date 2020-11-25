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
    td.text-xs-center {{ job.started_at | date('long', true, '-') }}
    progress-cell.text-xs-center(:pending="!job.completed_at && isStartedJob")
      span {{ job.launched_at | date('long', true, '-') }}
    progress-cell.text-xs-center(:pending="!job.completed_at && isStartedJob")
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
    isStartedJob() {
      return !!this.job.started_at;
    },

    isFailedJob() {
      return !!this.job.fail_reason;
    },

    isRunningJob() {
      return this.job.status === REMEDIATION_JOB_EXECUTION_STATUSES.running;
    },
  },
};
</script>
