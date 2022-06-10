<template lang="pug">
  g
    rect(
      v-bind="square.style",
      :x="square.x",
      :y="square.y",
      :width="square.size",
      :height="square.size"
    )
    square-shape-selection(
      :selected="selected",
      :square="square",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)",
      @resize="onResize"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import SquareShapeSelection from './square-shape-selection.vue';

export default {
  components: { SquareShapeSelection },
  mixins: [formBaseMixin],
  model: {
    prop: 'square',
    event: 'input',
  },
  props: {
    square: {
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
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.square;

      this.updateModel({
        ...this.square,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(square) {
      this.updateModel({
        ...this.square,
        x: square.x,
        y: square.y,
        size: square.size,
      });
    },
  },
};
</script>
