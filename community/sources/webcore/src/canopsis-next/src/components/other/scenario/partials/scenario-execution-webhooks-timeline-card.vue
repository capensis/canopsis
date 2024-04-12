<template>
  <alarm-test-query-execution-webhooks-timeline-card
    :step="step"
    :expandable="expandable"
    :expanded="expanded"
    :expanding="pending"
    :response="response"
    @expand="expandResponse"
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { EVENT_ENTITY_TYPES } from '@/constants';

import AlarmTestQueryExecutionWebhooksTimelineCard
  from '@/components/other/alarm/partials/alarm-test-query-execution-webhooks-timeline-card.vue';

const { mapActions: mapScenarioActions } = createNamespacedHelpers('scenario');

export default {
  components: { AlarmTestQueryExecutionWebhooksTimelineCard },
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      expanded: false,
      response: null,
    };
  },
  computed: {
    expandable() {
      return [EVENT_ENTITY_TYPES.webhookComplete, EVENT_ENTITY_TYPES.webhookFail].includes(this.step._t);
    },
  },
  methods: {
    ...mapScenarioActions({
      fetchTestScenarioExecutionWebhooksResponse: 'fetchTestScenarioExecutionWebhooksResponse',
    }),

    async expandResponse(value) {
      if (value) {
        if (!this.response) {
          this.pending = true;

          try {
            this.response = await this.fetchTestScenarioExecutionWebhooksResponse({ id: this.step.webhookId });
          } catch (err) {
            console.error(err);
          } finally {
            this.pending = false;
          }
        }

        this.expanded = true;
      } else {
        this.expanded = false;
      }
    },
  },
};
</script>
