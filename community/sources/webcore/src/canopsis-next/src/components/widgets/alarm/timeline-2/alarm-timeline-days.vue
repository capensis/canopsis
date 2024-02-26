<template>
  <v-layout column>
    <v-layout
      v-for="(steps, date) in days"
      :key="date"
      class="my-1"
      column
    >
      <span class="grey--text my-2">{{ date }}</span>
      <alarm-timeline-days-item-steps :steps="steps" />
    </v-layout>
  </v-layout>
</template>

<script>
import AlarmTimelineDaysItemSteps from './alarm-timeline-days-item-steps.vue';

export default {
  components: { AlarmTimelineDaysItemSteps },
  props: {
    days: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
