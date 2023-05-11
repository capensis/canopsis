<template lang="pug">
  v-flex.pa-4(v-else, md10, offset-md1, lg8, offset-lg2)
    c-information-block.mb-5(:title="$t('healthcheck.messagesLastHour')")
      healthcheck-last-hour-graph.mt-4(:max-messages-per-minute="maxMessagesPerMinute")
    c-information-block(:title="$t('healthcheck.messagesHistory')")
      healthcheck-history-graph.mt-4(:max-messages-per-minute="maxMessagesPerMinute")
</template>

<script>
import HealthcheckHistoryGraph from './partials/healthcheck-history-graph.vue';
import HealthcheckLastHourGraph from './partials/healthcheck-last-hour-graph.vue';

export default {
  components: { HealthcheckHistoryGraph, HealthcheckLastHourGraph },
  props: {
    maxMessagesPerMinute: {
      type: Number,
      required: true,
    },
  },
};
</script>
