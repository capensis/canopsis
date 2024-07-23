<template>
  <v-layout class="pa-3 gap-3" column align-center>
    <h3>{{ $t('eventsRecord.subheader') }}</h3>
    <v-layout v-if="progress" class="gap-4" align-center>
      <span class="font-italic">{{ $t('eventsRecord.inProgress') }}</span>
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
        {{ $t('eventsRecord.stop') }}
      </v-btn>
      <v-btn
        v-else
        key="start"
        color="warning"
        @click="start"
      >
        {{ $t('eventsRecord.launch') }}
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
    const start = () => emit('start');
    const stop = () => emit('stop');

    return {
      start,
      stop,
    };
  },
};
</script>
