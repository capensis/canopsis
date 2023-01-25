<template lang="pug">
  v-layout(column)
    declare-ticket-rule-execution-alarms(:is-one-execution="isOneExecution", :alarm-executions="alarmExecutions")
</template>

<script>
import DeclareTicketRuleExecutionAlarms from './declare-ticket-rule-execution-alarms.vue';
import DeclareTicketRuleExecutionStatus from './declare-ticket-rule-execution-status.vue';

export default {
  components: { DeclareTicketRuleExecutionStatus, DeclareTicketRuleExecutionAlarms },
  props: {
    alarmExecutions: {
      type: Array,
      default: () => [],
    },
    isOneExecution: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
