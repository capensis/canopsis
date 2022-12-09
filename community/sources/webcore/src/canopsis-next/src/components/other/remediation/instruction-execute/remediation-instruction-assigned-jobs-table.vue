<template lang="pug">
  div.jobs-assigned
    v-layout(row)
      span.subheading {{ $t('remediationInstructionExecute.jobs.title') }}
    v-layout(column)
      v-data-table(
        :items="jobs",
        :headers-length="4",
        item-key="_id",
        expand,
        hide-actions
      )
        template(#headers="")
          tr
            th
            th.text-xs-center.pre-line {{ $t('remediationInstructionExecute.jobs.startedAt') }}
            th.text-xs-center.pre-line {{ $t('remediationInstructionExecute.jobs.launchedAt') }}
            th.text-xs-center.pre-line {{ $t('remediationInstructionExecute.jobs.completedAt') }}
        template(#items="row")
          remediation-instruction-assigned-jobs-row(
            :job="row.item",
            :expanded.sync="row.expanded",
            @execute-job="$listeners['execute-job']",
            @cancel-job-execution="$listeners['cancel-job-execution']"
          )
        template(#expand="{ item }")
          remediation-instruction-assigned-jobs-expand-panel(:job="item")
</template>

<script>
import RemediationInstructionAssignedJobsRow from './partials/remediation-instruction-assigned-jobs-row.vue';
import RemediationInstructionAssignedJobsExpandPanel from './partials/remediation-instruction-assigned-jobs-expand-panel.vue';

export default {
  components: {
    RemediationInstructionAssignedJobsRow,
    RemediationInstructionAssignedJobsExpandPanel,
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
