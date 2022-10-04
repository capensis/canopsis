<template lang="pug">
  v-layout(column)
    v-layout.mb-4(row, align-center)
      span.subheading.mr-5 {{ $t('remediationInstructionExecute.jobs.title') }}
      v-btn.primary.ma-0(:loading="!jobsStarted")
        span {{ $t('remediationInstructionExecute.runJobs') }}
        v-icon(right) arrow_right
    table.instruction-execute-jobs-table
      thead
        tr
          th
          th.grey--text.lighten-3 {{ $t('remediationInstructionExecute.jobs.startedAt') }}
          th.grey--text.lighten-3 {{ $t('remediationInstructionExecute.jobs.launchedAt') }}
          th.grey--text.lighten-3 {{ $t('remediationInstructionExecute.jobs.completedAt') }}
      tbody
        remediation-instruction-execute-simple-job(v-for="job in jobs", :key="job.job_id", :job="job")
</template>

<script>
import RemediationInstructionExecuteSimpleJob from './partials/remediation-instruction-execute-simple-job.vue';

export default {
  components: { RemediationInstructionExecuteSimpleJob },
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      jobsStarted: false,
    };
  },
};
</script>

<style lang="scss" scoped>
.instruction-execute-jobs-table {
  tbody {
    tr, td {
      height: 40px;
    }
  }
}
</style>
