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
        v-if="current.isResending"
        key="stop"
        color="accent"
        @click="stopResending"
      >
        {{ $t('eventsRecord.stopResending') }}
      </v-btn>
      <v-btn
        v-else-if="current.isRecording"
        key="stop"
        color="accent"
        @click="stopRecording"
      >
        {{ $t('eventsRecord.stop') }}
      </v-btn>
      <v-btn
        v-else
        key="start"
        color="warning"
        @click="startRecording"
      >
        {{ $t('eventsRecord.launch') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    current: {
      type: Object,
      required: true,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const message = computed(() => ({
      [props.current.isRecording]: t('eventsRecord.inProgress'),
      [props.current.isResending]: t('eventsRecord.resendingInProgress'),
    }.true));

    const startRecording = () => emit('start:recording');
    const stopRecording = () => emit('stop:recording');

    const stopResending = () => emit('stop:resending');

    return {
      message,

      startRecording,
      stopRecording,

      stopResending,
    };
  },
};
</script>
