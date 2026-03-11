<template>
  <div>
    <c-page-header />
    <jobs-tabs
      :items="items"
      :pending="pending"
      :total-items="meta.total_count"
      :options="options"
      @update:options="updateOptions"
      @action="handleAction"
    />
  </div>
</template>

<script>
import { onMounted } from 'vue';

import { useTicketStatusJob } from '@/hooks/store/modules/ticket-status-job';
import { useFetchListWithoutStoreWithOptions } from '@/hooks/query/shared';

import JobsTabs from '@/components/other/jobs/jobs-tabs.vue';

export default {
  components: { JobsTabs },
  setup() {
    const { fetchTicketStatusJobsListWithoutStore } = useTicketStatusJob();

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
      },
      fetchListHandler: ({ params }) => fetchTicketStatusJobsListWithoutStore({ params }),
    });

    const handleAction = () => {
      // TODO: Implement start, stop, resume, pause, edit API calls when endpoints are available
    };

    onMounted(fetchList);

    return {
      items,
      meta,
      pending,
      options,
      updateOptions,
      handleAction,
    };
  },
};
</script>
