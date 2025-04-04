<template>
  <c-advanced-data-table
    :items="preparedRemediationInstructionExecutions"
    :headers="headers"
    :loading="pending"
    :options="options"
    :total-items="totalItems"
    :is-expandable-item="isExpandedItem"
    class="instruction-alarms-timeline-table"
    search
    expand
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar="">
      <v-layout align-center>
        <c-enabled-field
          v-model="showFailed"
          :label="$t('remediation.instructionStat.showFailedExecutions')"
          hide-details
        />
      </v-layout>
    </template>
    <template #alarm_display_name="{ item }">
      {{ item.alarm_display_name }}
    </template>
    <template #result="{ item }">
      <c-enabled
        v-if="item.alarm_id"
        :value="item.completed"
        :icon-success="item.result_icon"
        :icon-failed="item.result_icon"
      />
    </template>
    <template #result_alarm_state="{ item }">
      <c-alarm-state-chip v-if="item.alarm_id" :value="item.result_alarm_state" />
    </template>
    <template #started_at="{ item }">
      {{ item.started_at }}
    </template>
    <template #completed_at="{ item }">
      {{ item.completed_at }}
    </template>
    <template #duration="{ item }">
      {{ item.duration }}
    </template>
    <template #alarm_ok_at="{ item }">
      {{ item.alarm_ok_at }}
    </template>
    <template #timeout_after_execution="{ item }">
      {{ item.timeout_after_execution }}
    </template>
    <template #alarm_ok_before_completed="{ item }">
      {{ item.alarm_ok_before_completed }}
    </template>
    <template #timeline="{ item }">
      <span
        v-if="!item.alarm_id"
        class="text--secondary"
      >{{ $t('remediation.instructionStat.instructionChanged') }}</span>
      <alarm-horizontal-timeline
        v-else
        :steps="item.alarm_steps"
        class="my-2"
      />
    </template>
    <template #expand="{ item }">
      <remediation-instruction-stats-alarms-timeline-tab-expand-panel
        v-if="item.alarm_id"
        :alarm-id="item.alarm_id"
        :execution-id="item._id"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, ref, watch, onMounted } from 'vue';

import { DATETIME_FORMATS } from '@/constants';

import {
  prepareRemediationInstructionExecutionsForAlarmTimeline,
} from '@/helpers/entities/remediation/instruction-execution/list';
import { getQueryForList } from '@/helpers/entities/shared/query';
import { isInstructionTypeAuto } from '@/helpers/entities/remediation/instruction/form';
import { isInstructionExecutionCompleted } from '@/helpers/entities/remediation/instruction-execution/form';
import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { usePendingWithLocalQuery } from '@/hooks/query/shared';
import { useQueryOptions } from '@/hooks/query/options';
import { useRemdeitionInstructionStatsStore } from '@/hooks/store/modules/remediation-instruction-stats';

import AlarmHorizontalTimeline from '@/components/widgets/alarm/timeline/horizontal-timeline.vue';
import RemediationInstructionStatsAlarmsTimelineTabExpandPanel
  from '@/components/other/remediation/instruction-stats/partials/remediation-instruction-stats-alarms-timeline-tab-expand-panel.vue';

export default {
  components: { RemediationInstructionStatsAlarmsTimelineTabExpandPanel, AlarmHorizontalTimeline },
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
    interval: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const { fetchRemediationInstructionStatsExecutionsListWithoutStore } = useRemdeitionInstructionStatsStore();

    const remediationInstructionExecutions = ref([]);
    const totalItems = ref(0);
    const showFailed = ref(true);

    const preparedRemediationInstructionExecutions = computed(() => (
      remediationInstructionExecutions.value.map((execution = {}) => (!execution.alarm_id
        ? {}
        : ({
          ...execution,

          started_at: convertDateToString(execution.started_at, DATETIME_FORMATS.long, '-'),
          completed_at: convertDateToString(execution.completed_at, DATETIME_FORMATS.long, '-'),
          alarm_ok_at: convertDateToString(execution.alarm_ok_at, DATETIME_FORMATS.long, '-'),
          duration: convertDurationToString(execution.duration) || '-',
          timeout_after_execution: convertDurationToString(execution.timeout_after_execution) || '-',
          completed: isInstructionExecutionCompleted(execution),
          result_icon: isInstructionTypeAuto(execution.instruction_type)
            ? 'assignment'
            : '$vuetify.icons.manual_instruction',

          alarm_ok_before_completed: execution.alarm_ok_before_completed
            ? t('remediation.instructionStat.solveBeforeRemediationEnd')
            : convertDurationToString(execution.alarm_ok_timeout) || '-',
        })))
    ));

    const headers = computed(() => [
      {
        text: t('common.alarmId'),
        value: 'alarm_display_name',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.instructionResult'),
        value: 'result',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.alarmStateAfterTimeout'),
        value: 'result_alarm_state',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.remediationStart'),
        value: 'started_at',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.remediationEnd'),
        value: 'completed_at',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.remediationDuration'),
        value: 'duration',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.okStateDate'),
        value: 'alarm_ok_at',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.timeoutAfterExecution'),
        value: 'timeout_after_execution',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.afterRemediationEnd'),
        value: 'alarm_ok_before_completed',
        sortable: false,
      },
      {
        value: 'timeline',
        sortable: false,
      },
    ]);

    const {
      pending,
      query,
      fetchHandlerWithQuery: fetchList,
      updateQuery,
    } = usePendingWithLocalQuery({
      fetchHandler: async (fetchQuery) => {
        const params = getQueryForList(fetchQuery);

        params.from = props.interval.from;
        params.to = props.interval.to;
        params.show_failed = showFailed.value;

        const {
          data,
          meta,
        } = await fetchRemediationInstructionStatsExecutionsListWithoutStore({
          params,
          id: props.remediationInstruction._id,
        });

        remediationInstructionExecutions.value = prepareRemediationInstructionExecutionsForAlarmTimeline(data);
        totalItems.value = meta.total_count;
      },
    });

    const { options, updateOptions } = useQueryOptions(query, updateQuery);

    const isExpandedItem = item => !!item.alarm_id;

    watch(() => props.interval, fetchList);
    watch(showFailed, fetchList);

    onMounted(fetchList);

    return {
      preparedRemediationInstructionExecutions,
      totalItems,
      showFailed,
      headers,
      pending,
      options,
      updateOptions,
      isExpandedItem,
    };
  },
};
</script>

<style lang="scss">
.instruction-alarms-timeline-table td {
  &, & > * {
    white-space: nowrap;
  }
}
</style>
