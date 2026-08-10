<template>
  <span :class="chipClass" :style="chipStyle" class="c-state-count-changes-chip">
    <slot />
  </span>
</template>

<script>
import { computed } from 'vue';

import { ALARM_STATES, ALARM_STATES_CLASSES } from '@/constants';
import { getAlarmStateColor } from '@/helpers/entities/alarm/color';

export default {
  props: {
    state: {
      type: Number,
      default: ALARM_STATES.ok,
    },
  },

  setup(props) {
    const chipClass = computed(() => ALARM_STATES_CLASSES[props.state]);
    const chipStyle = computed(() => ({
      backgroundColor: getAlarmStateColor(props.state),
    }));

    return {
      chipClass,
      chipStyle,
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
  color: rgb(var(--v-theme-on-primary));
  width: 30px;
  height: 20px;
  border-radius: 10px;
}
</style>
