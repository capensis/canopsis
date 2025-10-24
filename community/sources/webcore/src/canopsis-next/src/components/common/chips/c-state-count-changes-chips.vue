<template>
  <span class="d-inline-flex align-center">
    <c-state-count-changes-chip :color="color">{{ affect.from }}</c-state-count-changes-chip>
    <v-icon size="16">arrow_forward</v-icon>
    <c-state-count-changes-chip :color="color">{{ affect.to }}</c-state-count-changes-chip></span>
</template>

<script>
import { CSS_COLORS_VARS } from '@/config';

export default {
  props: {
    affect: {
      type: Object,
      default: () => ({ from: 0, to: 0 }),
    },
    color: {
      type: String,
      default: CSS_COLORS_VARS.primary,
    },
  },
};
</script>
