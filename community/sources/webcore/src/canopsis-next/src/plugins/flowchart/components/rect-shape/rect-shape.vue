<template lang="pug">
  g(
    cursor="move",
    @mousedown="$emit('mousedown', $event)",
    @mouseup="$emit('mouseup', $event)",
    @dblclick.stop="",
    @click.stop=""
  )
    rect(
      v-bind="shape.style",
      :x="shape.x",
      :y="shape.y",
      :width="shape.width",
      :height="shape.height",
      pointer-events="all"
    )
    rect-shape-selection(
      v-if="selected",
      :rect="shape",
      @resize="onResize"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import RectShapeSelection from './rect-shape-selection.vue';

export default {
  components: { RectShapeSelection },
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
  methods: {
    move(newOffset, oldOffset) {
      const { x, y } = this.shape;

      this.updateModel({
        ...this.shape,

        x: (x - oldOffset.x) + newOffset.x,
        y: (y - oldOffset.y) + newOffset.y,
      });
    },

    onResize(rect) {
      this.updateModel({ ...this.shape, ...rect });
    },
  },
};
</script>
