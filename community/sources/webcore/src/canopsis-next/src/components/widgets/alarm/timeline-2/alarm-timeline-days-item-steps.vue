<template>
  <v-layout column>
    <alarm-timeline-days-item-steps-item
      v-for="step in steps"
      :key="step.key"
      :step="step"
      class="my-2"
    />
  </v-layout>
</template>

<script>
import AlarmTimelineDaysItemStepsItem from './alarm-timeline-days-item-steps-item.vue';

export default {
  components: { AlarmTimelineDaysItemStepsItem },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
};
</script>
