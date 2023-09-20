<template lang="pug">
  span.d-inline-flex.align-center
    c-state-count-changes-chip(:color="color") {{ affect.from }}
    v-icon(size="16") arrow_forward
    c-state-count-changes-chip(:color="color") {{ affect.to }}
</template>

<script>
import { CSS_COLOR_VARS } from '@/config';

export default {
  props: {
    affect: {
      type: Object,
      default: () => ({ from: 0, to: 0 }),
    },
    color: {
      type: String,
      default: CSS_COLOR_VARS.primary,
    },
  },
};
</script>
