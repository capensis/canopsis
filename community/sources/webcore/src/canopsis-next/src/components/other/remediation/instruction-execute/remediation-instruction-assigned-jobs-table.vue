<template>
  <v-data-table
    :items="jobs"
    :headers="headers"
    class="jobs-assigned-table"
    item-key="_id"
    show-expand
    hide-default-footer
    hide-default-header
  >
    <template #header="">
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
    <template #item="row">
      <remediation-instruction-assigned-jobs-row
        :job="row.item"
        :expanded="row.isExpanded"
        :executable="executable"
        :cancelable="cancelable"
        @execute-job="$emit('execute-job', $event)"
        @cancel-job-execution="$emit('cancel-job-execution', $event)"
        @update:expanded="row.expand"
      />
    </template>
    <template #expanded-item="{ item }">
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
  setup() {
    const headers = [
      { value: 'name' },
      { value: 'name' },
      { value: 'startedAt' },
      { value: 'launchedAt' },
      { value: 'completedAt' },
    ];

    return { headers };
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
