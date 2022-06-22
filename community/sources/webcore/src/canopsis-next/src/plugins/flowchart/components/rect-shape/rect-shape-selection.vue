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
      :x="rect.x",
      :y="rect.y",
      :width="rect.width",
      :height="rect.height",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      :connectable="connection",
      :connectors="rect.connectors",
      :resizable="selected && !connection",
      cursor="move",
      @start:resize="startResize"
    )
</template>

<script>
import { resizeRectangleShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RectSelection },
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
  },
  data() {
    return {
      direction: undefined,
      connection: false,
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
