<template lang="pug">
  span.d-inline-flex.align-center
    span.c-state-count-changes-chip(:style="style") {{ affect.from }}
    v-icon(size="16") arrow_forward
    span.c-state-count-changes-chip(:style="style") {{ affect.to }}
</template>

<script>
import { COLORS } from '@/config';

export default {
  props: {
    affect: {
      type: Object,
      default: () => ({ from: 0, to: 0 }),
    },
    color: {
      type: String,
      default: COLORS.primary,
    },
  },
  computed: {
    style() {
      return { backgroundColor: this.color };
    },
  },
};
</script>
