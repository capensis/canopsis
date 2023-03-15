<template lang="pug">
  v-data-table.jobs-assigned-table(
    :items="jobs",
    :headers-length="4",
    item-key="_id",
    expand,
    hide-actions
  )
    template(#headers="")
      tr
        th
        th.text-xs-center.pre-line {{ $t('remediation.instructionExecute.jobs.startedAt') }}
        th.text-xs-center.pre-line {{ $t('remediation.instructionExecute.jobs.launchedAt') }}
        th.text-xs-center.pre-line {{ $t('remediation.instructionExecute.jobs.completedAt') }}
    template(#items="row")
      remediation-instruction-assigned-jobs-row(
        :job="row.item",
        :expanded.sync="row.expanded",
        :executable="executable",
        :cancelable="cancelable",
        @execute-job="$emit('execute-job', $event)",
        @cancel-job-execution="$emit('cancel-job-execution', $event)"
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
    executable: {
      type: Boolean,
      default: false,
    },
    cancelable: {
      type: Boolean,
      default: false,
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
  .jobs-assigned-table {
    tr {
      border-bottom: none !important;
    }

    tbody tr:hover {
      background: unset !important;
    }
  }
</style>
