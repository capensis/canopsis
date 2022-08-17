<template lang="pug">
  point-icon.mermaid-point(
    v-on="$listeners",
    :class="{ 'mermaid-point--centered': isEntityObject }",
    :style="{ top: `${y}px`, left: `${x}px` }",
    :entity="entity",
    :size="size",
    :color-indicator="colorIndicator"
  )
</template>

<script>
import { isObject } from 'lodash';

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
  },
  computed: {
    isEntityObject() {
      return isObject(this.entity);
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

  &--centered {
    transform: translate(-50%, -50%);
  }
}
</style>
