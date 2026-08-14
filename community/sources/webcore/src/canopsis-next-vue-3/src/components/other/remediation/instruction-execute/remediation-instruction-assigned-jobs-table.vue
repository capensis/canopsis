<template>
  <table class="jobs-assigned-table">
    <thead>
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
    </thead>
    <tbody>
      <template v-for="job in jobs" :key="job._id || job.job_id || job.name">
        <remediation-instruction-assigned-jobs-row
          :job="job"
          :expanded="expandedJobIds.includes(job._id)"
          :executable="executable"
          :cancelable="cancelable"
          @execute-job="$emit('execute-job', $event)"
          @cancel-job-execution="$emit('cancel-job-execution', $event)"
          @update:expanded="toggleExpanded(job._id, $event)"
        />
        <tr v-if="expandedJobIds.includes(job._id)">
          <td colspan="4">
            <remediation-instruction-assigned-jobs-expand-panel :job="job" />
          </td>
        </tr>
      </template>
    </tbody>
  </table>
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

    const toggleExpanded = (id, expanded) => {
      expandedJobIds.value = expanded
        ? [...expandedJobIds.value, id]
        : expandedJobIds.value.filter(jobId => jobId !== id);
    };

    return {
      t,
      expandedJobIds,
      toggleExpanded,
    };
  },
};
</script>

<style lang="scss" scoped>
.jobs-assigned-table {
  width: 100%;
  border-collapse: collapse;
}
</style>
