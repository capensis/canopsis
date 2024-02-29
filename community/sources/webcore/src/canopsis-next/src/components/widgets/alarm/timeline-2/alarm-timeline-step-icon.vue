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
