<template>
  <c-page @refresh="refresh">
    <ticket-status-jobs-tabs />
  </c-page>
</template>

<script>
import TicketStatusJobsTabs from '@/components/other/ticket-status-job/ticket-status-jobs-tabs.vue';
import { useTicketStatusJob } from '@/hooks/store/modules/ticket-status-job';

export default {
  components: {
    TicketStatusJobsTabs,
  },

  setup() {
    const { fetchTicketStatusJobsListWithPreviousParams } = useTicketStatusJob();

    return {
      refresh: fetchTicketStatusJobsListWithPreviousParams,
    };
  },
};
</script>
