<template>
  <v-layout class="gap-2" column>
    <template v-for="step in steps">
      <alarm-timeline-step
        :key="step.key"
        :step="step"
        @expand="toggleExpandedForStep(step.key, $event)"
      />
      <alarm-timeline-step-expand
        :key="`${step.key}-expand`"
        :step="step"
        :expanded="expandedSteps[step.key]"
      />
    </template>
  </v-layout>
</template>

<script>
import { reactive, set } from 'vue';

import AlarmTimelineStep from './alarm-timeline-step.vue';
import AlarmTimelineStepExpand from './alarm-timeline-step-expand.vue';

export default {
  components: { AlarmTimelineStep, AlarmTimelineStepExpand },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
  setup() {
    const expandedSteps = reactive({});
    const toggleExpandedForStep = (key, value) => set(expandedSteps, key, value);

    return {
      expandedSteps,

      toggleExpandedForStep,
    };
  },
};
</script>
