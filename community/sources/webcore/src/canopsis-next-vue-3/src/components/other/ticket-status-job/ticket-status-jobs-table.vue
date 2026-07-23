<template>
  <c-advanced-data-table
    ref="tableElement"
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledItem"
    :search-label="t('jobs.searchByRuleName')"
    expand-content-class="secondary lighten-2"
    expand
    search
    select-all
    advanced-pagination
    @update:options="updateOptions"
  >
    <template #toolbar>
      <ticket-status-jobs-table-filters :options="options" @update:options="updateOptions" />
    </template>

    <template #mass-actions="{ selected, clearSelected }">
      <ticket-status-job-table-actions :items="selected" @clear:items="clearSelected" />
    </template>

    <template #header.data-table-select />

    <template #rule_type="{ item }">
      {{ t(`jobs.types.${item.rule_type}`) }}
    </template>

    <template #status="{ item }">
      <ticket-status-jobs-active-state-icon :status="item.status" />
    </template>

    <template #last_run_status="{ item }">
      <ticket-status-jobs-last-run-status-icon :status="item.last_run_status" />
    </template>

    <template #created_at="{ item }">
      {{ formatDate(item.created_at) }}
    </template>

    <template #checked_at="{ item }">
      {{ formatDate(item.checked_at, '-') }}
    </template>

    <template #next_check_at="{ item }">
      {{ formatDate(item.next_check_at, '-') }}
    </template>

    <template #fail_reason="{ item }">
      <span class="d-inline-block text-truncate ticket-status-jobs-table__fail-reason">
        {{ item.fail_reason || '-' }}
      </span>
    </template>

    <template #actions="{ item }">
      <ticket-status-job-table-actions :item="item" />
    </template>

    <template #expand="{ item }">
      <ticket-status-jobs-details-expand-panel :item="item" />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { isUndefined } from 'lodash';
import { computed, ref } from 'vue';

import { JOB_STATUS } from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import TicketStatusJobTableActions from './partials/ticket-status-job-table-actions.vue';
import TicketStatusJobsActiveStateIcon from './partials/ticket-status-jobs-active-state-icon.vue';
import TicketStatusJobsDetailsExpandPanel from './partials/ticket-status-jobs-details-expand-panel.vue';
import TicketStatusJobsLastRunStatusIcon from './partials/ticket-status-jobs-last-run-status-icon.vue';
import TicketStatusJobsTableFilters from './partials/ticket-status-jobs-table-filters.vue';

export default {
  components: {
    TicketStatusJobTableActions,
    TicketStatusJobsActiveStateIcon,
    TicketStatusJobsDetailsExpandPanel,
    TicketStatusJobsLastRunStatusIcon,
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

  emits: ['update:options'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const tableElement = ref(null);
    const firstSelectedStatus = computed(() => tableElement.value?.selectedItems?.[0]?.status);

    const isDisabledItem = (item) =>
      item.status === JOB_STATUS.stopped ||
      (!isUndefined(firstSelectedStatus.value) && item.status !== firstSelectedStatus.value);
    const updateOptions = (newOptions) => emit('update:options', newOptions);
    const formatDate = (date, defaultValue = '') =>
      convertDateToString(date, undefined, defaultValue);

    return {
      t,
      tableElement,
      isDisabledItem,
      updateOptions,
      formatDate,
    };
  },
};
</script>

<style lang="scss" scoped>
.ticket-status-jobs-table__fail-reason {
  max-width: 100%;
}
</style>
