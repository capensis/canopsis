<template>
  <div class="instruction-stats-alarms-timeline-tab-expand-panel__wrapper pa-3">
    <c-advanced-data-table
      :items="steps"
      :headers="headers"
      :loading="pending"
    >
      <template #status="{ item }">
        <remediation-instruction-stats-alarms-timeline-tab-expand-panel-status-icon
          :type="execution.type"
          :status="item.status"
          :name="item.name"
        />
      </template>
      <template #completed_at="{ item }">
        {{ item.completed_at | date('long', '-') }}
      </template>
    </c-advanced-data-table>
  </div>
</template>

<script>
import { computed, ref, onMounted } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { usePendingHandler } from '@/hooks/query/pending';
import { useRemediationInstructionExecution } from '@/hooks/store/modules/remediation-instruction-execution';

import RemediationInstructionStatsAlarmsTimelineTabExpandPanelStatusIcon
  from './remediation-instruction-stats-alarms-timeline-tab-expand-panel-status-icon.vue';

export default {
  components: { RemediationInstructionStatsAlarmsTimelineTabExpandPanelStatusIcon },
  props: {
    alarmId: {
      type: String,
      required: true,
    },
    executionId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const execution = ref({});

    const { t } = useI18n();

    const {
      fetchAlarmRemediationInstructionExecutionsWithoutStore,
    } = useRemediationInstructionExecution();

    const steps = computed(() => execution.value?.steps ?? []);

    const {
      pending,
      handler: fetchList,
    } = usePendingHandler(async () => {
      const response = await fetchAlarmRemediationInstructionExecutionsWithoutStore({
        alarmId: props.alarmId,
        params: {
          ids: [props.executionId],
        },
      });

      execution.value = response.data?.[0] ?? {};
    });

    const headers = computed(() => [
      {
        text: t('common.step'),
        value: 'name',
        sortable: false,
        width: '300px',
      },
      {
        text: t('common.status'),
        value: 'status',
        sortable: false,
        width: '100px',
      },
      {
        text: t('common.output'),
        value: 'fail_reason',
        sortable: false,
      },
      {
        text: t('remediation.instructionExecute.jobs.completedAt'),
        value: 'completed_at',
        sortable: false,
        width: '200px',
      },
    ]);

    onMounted(fetchList);

    return {
      execution,
      steps,
      pending,
      headers,
    };
  },
};
</script>

<style lang="scss">
.instruction-stats-alarms-timeline-tab-expand-panel__wrapper {
  --expand-panel-width: 94vw;

  width: var(--expand-panel-width);

  table {
    table-layout: fixed;
    width: 100%;

    * {
      word-break: break-all;
      word-wrap: break-word;
    }
  }
}
</style>
