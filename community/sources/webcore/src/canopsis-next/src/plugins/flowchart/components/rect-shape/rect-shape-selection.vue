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
      resizable,
      @start:resize="startResize"
    )
    rect-shape-connectors(
      v-if="!selected && connection",
      :x="rect.x",
      :y="rect.y",
      :rx="rect.rx",
      :ry="rect.ry",
      :width="rect.width",
      :height="rect.height"
    )
</template>

<script>
import { resizeRectangleShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';
import RectShapeConnectors from './rect-shape-connectors.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { RectShapeConnectors, RectSelection },
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
