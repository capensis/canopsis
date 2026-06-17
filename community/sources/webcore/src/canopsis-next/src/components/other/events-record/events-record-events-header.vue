<template>
  <v-layout class="gap-3 my-4" justify-center align-center>
    <v-progress-circular
      v-if="isRecording"
      size="24"
      width="2"
      indeterminate
    />
    <span class="text-subtitle-2">
      {{ $t('modals.eventsRecord.subtitle', { count }) }}
    </span>
    <c-action-btn
      v-if="isRecording"
      :tooltip="$t('eventsRecord.stop')"
      icon="stop"
      color="blue darken-3"
      @click="$emit('stop:recording')"
    />
    <c-action-btn
      v-else
      :tooltip="$t('modals.eventsRecord.buttonTooltip')"
      icon="delete"
      color="error"
      @click="$emit('remove')"
    />
  </v-layout>
</template>

<script>
export default {
  props: {
    count: {
      type: Number,
      default: 0,
    },
    isRecording: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
