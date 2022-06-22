<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    rect(
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      cursor="move",
      @start:resize="startResize"
    )
    rect-connectors(
      v-if="connecting",
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :padding="padding",
      :color="color",
      @connected="$emit('connected', $event)",
      @connecting="$emit('connecting', $event)",
      @unconnect="$emit('unconnect')"
    )
</template>

<script>
import { resizeRectangleShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';
import RectConnectors from '../common/rect-connectors.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RectConnectors, RectSelection },
  props: {
    rect: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    padding: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'blue',
    },
    pointerEvents: {
      type: String,
      default: 'all',
    },
    cornerRadius: {
      type: Number,
      default: 4,
    },
    connecting: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      direction: undefined,
    };
  },
  methods: {
    startResize(direction) {
      this.$mouseMove.register(this.onResize);
      this.$mouseUp.register(this.finishResize);

      this.direction = direction;
    },

    finishResize() {
      this.$mouseMove.unregister(this.onResize);
      this.$mouseUp.unregister(this.finishResize);
    },

    onResize(cursor) {
      const { direction, rect } = resizeRectangleShapeByDirection(
        this.rect,
        cursor,
        this.direction,
      );

      this.direction = direction;

      this.$emit('resize', rect);
    },
  },
};
</script>
