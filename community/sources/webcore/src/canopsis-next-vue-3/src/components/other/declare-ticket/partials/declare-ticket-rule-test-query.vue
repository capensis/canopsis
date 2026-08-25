<template>
  <alarm-webhook-execution
    v-model:alarm="alarm"
    :execution-status="executionStatus"
    :alarms-patterns-params="alarmsPatternsParams"
    :pending="pending"
    :has-errors="hasErrors"
    @run:execution="runTestExecution"
    @clear:execution="clearWebhookStatus"
  >
    <template #webhooks="{ webhooks }">
      <alarm-webhook-execution-timeline :webhooks="webhooks">
        <template #card="{ step }">
          <declare-ticket-rule-execution-webhooks-timeline-card :step="step" />
        </template>
      </alarm-webhook-execution-timeline>
    </template>
  </alarm-webhook-execution>
</template>

<script setup>
import { toRef } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { formToDeclareTicketRule } from '@/helpers/entities/declare-ticket/rule/form';

import { useDeclareTicketRule } from '@/hooks/store/modules/declare-ticket-rule';

import { useAlarmWebhookTestQuery } from '@/components/other/alarm/hooks/alarm-webhook-test-query';

import AlarmWebhookExecutionTimeline from '@/components/other/alarm/partials/alarm-webhook-execution-timeline.vue';
import AlarmWebhookExecution from '@/components/other/alarm/partials/alarm-webhook-execution.vue';

import DeclareTicketRuleExecutionWebhooksTimelineCard from './declare-ticket-rule-execution-webhooks-timeline-card.vue';

const props = defineProps({
  form: {
    type: Object,
    default: () => ({}),
  },
});

const { createTestDeclareTicketExecution, fetchDeclareTicketExecutionWithoutStore } =
  useDeclareTicketRule();
const {
  alarm,
  pending,
  executionStatus,
  hasErrors,
  alarmsPatternsParams,
  clearWebhookStatus,
  runTestExecution,
} = useAlarmWebhookTestQuery({
  form: toRef(props, 'form'),
  socketRoom: SOCKET_ROOMS.declareticket,
  createExecution: createTestDeclareTicketExecution,
  fetchExecution: fetchDeclareTicketExecutionWithoutStore,
  getExecutionData: ({ alarm: alarmId, form }) => ({
    alarms: [alarmId],
    ...formToDeclareTicketRule(form),
  }),
});
</script>
