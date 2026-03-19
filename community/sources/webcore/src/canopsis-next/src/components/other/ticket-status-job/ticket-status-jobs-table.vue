<template>
  <c-advanced-data-table
    ref="tableElement"
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledItem"
    :search-label="$t('jobs.searchByRuleName')"
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
    <template #mass-actions="{ selected: slotSelected }">
      <ticket-status-job-table-actions
        :items="slotSelected"
        class="mt-2"
      />
    </template>
    <template #header.data-table-select />
    <template #rule_type="{ item }">
      {{ $t(`jobs.types.${item.rule_type}`) }}
    </template>
    <template #status="{ item }">
      <ticket-status-jobs-active-state-icon :status="item.status" />
    </template>
    <template #last_run_status="{ item }">
      <ticket-status-jobs-last-run-status-icon :status="item.last_run_status" />
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
import { ref, computed } from 'vue';

import TicketStatusJobsActiveStateIcon from './partials/ticket-status-jobs-active-state-icon.vue';
import TicketStatusJobsDetailsExpandPanel from './partials/ticket-status-jobs-details-expand-panel.vue';
import TicketStatusJobsLastRunStatusIcon from './partials/ticket-status-jobs-last-run-status-icon.vue';
import TicketStatusJobsTableFilters from './partials/ticket-status-jobs-table-filters.vue';
import TicketStatusJobTableActions from './partials/ticket-status-job-table-actions.vue';

export default {
  components: {
    TicketStatusJobsActiveStateIcon,
    TicketStatusJobsDetailsExpandPanel,
    TicketStatusJobsLastRunStatusIcon,
    TicketStatusJobsTableFilters,
    TicketStatusJobTableActions,
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
    const tableElement = ref(null);

    const firstSelectedStatus = computed(() => tableElement.value?.selectedItems?.[0]?.status);

    /**
     * Disables table items that have a different status than the first selected item.
     * Ensures batch actions apply only to items with the same active state.
     *
     * @param {Object} item - Ticket status job entity
     * @returns {boolean} True if the item should be disabled for selection
     */
    const isDisabledItem = item => !isUndefined(firstSelectedStatus.value) && item.status !== firstSelectedStatus.value;

    /**
     * Emits options update to the parent.
     *
     * @param {Object} newOptions - New table options (pagination, sort, filters, etc.)
     */
    const updateOptions = newOptions => emit('update:options', newOptions);

    return {
      tableElement,
      isDisabledItem,
      updateOptions,
    };
  },
};
</script>
