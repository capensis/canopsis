<template lang="pug">
  span.c-horizontal-time-line-card
    v-tooltip(:disabled="!step.m", top)
      v-layout(slot="activator", align-center, column)
        c-alarm-chip(v-if="isState", :value="step.val")
        v-icon(v-else, :style="{ color: style.icon }") {{ style.icon }}
        span.c-horizontal-time-line-card__time {{ step.t | date('time') }}
      span(v-html="step.m")
</template>

<script>
import { formatStep } from '@/helpers/formatting';

export default {
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  computed: {
    isState() {
      return this.step._t.startsWith('state');
    },

    style() {
      return formatStep(this.step);
    },
  },
};
</script>

<style lang="scss" scoped>
.c-horizontal-time-line-card {
  display: inline-flex;
  align-items: center;

  &__time {
    font-size: 10px;
  }
}
</style>
