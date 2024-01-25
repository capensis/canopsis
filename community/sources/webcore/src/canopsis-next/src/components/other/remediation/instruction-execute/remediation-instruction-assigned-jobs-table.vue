<template>
  <v-data-table
    class="jobs-assigned-table"
    :items="jobs"
    :headers-length="4"
    item-key="_id"
    show-expand
    hide-default-footer
  >
    <template #headers="">
      <tr>
        <th />
        <th class="text-center pre-line">
          {{ $t('remediation.instructionExecute.jobs.startedAt') }}
        </th>
        <th class="text-center pre-line">
          {{ $t('remediation.instructionExecute.jobs.launchedAt') }}
        </th>
        <th class="text-center pre-line">
          {{ $t('remediation.instructionExecute.jobs.completedAt') }}
        </th>
      </tr>
    </template>
    <template #items="row">
      <remediation-instruction-assigned-jobs-row
        :job="row.item"
        :expanded.sync="row.expanded"
        :executable="executable"
        :cancelable="cancelable"
        @execute-job="$emit('execute-job', $event)"
        @cancel-job-execution="$emit('cancel-job-execution', $event)"
      />
    </template>
    <template #expand="{ item }">
      <remediation-instruction-assigned-jobs-expand-panel :job="item" />
    </template>
  </v-data-table>
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
