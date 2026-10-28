<template>
  <c-advanced-data-table
    ref="tableElement"
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledItem"
    :search-label="searchLabel"
    expand-content-class="secondary lighten-2"
    expand
    search
    select-all
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <ticket-status-jobs-table-filters
        v-if="isTicketStatusJob"
        :options="options"
        @update:options="updateOptions"
      />
      <job-management-table-filters
        v-else
        :options="options"
        :with-rule-type="jobType === JOBS_TABS.webhooks"
        @update:options="updateOptions"
      />
    </template>

    <template #mass-actions="{ selected, clearSelected }">
      <ticket-status-job-table-actions
        v-if="isTicketStatusJob"
        :items="selected"
        @clear:items="clearSelected"
      />
      <job-table-actions
        v-else
        :items="selected"
        :job-type="jobType"
        @clear:items="clearSelected"
      />
    </template>

    <template #[`header.data-table-select`] />

    <template #rule_type="{ item }">
      {{ t(`jobs.types.${item.rule_type}`) }}
    </template>

    <template #status="{ item }">
      <ticket-status-jobs-active-state-icon
        v-if="isTicketStatusJob"
        :status="item.status"
      />
      <job-execution-status-icon v-else :status="item.status" />
    </template>

    <template #last_run_status="{ item }">
      <ticket-status-jobs-last-run-status-icon :status="item.last_run_status" />
    </template>

    <template #started_at="{ item }">
      {{ formatDate(item.started_at ?? item.created_at) }}
    </template>

    <template #created_at="{ item }">
      {{ formatDate(item.created_at) }}
    </template>

    <template #completed_at="{ item }">
      {{ formatDate(item.completed_at, '-') }}
    </template>

    <template #checked_at="{ item }">
      {{ formatDate(item.checked_at, '-') }}
    </template>

    <template #next_check_at="{ item }">
      {{ formatDate(item.next_check_at, '-') }}
    </template>

    <template #expired_at="{ item }">
      {{ formatDate(item.expired_at, '-') }}
    </template>

    <template #fail_reason="{ item }">
      <c-ellipsis :text="item.fail_reason || '-'" :max-letters="70" hover />
    </template>

    <template #actions="{ item }">
      <ticket-status-job-table-actions v-if="isTicketStatusJob" :item="item" />
      <job-table-actions v-else :item="item" :job-type="jobType" />
    </template>

    <template #expand="{ item }">
      <job-details-expand-panel :item="item" :job-type="jobType" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { isUndefined } from 'lodash';
import { computed, ref } from 'vue';

import { JOB_EXECUTION_STATUS, JOB_STATUS, JOBS_TABS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import TicketStatusJobTableActions from '@/components/other/ticket-status-job/partials/ticket-status-job-table-actions.vue';
import TicketStatusJobsActiveStateIcon from '@/components/other/ticket-status-job/partials/ticket-status-jobs-active-state-icon.vue';
import TicketStatusJobsLastRunStatusIcon from '@/components/other/ticket-status-job/partials/ticket-status-jobs-last-run-status-icon.vue';
import TicketStatusJobsTableFilters from '@/components/other/ticket-status-job/partials/ticket-status-jobs-table-filters.vue';

import JobDetailsExpandPanel from './partials/job-details-expand-panel.vue';
import JobExecutionStatusIcon from './partials/job-execution-status-icon.vue';
import JobManagementTableFilters from './partials/job-management-table-filters.vue';
import JobTableActions from './partials/job-table-actions.vue';

export default {
  components: {
    JobDetailsExpandPanel,
    JobExecutionStatusIcon,
    JobManagementTableFilters,
    JobTableActions,
    TicketStatusJobTableActions,
    TicketStatusJobsActiveStateIcon,
    TicketStatusJobsLastRunStatusIcon,
    TicketStatusJobsTableFilters,
  },

  props: {
    jobType: {
      type: String,
      required: true,
    },
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

  emits: ['update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const tableElement = ref(null);
    const isTicketStatusJob = computed(() => props.jobType === JOBS_TABS.ticketStatus);
    const searchLabel = computed(() => {
      if (isTicketStatusJob.value) {
        return t('jobs.searchByRuleName');
      }

      return t(props.jobType === JOBS_TABS.authToken
        ? 'jobs.searchByAuthTokenName'
        : 'jobs.searchByRuleNameOnly');
    });
    const firstSelectedStatus = computed(() => tableElement.value?.selectedItems?.[0]?.status);
    const isDisabledItem = (item) => {
      if (!isTicketStatusJob.value) {
        return item.status !== JOB_EXECUTION_STATUS.running || item.stopped;
      }

      return item.status === JOB_STATUS.stopped ||
        (!isUndefined(firstSelectedStatus.value) && item.status !== firstSelectedStatus.value);
    };
    const updateOptions = newOptions => emit('update:options', newOptions);
    const formatDate = (date, defaultValue = '') => convertDateToString(date, undefined, defaultValue);

    return {
      t,
      JOBS_TABS,
      tableElement,
      isTicketStatusJob,
      searchLabel,
      isDisabledItem,
      updateOptions,
      formatDate,
    };
  },
};
</script>
