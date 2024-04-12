<template>
  <alarms-time-line-steps :steps="webhooksTimelineSteps">
    <template #card="{ step }">
      <slot :step="step" name="card" />
    </template>
  </alarms-time-line-steps>
</template>

<script>
import { EVENT_ENTITY_TYPES } from '@/constants';

import {
  isDeclareTicketExecutionFailed,
  isDeclareTicketExecutionSucceeded,
  isDeclareTicketExecutionWaiting,
} from '@/helpers/entities/declare-ticket/rule/form';

import AlarmsTimeLineSteps from '@/components/widgets/alarm/time-line/alarms-time-line-steps.vue';

export default {
  components: { AlarmsTimeLineSteps },
  props: {
    webhooks: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    webhooksTimelineSteps() {
      return this.webhooks.reduce((acc, webhook) => {
        if (isDeclareTicketExecutionWaiting(webhook)) {
          return acc;
        }

        const webhookId = webhook._id;

        acc.push({
          webhookId,
          t: webhook.launched_at > 0 ? webhook.launched_at : webhook.created_at,
          _t: EVENT_ENTITY_TYPES.webhookStart,
        });

        if (isDeclareTicketExecutionSucceeded(webhook)) {
          acc.push({
            webhookId,
            t: webhook.completed_at,
            _t: EVENT_ENTITY_TYPES.webhookComplete,
          });
        }

        if (isDeclareTicketExecutionFailed(webhook)) {
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
