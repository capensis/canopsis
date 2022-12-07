<template lang="pug">
  div.jobs-assigned
    v-layout(row)
      span.subheading {{ $t('remediation.instructionExecute.jobs.title') }}
    v-layout(column)
      v-data-table(:items="jobs", hide-actions)
        template(#headers="")
          tr
            th
            th.text-xs-center.pre-line {{ $t('remediation.instructionExecute.jobs.startedAt') }}
            th.text-xs-center.pre-line {{ $t('remediation.instructionExecute.jobs.launchedAt') }}
            th.text-xs-center.pre-line {{ $t('remediation.instructionExecute.jobs.completedAt') }}
        template(#items="{ item }")
          remediation-instruction-assigned-job(
            :job="item",
            :key="item.job_id",
            @execute-job="$listeners['execute-job']",
            @cancel-job-execution="$listeners['cancel-job-execution']"
          )
</template>

<script>
import RemediationInstructionAssignedJob from './partials/remediation-instruction-assigned-job.vue';

export default {
  components: {
    RemediationInstructionAssignedJob,
  },
  props: {
    jobs: {
      type: Array,
      default: () => [],
    },
  },
  methods: {
    async executeJob(job) {
      this.$emit('execute-job', { job });
    },
  },
};
</script>

<style lang="scss" scoped>
  .jobs-assigned {
    tr {
      border-bottom: none !important;
    }

    tbody tr:hover {
      background: unset !important;
    }
  }
</style>
