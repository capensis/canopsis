<template>
  <v-card-text>
    <ticket-status-jobs-table
      :headers="headers"
      :items="items"
      :pending="pending"
      :total-items="totalItems"
      :options="options"
      @update:options="updateOptions"
      @edit="edit"
      @play="play"
      @pause="pause"
      @stop="stop"
      @play:selected="playSelected"
      @pause:selected="pauseSelected"
      @stop:selected="stopSelected"
    />
  </v-card-text>
</template>

<script>
import { computed, onMounted } from 'vue';

import { MODALS, TICKET_STATUS_JOBS_TABS } from '@/constants';

import { useI18n } from '@/hooks/i18n';
import { useModals } from '@/hooks/modals';
import { useTicketStatusJob } from '@/hooks/store/modules/ticket-status-job';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import TicketStatusJobsTable from '@/components/other/ticket-status-job/ticket-status-jobs-table.vue';

export default {
  components: { TicketStatusJobsTable },
  props: {
    tabId: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const { t } = useI18n();
    const modals = useModals();
    const {
      fetchTicketStatusJobsListWithoutStore,
      updateTicketStatusJob,
      playTicketStatusJob,
      pauseTicketStatusJob,
      stopTicketStatusJob,
    } = useTicketStatusJob();

    const label = computed(() => ({
      [TICKET_STATUS_JOBS_TABS.instructions]: t('jobs.tabs.instructions'),
      [TICKET_STATUS_JOBS_TABS.webhooks]: t('jobs.tabs.webhooks'),
      [TICKET_STATUS_JOBS_TABS.ticketStatus]: t('jobs.tabs.ticketStatus'),
      [TICKET_STATUS_JOBS_TABS.authToken]: t('jobs.tabs.authToken'),
    })[props.tabId] ?? '');

    const headers = computed(() => {
      const isInstructionsTab = props.tabId === TICKET_STATUS_JOBS_TABS.instructions;
      const isWebhooksTab = props.tabId === TICKET_STATUS_JOBS_TABS.webhooks;
      const isTicketStatusTab = props.tabId === TICKET_STATUS_JOBS_TABS.ticketStatus;
      const isAuthTokenTab = props.tabId === TICKET_STATUS_JOBS_TABS.authToken;

      return [
        isAuthTokenTab && { value: 'authTokenName', text: t('jobs.authTokenName'), sortable: true },
        isTicketStatusTab && { value: 'ticket_system_name', text: t('jobs.ticketSystemName'), sortable: true },
        isTicketStatusTab && { value: 'ticket_id', text: t('jobs.ticketNumber'), sortable: true },
        isTicketStatusTab && { value: 'active', text: t('jobs.active'), sortable: true },
        (isInstructionsTab || isWebhooksTab) && { value: 'rule_name', text: t('jobs.ruleName'), sortable: true },
        isWebhooksTab && { value: 'ruleType', text: t('jobs.ruleType'), sortable: true },
        { value: 'status', text: t('jobs.statusLabel'), sortable: true },
        { value: 'created_at', text: t('jobs.startDate'), sortable: true },
        { value: 'checked_at', text: t('jobs.finishDate'), sortable: true },
        { value: 'fail_reason', text: t('jobs.failReason'), sortable: false },
        isAuthTokenTab && { value: 'next_check_at', text: t('jobs.expirationDate'), sortable: true },
        { value: 'actions', text: t('common.actionsLabel'), sortable: false },
      ].filter(Boolean);
    });

    const {
      data: items,
      meta,
      pending,
      options,
      updateOptions,
      fetchList,
    } = useFetchListWithoutStoreWithOptions({
      initialQuery: {
        page: 1,
        itemsPerPage: 10,
        search: '',
        status: undefined,
        last_run_status: undefined,
      },
      fetchListHandler: fetchTicketStatusJobsListWithoutStore,
    });

    const totalItems = computed(() => meta.value?.total_count ?? 0);

    const edit = ticketStatusJob => modals.show({
      name: MODALS.createTicketStatusJob,
      config: {
        ticketStatusJob,
        action: async (newTicketStatusJob) => {
          await updateTicketStatusJob({ id: ticketStatusJob._id, data: newTicketStatusJob });
          await fetchList();
        },
      },
    });
    const play = (item) => {
      playTicketStatusJob(item);
    };
    const pause = (item) => {
      pauseTicketStatusJob(item);
    };
    const stop = (item) => {
      stopTicketStatusJob(item);
    };

    const playSelected = (selected) => {
      playTicketStatusJob(selected);
    };
    const pauseSelected = (selected) => {
      pauseTicketStatusJob(selected);
    };
    const stopSelected = (selected) => {
      stopTicketStatusJob(selected);
    };

    onMounted(fetchList);

    return {
      items,
      pending,
      totalItems,
      options,
      label,
      headers,
      updateOptions,
      edit,
      play,
      pause,
      stop,
      playSelected,
      pauseSelected,
      stopSelected,
    };
  },
};
</script>
