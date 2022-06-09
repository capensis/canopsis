<template lang="pug">
  g(
    cursor="move",
    @mousedown="$emit('mousedown', $event)",
    @mouseup="$emit('mouseup', $event)",
    @dblclick.stop="",
    @click.stop=""
  )
    circle(
      v-bind="shape.style",
      :cx="centerX",
      :cy="centerY",
      :r="radius",
      pointer-events="all"
    )
    square-shape-selection(
      v-if="selected",
      :square="selection",
      @resize="onResize"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import SquareShapeSelection from '../square-shape/square-shape-selection.vue';

export default {
  components: { SquareShapeSelection },
  mixins: [formBaseMixin],
  model: {
    prop: 'shape',
    event: 'input',
  },
  props: {
    shape: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    cornerOffset: {
      type: Number,
      default: 0,
    },
  },
  computed: {
    radius() {
      return this.shape.diameter / 2;
    },

    centerX() {
      return this.shape.x + this.radius;
    },

    centerY() {
      return this.shape.y + this.radius;
    },

    selection() {
      return {
        x: this.shape.x,
        y: this.shape.y,
        size: this.shape.diameter,
      };
    },
  },
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.shape;

      this.updateModel({
        ...this.shape,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(square) {
      this.updateModel({
        ...this.shape,
        x: square.x,
        y: square.y,
        diameter: square.size,
      });
    },
  },
};
</script>
