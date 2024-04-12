<template>
  <v-layout column>
    <v-layout
      class="gap-3"
      align-center
      justify-space-between
    >
      <v-flex xs10>
        <c-alarm-field
          v-field="alarm"
          :disabled="pending || isExecutionRunning"
          :params="alarmsPatternsParams"
          name="alarms"
        />
      </v-flex>
      <v-flex xs2>
        <v-btn
          :disabled="hasErrors || !alarm"
          :loading="pending || isExecutionRunning"
          class="white--text"
          color="orange"
          block
          @click="runTestExecution"
        >
          {{ $t('alarm.runTest') }}
        </v-btn>
      </v-flex>
    </v-layout>
    <v-expand-transition>
      <v-layout
        v-if="executionStatus"
        column
      >
        <v-layout
          class="mb-4"
          align-center
        >
          <span class="text-subtitle-1 mr-5">{{ $t('declareTicket.webhookStatus') }}:</span>
          <alarm-test-query-execution-status
            :running="isExecutionRunning"
            :success="isExecutionSucceeded"
            :fail-reason="executionStatus.fail_reason"
          />
          <c-action-btn
            v-if="isExecutionSucceeded || isExecutionFailed"
            type="delete"
            @click="clearWebhookStatus"
          />
        </v-layout>
        <v-card v-if="isSomeOneWebhookStarted">
          <v-card-text>
            <slot :webhooks="executionStatus.webhooks" name="webhooks" />
          </v-card-text>
        </v-card>
      </v-layout>
    </v-expand-transition>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import {
  isWebhookExecutionFailed,
  isWebhookExecutionRunning,
  isWebhookExecutionSucceeded,
  isWebhookExecutionWaiting,
} from '@/helpers/entities/webhook-execution/entity';

import AlarmTestQueryExecutionStatus from './alarm-test-query-execution-status.vue';

export default {
  components: { AlarmTestQueryExecutionStatus },
  model: {
    prop: 'alarm',
    event: 'input',
  },
  props: {
    alarm: {
      type: String,
      required: false,
    },
    executionStatus: {
      type: Object,
      required: false,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    hasErrors: {
      type: Boolean,
      default: false,
    },
    alarmsPatternsParams: {
      type: Object,
      required: false,
    },
  },
  setup(props, { emit }) {
    const isExecutionRunning = computed(
      () => props.executionStatus && isWebhookExecutionRunning(props.executionStatus),
    );
    const isExecutionSucceeded = computed(() => isWebhookExecutionSucceeded(props.executionStatus));
    const isExecutionFailed = computed(() => isWebhookExecutionFailed(props.executionStatus));
    const isSomeOneWebhookStarted = computed(() => props.executionStatus?.webhooks.some(
      webhook => !isWebhookExecutionWaiting(webhook),
    ));

    const runTestExecution = () => emit('run:execution');
    const clearWebhookStatus = () => emit('run:execution');

    return {
      isExecutionRunning,
      isExecutionSucceeded,
      isExecutionFailed,
      isSomeOneWebhookStarted,

      runTestExecution,
      clearWebhookStatus,
    };
  },
};
</script>
