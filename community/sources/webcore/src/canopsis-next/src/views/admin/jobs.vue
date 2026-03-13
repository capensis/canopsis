<template>
  <div>
    <c-page-header />
    <ticket-status-jobs-tabs />
  </div>
</template>

<script>
import TicketStatusJobsTabs from '@/components/other/ticket-status-job/ticket-status-jobs-tabs.vue';

export default {
  components: { TicketStatusJobsTabs },
};
</script>
