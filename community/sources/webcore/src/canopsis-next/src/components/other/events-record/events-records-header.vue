<template>
  <v-layout class="pa-3 gap-3" column align-center>
    <h3>{{ $t('eventsRecord.subheader') }}</h3>
    <v-layout
      v-for="message in messages"
      :key="message"
      class="gap-4"
      align-center
    >
      <span class="font-italic">{{ message }}</span>
      <v-progress-circular
        color="primary"
        width="3"
        indeterminate
      />
    </v-layout>
    <v-flex>
      <v-btn
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
      {{ $t('eventsRecord.concurrentRecordings', { limit }) }}
    </span>
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    recordings: {
      type: Array,
      default: () => [],
    },
    resendings: {
      type: Array,
      default: () => [],
    },
    limit: {
      type: Number,
      default: 1,
    },
  },
  setup(props, { emit }) {
    const { t } = useI18n();

    const messages = computed(() => {
      const result = [];

      if (props.recordings?.length) {
        result.push(t('eventsRecord.inProgress'));
      }

      if (props.resendings?.length) {
        result.push(t('eventsRecord.resendingInProgress'));
      }

      return result;
    });

    const isLaunchDisabled = computed(() => !props.limit || props.recordings?.length >= props.limit);

    const startRecording = () => emit('start:recording');

    return {
      messages,
      isLaunchDisabled,
      startRecording,
    };
  },
};
</script>
