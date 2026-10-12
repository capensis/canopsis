<template>
  <alarm-webhook-execution-timeline-card
    :step="step"
    :expandable="expandable"
    :expanded="expanded"
    :expanding="pending"
    :response="response"
    @expand="expandResponse"
  />
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import { useScenarioStore } from '@/stores/modules/scenario';

import AlarmWebhookExecutionTimelineCard from '@/components/other/alarm/partials/alarm-webhook-execution-timeline-card.vue';

export default {
  components: { AlarmWebhookExecutionTimelineCard },
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
      return [EVENT_ENTITY_TYPES.webhookComplete, EVENT_ENTITY_TYPES.webhookFail].includes(
        this.step._t,
      );
    },
  },
  methods: {
    async expandResponse(value) {
      if (value) {
        if (!this.response) {
          this.pending = true;

          try {
            this.response = await useScenarioStore().fetchTestScenarioExecutionWebhooksResponse({
              id: this.step.webhookId,
            });
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
