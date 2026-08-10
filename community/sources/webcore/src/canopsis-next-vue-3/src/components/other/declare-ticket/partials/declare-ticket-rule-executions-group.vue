<template>
  <c-layout column>
    <c-layout v-if="showStatus" class="mb-4" align="center">
      <v-chip
        v-if="showRuleName"
        class="mr-2"
        color="primary"
        size="small"
      >
        {{ ruleName }}
      </v-chip>
      <span class="text-subtitle-1 mr-5"> {{ $t('common.webhookStatus') }}: </span>
      <alarm-webhook-execution-status
        :running="isExecutionsRunning"
        :success="isExecutionsSucceeded"
        :fail-reason="failReason"
      />
    </c-layout>
    <declare-ticket-rule-execution-alarms
      :alarm-executions="executions"
      :is-one-execution="isOneExecution"
    />
  </c-layout>
</template>

<script>
import { keyBy } from 'lodash';

import AlarmWebhookExecutionStatus from '@/components/other/alarm/partials/alarm-webhook-execution-status.vue';
import {
  isWebhookExecutionRunning,
  isWebhookExecutionSucceeded,
} from '@/helpers/entities/webhook-execution/entity';

import DeclareTicketRuleExecutionAlarms from './declare-ticket-rule-execution-alarms.vue';

export default {
  components: {
    DeclareTicketRuleExecutionAlarms,
    AlarmWebhookExecutionStatus,
  },
  props: {
    executions: {
      type: Array,
      required: true,
    },
    showRuleName: {
      type: Boolean,
      default: false,
    },
    showStatus: {
      type: Boolean,
      default: false,
    },
    isOneExecution: {
      type: Boolean,
      default: false,
    },
    ruleName: {
      type: String,
      required: false,
    },
  },
  computed: {
    isExecutionsRunning() {
      return this.executions.some(isWebhookExecutionRunning);
    },

    isExecutionsSucceeded() {
      return this.executions.every(isWebhookExecutionSucceeded);
    },

    failReason() {
      return Object.values(keyBy(this.executions, 'executionId'))
        .map(execution => execution.fail_reason)
        .join('\n');
    },
  },
};
</script>
