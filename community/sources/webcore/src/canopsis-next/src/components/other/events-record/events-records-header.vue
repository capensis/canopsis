<template>
  <v-layout class="pa-3 gap-3" column align-center>
    <h3>{{ $t('eventsRecord.subheader') }}</h3>
    <v-layout v-if="message" class="gap-4" align-center>
      <span class="font-italic">{{ message }}</span>
      <v-progress-circular
        color="primary"
        width="3"
        indeterminate
      />
    </v-layout>
    <v-flex>
      <v-btn
        v-if="current.is_resending"
        key="stop-resending"
        color="accent"
        @click="stopResending"
      >
        {{ $t('eventsRecord.stopResending') }}
      </v-btn>
      <v-btn
        v-else-if="current.is_recording"
        key="stop-recording"
        color="accent"
        @click="stopRecording"
      >
        {{ $t('eventsRecord.stop') }}
      </v-btn>
      <v-btn
        v-else
        key="start"
        :disabled="isLaunchDisabled"
        color="warning"
        @click="startRecording"
      >
        {{ $t('eventsRecord.launch') }}
      </v-btn>
    </v-flex>
    <span class="text-caption grey--text text--darken-1">
      <strong>{{ $t('common.limit') }}:</strong>&nbsp;
      {{ $t('eventsRecord.concurrentRecordings', { limit: eventRecordsConcurrentLimit }) }}
    </span>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useInfo } from '@/hooks/store/modules/info';

export default {
  props: {
    current: {
      type: Object,
      required: true,
    },
    inProgressCount: {
      type: Number,
      default: 0,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();
    const { eventRecordsConcurrentLimit } = useInfo();

    const message = computed(() => ({
      [props.current.is_recording]: t('eventsRecord.inProgress'),
      [props.current.is_resending]: t('eventsRecord.resendingInProgress'),
    }.true));

    const isLaunchDisabled = computed(() => {
      const limit = eventRecordsConcurrentLimit.value;

      return limit != null && limit > 0 && props.inProgressCount >= limit;
    });

    const startRecording = () => emit('start:recording');
    const stopRecording = () => emit('stop:recording');

    const stopResending = () => emit('stop:resending');

    return {
      message,
      eventRecordsConcurrentLimit,

      isLaunchDisabled,

      startRecording,
      stopRecording,

      stopResending,
    };
  },
};
</script>
