<template>
  <v-layout column>
    <v-layout
      v-for="(steps, date) in days"
      :key="date"
      class="my-1"
      column
    >
      <span class="grey--text my-2">{{ date }}</span>
      <alarm-timeline-steps :steps="steps" :is-html-enabled="isHtmlEnabled" />
    </v-layout>
  </v-layout>
</template>

<script>
import AlarmTimelineSteps from './alarm-timeline-steps.vue';

export default {
  components: { AlarmTimelineSteps },
  props: {
    days: {
      type: Object,
      default: () => ({}),
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
