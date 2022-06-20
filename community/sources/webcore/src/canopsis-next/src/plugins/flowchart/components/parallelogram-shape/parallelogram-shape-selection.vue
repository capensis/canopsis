<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    parallelogram-figure(
      :x="parallelogram.x",
      :y="parallelogram.y",
      :width="parallelogram.width",
      :height="parallelogram.height",
      :pointer-events="pointerEvents",
      fill="transparent",
      cursor="move",
      @mousedown.stop="$listeners.mousedown",
      @mouseup="$listeners.mouseup"
    )
    rect-selection(
      v-if="selected",
      :rect="parallelogram",
      :padding="padding",
      :color="color",
      :corner-radius="cornerRadius",
      resizable,
      @start:resize="startResize"
    )
</template>

<script>
import { resizeRectangleShapeByDirection } from '../../utils/resize';

import RectSelection from '../common/rect-selection.vue';
import ParallelogramFigure from '../common/parallelogram-figure.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { ParallelogramFigure, RectSelection },
  props: {
    parallelogram: {
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
      const { rect, direction } = resizeRectangleShapeByDirection(
        this.parallelogram,
        cursor,
        this.direction,
      );

      this.direction = direction;
      this.$emit('resize', rect);
    },
  },
};
</script>
