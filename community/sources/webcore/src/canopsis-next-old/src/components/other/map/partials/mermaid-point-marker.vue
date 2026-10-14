<template>
  <point-icon
    :class="{ 'mermaid-point--centered': isFullEntity }"
    :style="{ top: `${y}px`, left: `${x}px` }"
    :entity="entity"
    :size="size"
    :color-indicator="colorIndicator"
    :pbehavior-enabled="pbehaviorEnabled"
    class="mermaid-point"
    v-on="$listeners"
  />
</template>

<script>
import { isNumber } from 'lodash';

import PointIcon from './point-icon.vue';

export default {
  components: { PointIcon },
  props: {
    x: {
      type: Number,
      required: true,
    },
    y: {
      type: Number,
      required: true,
    },
    entity: {
      type: [String, Object],
      required: false,
    },
    size: {
      type: Number,
      default: 24,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    pbehaviorEnabled: {
      type: Boolean,
      required: false,
    },
  },
  computed: {
    isFullEntity() {
      return isNumber(this.entity?.state);
    },
  },
};
</script>

<style lang="scss">
.mermaid-point {
  position: absolute;
  transform: translate(-50%, -100%);
  user-select: none;
  cursor: pointer;
  transition: none;
  pointer-events: initial;

  &--centered {
    transform: translate(-50%, -50%);
  }
}
</style>
