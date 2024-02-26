<template>
  <v-layout column>
    <alarm-timeline-days-item-steps :steps="steps" />
  </v-layout>
</template>

<script>
import AlarmTimelineDaysItemSteps from './alarm-timeline-days-item-steps.vue';

export default {
  components: { AlarmTimelineDaysItemSteps },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
