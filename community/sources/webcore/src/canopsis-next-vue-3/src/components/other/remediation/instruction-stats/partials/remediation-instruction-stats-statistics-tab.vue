<template>
  <c-advanced-data-table
    :items="changes"
    :headers="headers"
    :loading="pending"
    :options="options"
    :total-items="totalItems"
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #modified_on="{ item }">
      <span>{{ dateToString(item.modified_on) }}</span>
    </template>
    <template #avg_complete_time="{ item }">
      <span v-if="item.execution_count">{{ durationToString(item.avg_complete_time) }}</span>
      <span v-else>{{ t('common.notAvailable') }}</span>
    </template>
    <template #avg_alarm_ok_timeout="{ item }">
      <span v-if="item.avg_alarm_ok_timeout">{{
        durationToString(item.avg_alarm_ok_timeout)
      }}</span>
      <span v-else>{{ t('common.notAvailable') }}</span>
    </template>
    <template #avg_successful="{ item }">
      <span>{{ item.avg_successful }}%</span>
    </template>
    <template #avg_successful_state_ok="{ item }">
      <span>{{ item.avg_successful_state_ok }}%</span>
    </template>
    <template #alarm_states="{ item }">
      <affect-alarm-states v-if="item.execution_count" :alarm-states="item.alarm_states" />
      <template v-else>
        -
      </template>
    </template>
    <template #ok_alarm_states="{ item }">
      <c-state-count-changes-chip v-if="item.execution_count">
        {{ item.ok_alarm_states }}
      </c-state-count-changes-chip>
      <template v-else>
        -
      </template>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed, onMounted, ref, watch } from 'vue';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { useI18n } from '@/hooks/i18n';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';
import { useRemediationInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

import AffectAlarmStates from './affect-alarm-states.vue';

export default {
  components: { AffectAlarmStates },

  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { fetchRemediationInstructionStatsChangesWithoutStore } =
      useRemediationInstructionStats();
    const totalItems = ref(0);

    const {
      data: changes,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      fetchListHandler: ({ params }) => fetchRemediationInstructionStatsChangesWithoutStore({
          id: props.remediationInstruction._id,
          params,
        }),
    });

    watch(meta, (value) => {
      totalItems.value = value?.total_count ?? 0;
    });

    const headers = computed(() => [
      { text: t('remediation.instructionStat.modifiedOn'), value: 'modified_on', sortable: false },
      {
        text: t('remediation.instructionStat.averageCompletionTime'),
        value: 'avg_complete_time',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.averageAlarmOkTimeout'),
        value: 'avg_alarm_ok_timeout',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.executionCount'),
        value: 'execution_count',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.averageSuccessfulAll'),
        value: 'avg_successful',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.averageSuccessfulOk'),
        value: 'avg_successful_state_ok',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.alarmStates'),
        value: 'alarm_states',
        sortable: false,
      },
      {
        text: t('remediation.instructionStat.okAlarmStates'),
        value: 'ok_alarm_states',
        sortable: false,
      },
    ]);
    const dateToString = value => convertDateToString(value);
    const durationToString = value => convertDurationToString(value);

    onMounted(fetchList);

    return {
      t,
      changes,
      pending,
      options,
      totalItems,
      headers,
      updateOptions,
      dateToString,
      durationToString,
    };
  },
};
</script>
