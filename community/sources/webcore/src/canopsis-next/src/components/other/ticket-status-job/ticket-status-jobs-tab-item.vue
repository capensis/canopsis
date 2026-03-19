<template>
  <v-card-text>
    <ticket-status-jobs-table
      :headers="headers"
      :items="ticketStatusJobs"
      :pending="ticketStatusJobsPending"
      :total-items="ticketStatusJobsMeta.total_count"
      :options="options"
      @update:options="updateOptions"
    />
  </v-card-text>
</template>

<script>
import { computed, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { JOBS_TABS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';
import { useTicketStatusJob } from '@/hooks/store/modules/ticket-status-job';

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
    const {
      ticketStatusJobs,
      ticketStatusJobsMeta,
      ticketStatusJobsPending,
      fetchTicketStatusJobsList,
    } = useTicketStatusJob();

    const label = computed(() => ({
      [JOBS_TABS.ticketStatus]: t('jobs.tabs.ticketStatus'),
    })[props.tabId] ?? '');

    const headers = computed(() => {
      const isTicketStatusTab = props.tabId === JOBS_TABS.ticketStatus;

      return [
        { value: 'rule_name', text: t('jobs.ruleName'), sortable: true },
        { value: 'rule_type', text: t('jobs.ruleType'), sortable: true },
        isTicketStatusTab && { value: 'ticket_system_name', text: t('jobs.ticketSystemName'), sortable: true },
        isTicketStatusTab && { value: 'ticket_id', text: t('jobs.ticketNumber'), sortable: true },
        isTicketStatusTab && { value: 'status', text: t('jobs.activeState'), sortable: true },
        { value: 'last_run_status', text: t('jobs.lastStatus'), sortable: true },
        { value: 'created_at', text: t('jobs.startDate'), sortable: true },
        { value: 'checked_at', text: t('jobs.finishDate'), sortable: true },
        { value: 'fail_reason', text: t('jobs.failReason'), sortable: false },
        { value: 'actions', text: t('common.actionsLabel'), sortable: false },
      ].filter(Boolean);
    });

    const initialQuery = {
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      search: '',
      sortBy: [],
      sortDesc: [],
    };

    const fetchList = (fetchQuery = initialQuery) => fetchTicketStatusJobsList({
      params: convertQueryToRequest(fetchQuery),
    });

    const {
      options,
      updateOptions,
    } = useLocalQueryWithOptions({
      initialQuery,
      onUpdate: fetchList,
    });

    onMounted(fetchList);

    return {
      ticketStatusJobs,
      ticketStatusJobsMeta,
      ticketStatusJobsPending,
      options,
      label,
      headers,
      updateOptions,
    };
  },
};
</script>
