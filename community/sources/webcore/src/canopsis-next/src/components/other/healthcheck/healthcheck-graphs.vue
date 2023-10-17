<template>
  <v-flex
    class="pa-4"
    md10
    offset-md1
    lg8
    offset-lg2
  >
    <c-information-block
      class="mb-5"
      :title="$t('healthcheck.messagesLastHour')"
    >
      <healthcheck-last-hour-graph
        class="mt-4"
        :max-messages-per-minute="maxMessagesPerMinute"
      />
    </c-information-block>
    <c-information-block :title="$t('healthcheck.messagesHistory')">
      <healthcheck-history-graph
        class="mt-4"
        :max-messages-per-minute="maxMessagesPerMinute"
      />
    </c-information-block>
  </v-flex>
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
