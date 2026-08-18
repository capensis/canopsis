<template>
  <alarm-webhook-execution-timeline-steps :steps="webhooksTimelineSteps">
    <template #card="{ step }">
      <slot :step="step" name="card" />
    </template>
  </alarm-webhook-execution-timeline-steps>
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import {
  isWebhookExecutionFailed,
  isWebhookExecutionSucceeded,
  isWebhookExecutionWaiting,
} from '@/helpers/entities/webhook-execution/entity';

import AlarmWebhookExecutionTimelineSteps from './alarm-webhook-execution-timeline-steps.vue';

export default {
  components: {
    AlarmWebhookExecutionTimelineSteps,
  },
  props: {
    webhooks: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    webhooksTimelineSteps() {
      return this.webhooks.reduce((acc, webhook) => {
        if (isWebhookExecutionWaiting(webhook)) {
          return acc;
        }

        const webhookId = webhook._id;

        acc.push({
          webhookId,
          t: webhook.launched_at > 0 ? webhook.launched_at : webhook.created_at,
          _t: EVENT_ENTITY_TYPES.webhookStart,
        });

        if (isWebhookExecutionSucceeded(webhook)) {
          acc.push({
            webhookId,
            t: webhook.completed_at,
            _t: EVENT_ENTITY_TYPES.webhookComplete,
          });
        }

        if (isWebhookExecutionFailed(webhook)) {
          acc.push({
            webhookId,
            t: webhook.completed_at,
            _t: EVENT_ENTITY_TYPES.webhookFail,
            message: webhook.fail_reason,
          });
        }

        return acc;
      }, []);
    },
  },
};
</script>
