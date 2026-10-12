<template>
  <v-badge
    :value="isFailedStep"
    class="time-line-flag"
    color="transparent"
    overlap
  >
    <template #badge="">
      <v-icon
        class="time-line-flag__badge-icon"
        color="error"
        size="14"
      >
        error
      </v-icon>
    </template>
    <v-icon :color="stepColor">
      {{ stepIcon }}
    </v-icon>
  </v-badge>
</template>

<script>
import { computed } from 'vue';

import { getAlarmStepColor, getAlarmStepIcon } from '@/helpers/entities/alarm/step/formatting';
import { isFailStepType } from '@/helpers/entities/alarm/step/entity';

/**
 * Component for the flag on the alarms list's timeline
 *
 * @module alarm
 *
 * @prop {Object} step - step object
 */
export default {
  props: {
    step: {
      type: Object,
      required: true,
    },
    error: {
      type: Boolean,
      default: false,
    },
  },
  setup(props) {
    const stepIcon = computed(() => getAlarmStepIcon(props.step._t));
    const stepColor = computed(() => getAlarmStepColor(props.step._t));
    const isFailedStep = computed(() => isFailStepType(props.step._t));

    return {
      stepIcon,
      stepColor,
      isFailedStep,
    };
  },
};
</script>

<style lang="scss" scoped>
.time-line-flag {
  &__badge-icon {
    width: 14px !important;
    max-width: unset !important;
    height: 14px !important;
    background: white;
    border-radius: 50%;
  }
}
</style>
