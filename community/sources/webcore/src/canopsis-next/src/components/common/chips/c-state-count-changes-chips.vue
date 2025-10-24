<template>
  <span class="d-inline-flex align-center">
    <c-state-count-changes-chip :state="state">{{ affect.from }}</c-state-count-changes-chip>
    <v-icon size="16">arrow_forward</v-icon>
    <c-state-count-changes-chip :state="state">{{ affect.to }}</c-state-count-changes-chip></span>
</template>

<script>
import { ALARM_STATES } from '@/constants';

export default {
  props: {
    affect: {
      type: Object,
      default: () => ({ from: 0, to: 0 }),
    },
    state: {
      type: Number,
      default: ALARM_STATES.ok,
    },
  },
};
</script>
