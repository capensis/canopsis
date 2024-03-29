<template>
  <v-layout class="gap-2" column>
    <alarm-timeline-step
      v-for="step in steps"
      :key="step.key"
      :step="step"
    />
  </v-layout>
</template>

<script>
import AlarmTimelineStep from './alarm-timeline-step.vue';

export default {
  components: { AlarmTimelineStep },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
