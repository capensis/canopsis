<template>
  <span
    :style="style"
    class="c-state-count-changes-chip"
  >
    <slot />
  </span>
</template>

<script>
import { CSS_COLORS_VARS } from '@/config';

export default {
  props: {
    color: {
      type: String,
      default: CSS_COLORS_VARS.primary,
    },
  },
  computed: {
    style() {
      return { backgroundColor: this.color };
    },
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
