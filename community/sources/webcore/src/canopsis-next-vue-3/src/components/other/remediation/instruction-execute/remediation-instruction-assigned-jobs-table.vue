<template>
  <v-data-table
    v-model:expanded="expandedJobIds"
    :headers="headers"
    :items="jobs"
    class="jobs-assigned-table"
    item-value="_id"
    hide-default-footer
    hide-default-header
    show-expand
  >
    <template #headers>
      <tr>
        <th />
        <th class="text-center pre-line">
          {{ t('remediation.instructionExecute.jobs.startedAt') }}
        </th>
        <th class="text-center pre-line">
          {{ t('remediation.instructionExecute.jobs.launchedAt') }}
        </th>
        <th class="text-center pre-line">
          {{ t('remediation.instructionExecute.jobs.completedAt') }}
        </th>
      </tr>
    </template>

    <template #item="{ item }">
      <remediation-instruction-assigned-jobs-row
        :job="getRawJob(item)"
        :expanded="isJobExpanded(item)"
        :executable="executable"
        :cancelable="cancelable"
        @execute-job="$emit('execute-job', $event)"
        @cancel-job-execution="$emit('cancel-job-execution', $event)"
        @update:expanded="toggleExpanded(item, $event)"
      />
    </template>

    <template #expanded-row="{ item, columns }">
      <tr>
        <td :colspan="columns.length" class="jobs-assigned-table__expand-cell">
          <remediation-instruction-assigned-jobs-expand-panel :job="getRawJob(item)" />
        </td>
      </tr>
    </template>
  </v-data-table>
</template>

<script>
import { ref } from 'vue';

import { useI18n } from '@/hooks/i18n';

import RemediationInstructionAssignedJobsExpandPanel from './partials/remediation-instruction-assigned-jobs-expand-panel.vue';
import RemediationInstructionAssignedJobsRow from './partials/remediation-instruction-assigned-jobs-row.vue';

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

  emits: ['execute-job', 'cancel-job-execution'],

  setup() {
    const { t } = useI18n();
    const expandedJobIds = ref([]);
    const headers = [
      { title: '', key: 'name', sortable: false },
      { title: '', key: 'startedAt', sortable: false },
      { title: '', key: 'launchedAt', sortable: false },
      { title: '', key: 'completedAt', sortable: false },
    ];

    const getRawJob = job => job?.raw ?? job;
    const getJobId = (job) => {
      const rawJob = getRawJob(job);

      return rawJob._id || rawJob.job_id || rawJob.name;
    };
    const isJobExpanded = job => expandedJobIds.value.includes(getJobId(job));

    const toggleExpanded = (job, expanded) => {
      const id = getJobId(job);

      expandedJobIds.value = expanded
        ? [...expandedJobIds.value, id]
        : expandedJobIds.value.filter(jobId => jobId !== id);
    };

    return {
      t,
      headers,
      expandedJobIds,
      getRawJob,
      isJobExpanded,
      toggleExpanded,
    };
  },
};
</script>

<style lang="scss" scoped>
.jobs-assigned-table {
  :deep(tr) {
    border-bottom: none !important;
  }

  :deep(tbody tr:hover) {
    background: unset !important;
  }

  &__expand-cell {
    padding: 0;
  }
}
</style>
