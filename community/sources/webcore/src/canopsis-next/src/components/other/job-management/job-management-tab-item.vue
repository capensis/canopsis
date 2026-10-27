<template>
  <job-management-table
    :job-type="tabId"
    :headers="headers"
    :items="jobs"
    :pending="jobsPending"
    :total-items="jobsMeta.total_count"
    :options="options"
    @update:options="updateOptions"
  />
</template>

<script>
import { computed, onMounted } from 'vue';

import { PAGINATION_LIMIT } from '@/config';
import { JOBS_TABS } from '@/constants';

import { convertQueryToRequest } from '@/helpers/query';

import { useI18n } from '@/hooks/i18n';
import { useLocalQueryWithOptions } from '@/hooks/query/shared';

import { useJobManagement } from './hooks/job-management';
import JobManagementTable from './job-management-table.vue';

export default {
  components: { JobManagementTable },

  props: {
    tabId: {
      type: String,
      required: true,
    },
  },

  setup(props, { expose }) {
    const { t } = useI18n();
    const {
      jobs,
      jobsMeta,
      jobsPending,
      fetchJobsList,
      fetchJobsListWithPreviousParams,
    } = useJobManagement(props.tabId);

    const headers = computed(() => {
      if (props.tabId === JOBS_TABS.ticketStatus) {
        return [
          { value: 'rule_name', text: t('jobs.ruleName'), sortable: true },
          { value: 'rule_type', text: t('jobs.ruleType'), sortable: true },
          { value: 'ticket_system_name', text: t('jobs.ticketSystemName'), sortable: true },
          { value: 'ticket_id', text: t('jobs.ticketNumber'), sortable: true },
          { value: 'status', text: t('jobs.activeState'), sortable: true },
          { value: 'last_run_status', text: t('jobs.lastStatus'), sortable: true },
          { value: 'created_at', text: t('jobs.startDate'), sortable: true },
          { value: 'checked_at', text: t('jobs.finishDate'), sortable: true },
          { value: 'fail_reason', text: t('jobs.failReason'), sortable: false },
          { value: 'actions', text: t('common.actionsLabel'), sortable: false },
        ];
      }

      const startDateField = props.tabId === JOBS_TABS.instructions ? 'started_at' : 'created_at';

      return [
        {
          value: 'rule_name',
          text: t(props.tabId === JOBS_TABS.authToken ? 'jobs.authTokenName' : 'jobs.ruleName'),
          sortable: true,
        },
        props.tabId === JOBS_TABS.webhooks && {
          value: 'rule_type',
          text: t('jobs.ruleType'),
          sortable: true,
        },
        { value: 'status', text: t('common.status'), sortable: true },
        { value: startDateField, text: t('jobs.startDate'), sortable: true },
        { value: 'completed_at', text: t('jobs.finishDate'), sortable: true },
        props.tabId === JOBS_TABS.authToken && {
          value: 'expired_at',
          text: t('jobs.expirationDate'),
          sortable: true,
        },
        { value: 'fail_reason', text: t('jobs.failReason'), sortable: false },
        { value: 'actions', text: t('common.actionsLabel'), sortable: false },
      ].filter(Boolean);
    });

    const initialQuery = {
      page: 1,
      itemsPerPage: PAGINATION_LIMIT,
      search: '',
      sortBy: props.tabId === JOBS_TABS.ticketStatus ? [] : ['status'],
      sortDesc: props.tabId === JOBS_TABS.ticketStatus ? [] : [false],
    };
    const fetchList = (fetchQuery = initialQuery) => fetchJobsList({
        params: convertQueryToRequest(fetchQuery),
      });
    const { options, updateOptions } = useLocalQueryWithOptions({
      initialQuery,
      onUpdate: fetchList,
    });

    expose({ refresh: fetchJobsListWithPreviousParams });
    onMounted(fetchList);

    return {
      jobs,
      jobsMeta,
      jobsPending,
      headers,
      options,
      updateOptions,
    };
  },
};
</script>
