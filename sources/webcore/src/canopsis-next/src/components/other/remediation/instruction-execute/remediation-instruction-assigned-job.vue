<template lang="pug">
  tr
    td.pa-0
      v-tooltip(right, :disabled="!isFailedJob")
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
    td.text-xs-center {{ job.launched_at | date('long', true, '-') }}
    td.text-xs-center {{ job.completed_at | date('long', true, '-') }}
</template>

<script>
import { REMEDIATION_JOB_EXECUTION_STATUSES } from '@/constants';

export default {
  props: {
    job: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isFailedJob() {
      return !!this.job.fail_reason;
    },

    isRunningJob() {
      return this.job.status === REMEDIATION_JOB_EXECUTION_STATUSES.running;
    },
  },
};
</script>
