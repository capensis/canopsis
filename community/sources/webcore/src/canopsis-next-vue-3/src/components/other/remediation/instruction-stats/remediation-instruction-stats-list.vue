<template>
  <div class="c-remediation-instruction-stats">
    <c-advanced-data-table
      :headers="headers"
      :items="remediationInstructionStats"
      :loading="pending"
      :total-items="totalItems"
      :options="options"
      table-class="c-remediation-instruction-stats__table"
      expand
      search
      advanced-pagination
      @update:options="$emit('update:options', $event)"
    >
      <template #toolbar>
        <c-layout align="center">
          <c-quick-date-interval-field
            :interval="options.interval"
            :accumulated-before="accumulatedBefore"
            @input="updateInterval"
          />
        </c-layout>
      </template>
      <template #type="{ item }">
        {{ t(`remediation.instruction.types.${item.type}`) }}
      </template>
      <template #last_executed_on="{ item }">
        {{ dateToString(item.last_executed_on) }}
      </template>
      <template #last_modified="{ item }">
        {{ dateToString(item.last_modified) }}
      </template>
      <template #avg_complete_time="{ item }">
        {{ durationToString(item.avg_complete_time) }}
      </template>
      <template #avg_alarm_ok_timeout="{ item }">
        {{ durationToString(item.avg_alarm_ok_timeout) }}
      </template>
      <template #avg_successful="{ item }">
        {{ item.avg_successful }}% / {{ item.avg_successful_state_ok }}%
      </template>
      <template #alarm_states="{ item }">
        <affect-alarm-states :alarm-states="item.alarm_states" />
      </template>
      <template #ok_alarm_states="{ item }">
        <c-state-count-changes-chip>{{ item.ok_alarm_states }}</c-state-count-changes-chip>
      </template>
      <template #rating="{ item }">
        <rating-field :model-value="item.rating" readonly />
      </template>
      <template #actions="{ item }">
        <c-layout justify="end">
          <c-action-btn
            v-if="item.rate_notify"
            :tooltip="t('remediation.instructionStat.actions.needRate')"
            icon="notification_important"
            color="error"
            @click="$emit('rate', item)"
          />
          <c-action-btn
            :tooltip="t('remediation.instructionStat.actions.rate')"
            :disabled-button="!item.ratable"
            icon="thumbs_up_down"
            @click="$emit('rate', item)"
          />
        </c-layout>
      </template>
      <template #expand="{ item }">
        <remediation-instruction-stats-list-expand-panel
          :interval="interval"
          :remediation-instruction-stats-item="item"
        />
      </template>
    </c-advanced-data-table>
  </div>
</template>

<script>
import { computed } from 'vue';

import { USER_PERMISSIONS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';
import { convertDurationToString } from '@/helpers/date/duration';

import { useCRUDPermissions } from '@/hooks/auth';
import { useI18n } from '@/hooks/i18n';

import RatingField from '@/components/forms/fields/rating-field.vue';

import AffectAlarmStates from './partials/affect-alarm-states.vue';
import RemediationInstructionStatsListExpandPanel from './partials/remediation-instruction-stats-list-expand-panel.vue';

export default {
  components: {
    AffectAlarmStates,
    RatingField,
    RemediationInstructionStatsListExpandPanel,
  },

  props: {
    remediationInstructionStats: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    accumulatedBefore: {
      type: Number,
      required: false,
    },
    options: {
      type: Object,
      required: true,
    },
    interval: {
      type: Object,
      required: true,
    },
  },

  emits: ['update:options', 'rate'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { hasCreateAccess } = useCRUDPermissions(
      USER_PERMISSIONS.technical.exploitation.remediationInstruction,
    );

    const headers = computed(() => [
        { text: t('common.name'), value: 'name', width: 150 },
        hasCreateAccess.value && { text: t('common.type'), value: 'type', width: 100 },
        {
          text: t('remediation.instructionStat.lastExecutedOn'),
          value: 'last_executed_on',
          width: 160,
        },
        { text: t('common.lastModifiedOn'), value: 'last_modified', width: 160 },
        {
          text: t('remediation.instructionStat.averageCompletionTime'),
          value: 'avg_complete_time',
          width: 140,
        },
        {
          text: t('remediation.instructionStat.averageAlarmOkTimeout'),
          value: 'avg_alarm_ok_timeout',
          width: 100,
        },
        {
          text: t('remediation.instructionStat.executionCount'),
          value: 'execution_count',
          width: 100,
        },
        {
          text: t('remediation.instructionStat.averageSuccessfulAllOk'),
          value: 'avg_successful',
          width: 100,
        },
        {
          text: t('remediation.instructionStat.alarmStates'),
          value: 'alarm_states',
          width: 300,
        },
        {
          text: t('remediation.instructionStat.okAlarmStates'),
          value: 'ok_alarm_states',
          width: 90,
        },
        { text: t('common.rating'), value: 'rating', sortable: false, width: 230 },
        { text: t('common.actionsLabel'), value: 'actions', sortable: false, width: 80 },
      ].filter(Boolean),
    );

    const updateInterval = interval => emit('update:options', {
        ...props.options,
        interval,
      });
    const dateToString = value => convertDateToString(value);
    const durationToString = value => convertDurationToString(value) || '-';

    return {
      t,
      headers,
      updateInterval,
      dateToString,
      durationToString,
    };
  },
};
</script>

<style lang="scss">
.c-remediation-instruction-stats__table {
  table {
    table-layout: fixed !important;
  }
}
</style>
