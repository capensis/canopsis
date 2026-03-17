<template>
  <c-advanced-data-table
    ref="tableElement"
    :headers="headers"
    :items="items"
    :loading="pending"
    :total-items="totalItems"
    :options="options"
    :is-disabled-item="isDisabledItem"
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
        @edit="$emit('edit', $event)"
        @play="$emit('play', $event)"
        @stop="$emit('stop', $event)"
        @retry="$emit('retry', $event)"
        @pause="$emit('pause', $event)"
      />
    </template>
    <template #header.data-table-select />
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
      <ticket-status-job-table-actions
        :item="item"
        @edit="$emit('edit', $event)"
        @play="$emit('play', $event)"
        @stop="$emit('stop', $event)"
        @retry="$emit('retry', $event)"
        @pause="$emit('pause', $event)"
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { isUndefined } from 'lodash';
import { ref, computed } from 'vue';

import TicketStatusJobsActiveStateIcon from './partials/ticket-status-jobs-active-state-icon.vue';
import TicketStatusJobsDetailsExpandPanel from './partials/ticket-status-jobs-details-expand-panel.vue';
import TicketStatusJobsRunStatusIcon from './partials/ticket-status-jobs-run-status-icon.vue';
import TicketStatusJobsTableFilters from './partials/ticket-status-jobs-table-filters.vue';
import TicketStatusJobTableActions from './partials/ticket-status-job-table-actions.vue';

export default {
  components: {
    TicketStatusJobsActiveStateIcon,
    TicketStatusJobsDetailsExpandPanel,
    TicketStatusJobsRunStatusIcon,
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

    const isDisabledItem = item => !isUndefined(firstSelectedStatus.value) && item.status !== firstSelectedStatus.value;

    const itemsWithSelectable = computed(() => (
      firstSelectedStatus.value
        ? props.items.map(item => ({ ...item, isSelectable: item.status === firstSelectedStatus.value }))
        : props.items
    ));

    const updateOptions = newOptions => emit('update:options', newOptions);

    return {
      tableElement,
      isDisabledItem,
      itemsWithSelectable,
      updateOptions,
    };
  },
};
</script>
