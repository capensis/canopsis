<template>
  <c-advanced-data-table
    :items="steps"
    :headers="headers"
    hide-actions
  >
    <template #item="props">
      <remediation-instruction-executions-expand-panel-row :step="props.item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

import RemediationInstructionExecutionsExpandPanelRow from './remediation-instruction-executions-expand-panel-row.vue';

export default {
  components: { RemediationInstructionExecutionsExpandPanelRow },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const { t } = useI18n();
    const headers = computed(() => [
      { text: '', value: 'icon', sortable: false },
      { text: t('common.step'), value: 'name', sortable: false },
      { text: t('common.output'), value: 'fail_reason', sortable: false },
      { text: t('remediation.instructionExecute.jobs.completedAt'), value: 'completed_at', sortable: false },
    ]);

    return { headers };
  },
};
</script>
