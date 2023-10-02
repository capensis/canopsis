<template lang="pug">
  v-layout(column)
    v-layout.mb-4(v-if="showStatus", row, align-center)
      v-chip.mr-2(v-if="showRuleName", color="primary", text-color="white", small) {{ ruleName }}
      span.subheading.mr-5 {{ $t('declareTicket.webhookStatus') }}:
      declare-ticket-rule-execution-status(
        :running="isExecutionsRunning",
        :success="isExecutionsSucceeded",
        :fail-reason="failReason"
      )
    declare-ticket-rule-execution-alarms(:alarm-executions="executions", :is-one-execution="isOneExecution")
</template>

<script>
import { MODALS } from '@/constants';

import {
  isDeclareTicketExecutionRunning,
  isDeclareTicketExecutionSucceeded,
} from '@/helpers/entities/declare-ticket/rule/form';

import DeclareTicketRuleExecutionStatus from './declare-ticket-rule-execution-status.vue';
import DeclareTicketRuleExecutionAlarms from './declare-ticket-rule-execution-alarms.vue';

/**
 * Modal to execute declare tickets
 */
export default {
  name: MODALS.executeDeclareTickets,
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
      return this.executions.map(execution => execution.fail_reason).join('\n');
    },
  },
};
</script>
