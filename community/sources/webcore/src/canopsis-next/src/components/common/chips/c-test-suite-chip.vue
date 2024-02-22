<template>
  <div
    :style="{ background: color }"
    class="test-suite-chip white--text"
    v-on="$listeners"
  >
    <slot>{{ $t(`testSuite.statuses.${value}`) }}</slot>
  </div>
</template>

<script>
import { TEST_SUITE_COLORS, TEST_SUITE_STATUSES } from '@/constants';

export default {
  props: {
    value: {
      type: [String, Number],
      default: TEST_SUITE_STATUSES.passed,
    },
  },
  computed: {
    color() {
      return TEST_SUITE_COLORS[this.value];
    },
  },
};
</script>

<style lang="scss">
.test-suite-chip {
  font-size: 12px;
  padding: 3px 7px;
  display: inline-block;
  border-radius: 12px;
}
</style>
