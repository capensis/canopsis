<template>
  <v-layout class="pa-3 gap-3" column align-center>
    <h3>{{ $t('eventsRecording.subheader') }}</h3>
    <v-layout v-if="progress" class="gap-4" align-center>
      <span class="font-italic">{{ $t('eventsRecording.inProgress') }}</span>
      <v-progress-circular
        color="primary"
        width="3"
        indeterminate
      />
    </v-layout>
    <v-flex>
      <v-btn
        v-if="progress"
        key="stop"
        color="accent"
        @click="stop"
      >
        {{ $t('eventsRecording.stop') }}
      </v-btn>
      <v-btn
        v-else
        key="launch"
        color="warning"
        @click="launch"
      >
        {{ $t('eventsRecording.launch') }}
      </v-btn>
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  props: {
    progress: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, { emit }) {
    const launch = () => emit('launch');
    const stop = () => emit('stop');

    return {
      launch,
      stop,
    };
  },
};
</script>
