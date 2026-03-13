<template>
  <c-advanced-data-table
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    search
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <v-flex xs4>
        <c-select-field
          :value="options.status"
          :label="$t('jobs.filterByStatus')"
          :items="statusItems"
          item-value="value"
          item-text="text"
          clearable
          hide-details
          @input="updateStatus"
        />
      </v-flex>
    </template>

    <template #ruleName="{ item }">
      {{ getRuleName(item) }}
    </template>
    <template #authTokenName="{ item }">
      {{ getAuthTokenName(item) }}
    </template>
    <template #ticketSystemName="{ item }">
      {{ item.ticket_system_name || '-' }}
    </template>
    <template #ruleType="{ item }">
      {{ ruleTypeLabel(item.rule_type) }}
    </template>
    <template #active="{ item }">
      <v-tooltip right>
        <template #activator="{ on }">
          <v-icon
            :color="statusColor(item.status)"
            v-on="on"
          >
            {{ statusIcon(item.status) }}
          </v-icon>
        </template>
        <span>{{ activeStateLabel(item.status) }}</span>
      </v-tooltip>
    </template>
    <template #status="{ item }">
      {{ runStatusLabel(item.last_run_status) }}
    </template>
    <template #startDate="{ item }">
      {{ item.created_at | date }}
    </template>
    <template #finishDate="{ item }">
      <template v-if="shouldShowFinishDate(item)">
        {{ item.checked_at | date }}
      </template>
      <template v-else>
        -
      </template>
    </template>
    <template #failReason="{ item }">
      <template v-if="shouldShowFailReason(item)">
        {{ item.fail_reason || '-' }}
      </template>
      <template v-else>
        -
      </template>
    </template>
    <template #expirationDate="{ item }">
      <template v-if="shouldShowExpirationDate(item)">
        {{ item.expiration_at | date }}
      </template>
      <template v-else>
        -
      </template>
    </template>
    <template #expand="{ item }">
      <ticket-status-jobs-details-expand-panel :item="item" />
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          v-if="canStart(item)"
          :tooltip="$t('jobs.actions.start')"
          icon="play_arrow"
          color="primary"
          @click="emitAction('start', item)"
        />
        <c-action-btn
          v-if="canStop(item)"
          :tooltip="$t('jobs.actions.stop')"
          icon="stop"
          color="error"
          @click="emitAction('stop', item)"
        />
        <c-action-btn
          v-if="canResume(item)"
          :tooltip="$t('jobs.actions.resume')"
          icon="play_arrow"
          color="primary"
          @click="emitAction('resume', item)"
        />
        <c-action-btn
          v-if="canPause(item)"
          :tooltip="$t('jobs.actions.pause')"
          icon="pause"
          color="warning"
          @click="emitAction('pause', item)"
        />
        <c-action-btn
          :tooltip="$t('jobs.actions.edit')"
          icon="edit"
          color="#6A6A6A"
          @click="emitAction('edit', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { computed } from 'vue';

import { JOB_STATE, JOB_RUN_STATUS, JOB_RULE_TYPE } from '@/constants';

import { useI18n } from '@/hooks/i18n';

import TicketStatusJobsDetailsExpandPanel from './partials/ticket-status-jobs-details-expand-panel.vue';

export default {
  components: { TicketStatusJobsDetailsExpandPanel },
  props: {
    headers: {
      type: Array,
      required: true,
    },
    items: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      default: 0,
    },
    options: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const isTicketStatusTab = computed(() => props.columns.some(col => col.value === 'active'));

    const statusItems = computed(() => {
      const i18nBase = isTicketStatusTab.value ? 'jobs.activeState' : 'jobs.status';

      return Object.values(JOB_STATE).map(value => ({
        value,
        text: t(`${i18nBase}.${value}`),
      }));
    });

    const getRuleName = item => item.rule_name ?? ((`${item.ticket_system_name || ''} - ${item.ticket_id || ''}`.trim() || '-'));
    const getAuthTokenName = item => item.auth_token_name ?? item.rule_name ?? '-';

    const ruleTypeLabel = (ruleType) => {
      if (!ruleType || !Object.values(JOB_RULE_TYPE).includes(ruleType)) {
        return '-';
      }

      return t(`jobs.ruleTypeValues.${ruleType}`);
    };

    const activeStateLabel = status => t(`jobs.activeState.${status}`) || t('jobs.status.unknown');

    const runStatusLabel = lastRunStatus => (
      (lastRunStatus === undefined || lastRunStatus === null)
        ? t('jobs.runStatus.inProgress')
        : (t(`jobs.runStatus.${lastRunStatus}`) || t('jobs.runStatus.inProgress'))
    );

    const statusColor = (status) => {
      const colors = {
        [JOB_STATE.running]: 'primary',
        [JOB_STATE.paused]: 'warning',
        [JOB_STATE.stopped]: 'error',
      };

      return colors[status] ?? 'grey';
    };

    const statusIcon = (status) => {
      const icons = {
        [JOB_STATE.running]: 'play_arrow',
        [JOB_STATE.paused]: 'pause',
        [JOB_STATE.stopped]: 'close',
      };

      return icons[status] ?? 'help_outline';
    };

    const shouldShowFinishDate = (item) => {
      const status = item.last_run_status ?? item.lastRunStatus;

      return status === JOB_RUN_STATUS.succeed || status === JOB_RUN_STATUS.failed;
    };

    const shouldShowFailReason = (item) => {
      const status = item.last_run_status ?? item.lastRunStatus;

      return status === JOB_RUN_STATUS.failed;
    };

    const shouldShowExpirationDate = (item) => {
      const status = item.last_run_status ?? item.lastRunStatus;

      return status === JOB_RUN_STATUS.succeed;
    };

    const canStart = item => item.status === JOB_STATE.stopped;
    const canStop = item => item.status === JOB_STATE.running;
    const canResume = item => item.status === JOB_STATE.paused;
    const canPause = item => item.status === JOB_STATE.running;

    const updateOptions = newOptions => emit('update:options', newOptions);

    const updateStatus = status => updateOptions({
      ...props.options,
      status: status ?? undefined,
      page: 1,
    });

    const emitAction = (action, item) => emit('action', { action, item });

    return {
      statusItems,
      getRuleName,
      getAuthTokenName,
      ruleTypeLabel,
      activeStateLabel,
      runStatusLabel,
      statusColor,
      statusIcon,
      shouldShowFinishDate,
      shouldShowFailReason,
      shouldShowExpirationDate,
      canStart,
      canStop,
      canResume,
      canPause,
      updateOptions,
      updateStatus,
      emitAction,
    };
  },
};
</script>
