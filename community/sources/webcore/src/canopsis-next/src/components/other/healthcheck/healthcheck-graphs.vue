<template>
  <v-flex
    class="pa-4"
    md10
    offset-md1
    lg8
    offset-lg2
  >
    <c-information-block
      :title="$t('healthcheck.messagesLastHour')"
      class="mb-5"
    >
      <healthcheck-last-hour
        :max-messages-per-minute="maxMessagesPerMinute"
        class="mt-4"
      />
    </c-information-block>
    <c-information-block :title="$t('healthcheck.messagesHistory')">
      <healthcheck-history-graph
        :max-messages-per-minute="maxMessagesPerMinute"
        class="mt-4"
      />
    </c-information-block>
  </v-flex>
</template>

<script>
import HealthcheckLastHour from '@/components/other/healthcheck/partials/healthcheck-last-hour.vue';

import HealthcheckHistoryGraph from './partials/healthcheck-history-graph.vue';

export default {
  components: { HealthcheckLastHour, HealthcheckHistoryGraph },
  props: {
    maxMessagesPerMinute: {
      type: Number,
      required: true,
    },
  },
};
</script>
