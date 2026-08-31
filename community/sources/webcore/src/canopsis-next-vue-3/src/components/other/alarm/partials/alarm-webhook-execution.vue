<template>
  <c-layout column>
    <v-row align="center" justify="space-between" density="comfortable">
      <v-col cols="10">
        <c-alarm-field
          :value="alarm"
          :disabled="pending || isExecutionRunning"
          :params="alarmsPatternsParams"
          name="alarms"
          @update:value="$emit('update:alarm', $event)"
        />
      </v-col>
      <v-col cols="2">
        <v-btn
          :disabled="hasErrors || !alarm"
          :loading="pending || isExecutionRunning"
          class="text-white"
          color="orange"
          block
          @click="runTestExecution"
        >
          {{ $t('alarm.runTest') }}
        </v-btn>
      </v-col>
    </v-row>
    <v-expand-transition>
      <c-layout v-if="executionStatus" column>
        <c-layout class="mb-4" align="center">
          <span class="text-subtitle-1 mr-5"> {{ $t('common.webhookStatus') }}: </span>
          <alarm-webhook-execution-status
            :running="isExecutionRunning"
            :success="isExecutionSucceeded"
            :fail-reason="executionStatus.fail_reason"
          />
          <c-action-btn
            v-if="isExecutionSucceeded || isExecutionFailed"
            type="delete"
            @click="clearWebhookStatus"
          />
        </c-layout>
        <v-card v-if="isSomeOneWebhookStarted">
          <v-card-text>
            <slot :webhooks="executionStatus.webhooks" name="webhooks" />
          </v-card-text>
        </v-card>
      </c-layout>
    </v-expand-transition>
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import {
  isWebhookExecutionFailed,
  isWebhookExecutionRunning,
  isWebhookExecutionSucceeded,
  isWebhookExecutionWaiting,
} from '@/helpers/entities/webhook-execution/entity';

import AlarmWebhookExecutionStatus from './alarm-webhook-execution-status.vue';

export default {
  components: { AlarmWebhookExecutionStatus },
  emits: ['update:alarm', 'run:execution', 'clear:execution'],
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
    const isSomeOneWebhookStarted = computed(() => props.executionStatus?.webhooks.some(webhook => !isWebhookExecutionWaiting(webhook)),
    );

    const runTestExecution = () => emit('run:execution');
    const clearWebhookStatus = () => emit('clear:execution');

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
