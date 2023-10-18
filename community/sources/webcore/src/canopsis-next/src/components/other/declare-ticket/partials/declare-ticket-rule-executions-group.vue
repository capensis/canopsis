<template>
  <v-layout column>
    <v-layout
      class="mb-4"
      v-if="showStatus"
      align-center
    >
      <v-chip
        class="mr-2"
        v-if="showRuleName"
        color="primary"
        text-color="white"
        small
      >
        {{ ruleName }}
      </v-chip>
      <span class="text-subtitle-1 mr-5">{{ $t('declareTicket.webhookStatus') }}:</span>
      <declare-ticket-rule-execution-status
        :running="isExecutionsRunning"
        :success="isExecutionsSucceeded"
        :fail-reason="failReason"
      />
    </v-layout>
    <declare-ticket-rule-execution-alarms
      :alarm-executions="executions"
      :is-one-execution="isOneExecution"
    />
  </v-layout>
</template>

<script>
import { keyBy } from 'lodash';

import {
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
} from '@/helpers/entities/declare-ticket/rule/form';

import DeclareTicketRuleExecutionStatus from './declare-ticket-rule-execution-status.vue';
import DeclareTicketRuleExecutionAlarms from './declare-ticket-rule-execution-alarms.vue';

export default {
  components: {
    DeclareTicketRuleExecutionAlarms,
    DeclareTicketRuleExecutionStatus,
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
      return this.executions.some(isDeclareTicketExecutionRunning);
    },

    isExecutionsSucceeded() {
      return this.executions.every(isDeclareTicketExecutionSucceeded);
    },

    failReason() {
      return Object.values(keyBy(this.executions, 'executionId')).map(execution => execution.fail_reason).join('\n');
    },
  },
};
</script>
