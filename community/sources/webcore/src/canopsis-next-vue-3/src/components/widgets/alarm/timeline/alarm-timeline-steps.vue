<template>
  <c-layout gap="2" column>
    <template v-for="step in steps" :key="step._id">
      <alarm-timeline-step
        :step="step"
        :is-html-enabled="isHtmlEnabled"
        @expand="toggleExpandedForStep(step._id, $event)"
      />
      <alarm-timeline-step-expand
        :step="step"
        :expanded="expandedSteps[step._id]"
      />
    </template>
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import AlarmTimelineStepExpand from './alarm-timeline-step-expand.vue';
import AlarmTimelineStep from './alarm-timeline-step.vue';

export default {
  components: { AlarmTimelineStep, AlarmTimelineStepExpand },

  props: {
    steps: {
      type: Array,
      default: () => [],
    },
    isHtmlEnabled: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const expandedSteps = ref({});

    const toggleExpandedForStep = (key, value) => {
      expandedSteps.value = {
        ...expandedSteps.value,
        [key]: value,
      };
    };

    return {
      expandedSteps,
      toggleExpandedForStep,
    };
  },
};
</script>
