<template>
  <v-row class="healthcheck-graphs pa-4" justify="center" no-gutters>
    <v-col cols="12" md="10" lg="8">
      <c-information-block :title="$t('healthcheck.messagesLastHour')" class="mb-5">
        <healthcheck-last-hour :max-messages-per-minute="maxMessagesPerMinute" class="mt-4" />
      </c-information-block>

      <c-information-block :title="$t('healthcheck.messagesHistory')">
        <healthcheck-history-graph :max-messages-per-minute="maxMessagesPerMinute" class="mt-4" />
      </c-information-block>
    </v-col>
  </v-row>
</template>

<script>
import HealthcheckLastHour from './partials/healthcheck-last-hour.vue';
import HealthcheckHistoryGraph from './partials/healthcheck-history-graph.vue';

export default {
  components: {
    HealthcheckHistoryGraph,
    HealthcheckLastHour,
  },
  props: {
    maxMessagesPerMinute: {
      type: Number,
      required: true,
    },
  },
};
</script>
