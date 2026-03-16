<template>
  <c-advanced-data-table
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    expand
    search
    select-all
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <ticket-status-jobs-table-filters
        :options="options"
        @update:options="updateOptions"
      />
    </template>
    <template #mass-actions="{ selected }">
      <c-action-btn
        v-if="removable"
        type="delete"
        @click="$emit('remove-selected', selected)"
      />
    </template>
    <template #active="{ item }">
      <ticket-status-jobs-active-state-icon :status="item.status" />
    </template>
    <template #status="{ item }">
      <ticket-status-jobs-run-status-icon :status="item.last_run_status" />
    </template>
    <template #created_at="{ item }">
      {{ item.created_at | date }}
    </template>
    <template #checked_at="{ item }">
      {{ item.checked_at | date('long', '-') }}
    </template>
    <template #next_check_at="{ item }">
      {{ item.next_check_at | date('long', '-') }}
    </template>
    <template #fail_reason="{ item }">
      {{ item.fail_reason || '-' }}
    </template>
    <template #expand="{ item }">
      <ticket-status-jobs-details-expand-panel :item="item" />
    </template>
    <template #actions="{ item }">
      <v-layout>
        <c-action-btn
          :tooltip="$t('jobs.actions.editJob')"
          icon="edit"
          @click="emitAction('edit', item)"
        />
        <c-action-btn
          v-if="canStart(item)"
          :tooltip="$t('jobs.actions.startJob')"
          icon="play_arrow"
          @click="emitAction('start', item)"
        />
        <c-action-btn
          v-if="canStop(item)"
          :tooltip="$t('jobs.actions.stopJob')"
          icon="stop"
          @click="emitAction('stop', item)"
        />
        <c-action-btn
          v-if="canResume(item)"
          :tooltip="$t('jobs.actions.repeatJob')"
          icon="refresh"
          @click="emitAction('retry', item)"
        />
        <c-action-btn
          v-if="canPause(item)"
          :tooltip="$t('jobs.actions.pauseJob')"
          icon="pause"
          @click="emitAction('pause', item)"
        />
      </v-layout>
    </template>
  </c-advanced-data-table>
</template>

<script>
import { JOB_STATUS } from '@/constants';

import TicketStatusJobsActiveStateIcon from './partials/ticket-status-jobs-active-state-icon.vue';
import TicketStatusJobsDetailsExpandPanel from './partials/ticket-status-jobs-details-expand-panel.vue';
import TicketStatusJobsRunStatusIcon from './partials/ticket-status-jobs-run-status-icon.vue';
import TicketStatusJobsTableFilters from './partials/ticket-status-jobs-table-filters.vue';

export default {
  components: {
    TicketStatusJobsActiveStateIcon,
    TicketStatusJobsDetailsExpandPanel,
    TicketStatusJobsRunStatusIcon,
    TicketStatusJobsTableFilters,
  },
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
    const canStart = item => true || item.status === JOB_STATUS.stopped;
    const canStop = item => true || item.status === JOB_STATUS.running;
    const canResume = item => true || item.status === JOB_STATUS.paused;
    const canPause = item => true || item.status === JOB_STATUS.running;

    const updateOptions = newOptions => emit('update:options', newOptions);

    const emitAction = (action, item) => emit('action', { action, item });

    return {
      canStart,
      canStop,
      canResume,
      canPause,
      updateOptions,
      emitAction,
    };
  },
};
</script>
