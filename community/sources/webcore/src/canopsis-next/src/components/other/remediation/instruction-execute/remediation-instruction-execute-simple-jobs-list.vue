<template lang="pug">
  table.instruction-execute-jobs-table
    thead
      tr
        th
        th.grey--text.lighten-3 {{ $t('remediation.instructionExecute.jobs.startedAt') }}
        th.grey--text.lighten-3 {{ $t('remediation.instructionExecute.jobs.launchedAt') }}
        th.grey--text.lighten-3 {{ $t('remediation.instructionExecute.jobs.completedAt') }}
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
