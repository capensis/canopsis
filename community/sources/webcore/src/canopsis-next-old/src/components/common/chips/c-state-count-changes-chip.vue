<template>
  <span
    :class="chipClass"
    class="c-state-count-changes-chip"
  >
    <slot />
  </span>
</template>

<script>
import { computed } from 'vue';

import { ALARM_STATES, ALARM_STATES_CLASSES } from '@/constants';

export default {
  props: {
    state: {
      type: Number,
      default: ALARM_STATES.ok,
    },
  },
  setup(props) {
    const chipClass = computed(() => ALARM_STATES_CLASSES[props.state]);

    return {
      chipClass,
    };
  },
};
</script>

<style lang="scss">
.c-state-count-changes-chip {
  font-size: 12px;
  display: inline-flex;
  justify-content: center;
  align-items: center;
  width: 30px;
  height: 20px;
  border-radius: 10px;
}
</style>
