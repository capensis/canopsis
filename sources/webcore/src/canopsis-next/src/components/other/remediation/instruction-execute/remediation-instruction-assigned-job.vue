<template lang="pug">
  tr
    td.pa-0
      v-btn.primary(
        :disabled="isRunningJob",
        round,
        small,
        block,
        @click="$emit('execute-job', job)"
      ) {{ job.name }}
    td.text-xs-center {{ job.started_at | date('long', true, '-') }}
    td.text-xs-center {{ job.launched_at | date('long', true, '-') }}
    td.text-xs-center {{ job.completed_at | date('long', true, '-') }}
</template>

<script>
import { REMEDIATION_JOB_EXECUTION_STATUSES } from '@/constants';

import entitiesRemediationJobsExecutionsMixin from '@/mixins/entities/remediation/jobs-executions';

export default {
  mixins: [entitiesRemediationJobsExecutionsMixin],
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
  },
};
</script>
