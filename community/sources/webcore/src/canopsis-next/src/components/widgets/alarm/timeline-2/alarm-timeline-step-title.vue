<template>
  <div v-if="!isStateCounter">
    <c-alarm-chip
      v-if="isChangeState"
      :value="step.val"
    />
  </div>
</template>

<script>
import { ALARM_LIST_STEPS } from '@/constants';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    message() {
      switch (this.step._t) {
        case ALARM_LIST_STEPS.stateCounter:
          return {
            title: 'Cropped severity (since the last change of the status)',
          };
        case ALARM_LIST_STEPS.changeState:
          return {
            icon: 'c-alarm-chip',
            title: 'State changed by %username%',
          };
        case ALARM_LIST_STEPS.stateinc:
          return {
            icon: 'c-alarm-chip',
            title: 'State increased by centreon.centreon_0',
          };
        case ALARM_LIST_STEPS.statedec:
          return {
            icon: 'c-alarm-chip',
            title: 'State decreased by centreon.centreon_0',
          };

        case ALARM_LIST_STEPS.statusinc:
          return {
            icon: 'swap_vert',
            title: 'Status changed to flapping by the system',
          };

        default:
          return {};
      }
    },

    isStateCounter() {
      return this.step._t === ALARM_LIST_STEPS.stateCounter;
    },

    isChangeState() {
      return this.step._t === ALARM_LIST_STEPS.changeState;
    },

    isChangeStatus() {
      return [ALARM_LIST_STEPS.statusinc, ALARM_LIST_STEPS.statusdec].includes(this.step._t);
    },
  },
};
</script>
