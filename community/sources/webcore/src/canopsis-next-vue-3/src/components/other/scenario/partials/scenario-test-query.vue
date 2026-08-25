<template>
  <alarm-webhook-execution
    v-model:alarm="alarm"
    :execution-status="executionStatus"
    :alarms-patterns-params="alarmsPatternsParams"
    :pending="pending"
    :has-errors="hasErrors"
    class="mt-3"
    @run:execution="runTestExecution"
    @clear:execution="clearWebhookStatus"
  >
    <template #webhooks="{ webhooks }">
      <alarm-webhook-execution-timeline :webhooks="webhooks">
        <template #card="{ step }">
          <scenario-execution-webhooks-timeline-card :step="step" />
        </template>
      </alarm-webhook-execution-timeline>
    </template>
  </alarm-webhook-execution>
</template>

<script setup>
import { toRef } from 'vue';

import { SOCKET_ROOMS } from '@/config';

import { formToScenario } from '@/helpers/entities/scenario/form';

import { useScenario } from '@/hooks/store/modules/scenario';

import { useAlarmWebhookTestQuery } from '@/components/other/alarm/hooks/alarm-webhook-test-query';

import AlarmWebhookExecutionTimeline from '@/components/other/alarm/partials/alarm-webhook-execution-timeline.vue';
import AlarmWebhookExecution from '@/components/other/alarm/partials/alarm-webhook-execution.vue';

import ScenarioExecutionWebhooksTimelineCard from './scenario-execution-webhooks-timeline-card.vue';

const props = defineProps({
  form: {
    type: Object,
    default: () => ({}),
  },
});

const { createTestScenarioExecution, fetchTestScenarioExecutionWithoutStore } = useScenario();
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
  socketRoom: SOCKET_ROOMS.testscenario,
  createExecution: createTestScenarioExecution,
  fetchExecution: fetchTestScenarioExecutionWithoutStore,
  getExecutionData: ({ alarm: alarmId, form }) => ({
    alarm: alarmId,
    ...formToScenario(form),
  }),
});
</script>
