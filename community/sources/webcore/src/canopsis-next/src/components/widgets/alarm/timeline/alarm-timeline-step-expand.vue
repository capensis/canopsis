<template>
  <v-expand-transition>
    <v-layout v-if="expanded" class="gap-2" column>
      <alarm-timeline-step
        v-for="childrenStep in childrenSteps"
        :key="childrenStep.key"
        :step="childrenStep"
        deep
      />
    </v-layout>
  </v-expand-transition>
</template>

<script>
import { computed } from 'vue';

import { addKeyInEntities } from '@/helpers/array';

import AlarmTimelineStep from './alarm-timeline-step.vue';

export default {
  components: { AlarmTimelineStep },
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
    expanded: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const childrenSteps = computed(() => addKeyInEntities(props.step.steps ?? []));

    return { childrenSteps };
  },
};
</script>
