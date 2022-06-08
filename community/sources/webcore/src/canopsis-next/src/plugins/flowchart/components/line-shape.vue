<template lang="pug">
  g(
    cursor="move",
    @mousedown="$emit('mousedown', $event)",
    @mouseup="$emit('mouseup', $event)",
    @dblclick.stop="",
    @click.stop=""
  )
    line(
      v-bind="shape.style",
      :x1="shape.x1",
      :y1="shape.y1",
      :x2="shape.x2",
      :y2="shape.y2",
      pointer-events="all"
    )
    line(
      :x1="shape.x1",
      :y1="shape.y1",
      :x2="shape.x2",
      :y2="shape.y2",
      stroke-width="10",
      pointer-events="all"
    )
    line-shape-selection(
      v-if="selected",
      :x1="shape.x1",
      :y1="shape.y1",
      :x2="shape.x2",
      :y2="shape.y2",
      @start:move-from="onStartMoveFrom",
      @start:move-to="onStartMoveTo"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import LineShapeSelection from './line-shape-selection.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { LineShapeSelection },
  mixins: [formBaseMixin],
  model: {
    value: 'shape',
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
  data() {
    return {
      firstPointMoving: false,
      secondPointMoving: false,
    };
  },
  created() {
    this.$mouseMove.register(this.onMouseMove);
    this.$mouseUp.register(this.onMouseUp);
  },
  beforeDestroy() {
    this.$mouseMove.unregister(this.onMouseMove);
    this.$mouseUp.unregister(this.onMouseUp);
  },
  methods: {
    move(newOffset, oldOffset) {
      const { x1, x2, y1, y2 } = this.shape;

      this.updateModel({
        ...this.shape,

        x1: (x1 - oldOffset.x) + newOffset.x,
        x2: (x2 - oldOffset.x) + newOffset.x,
        y1: (y1 - oldOffset.y) + newOffset.y,
        y2: (y2 - oldOffset.y) + newOffset.y,
      });
    },

    onMouseMove({ x, y }) {
      if (this.firstPointMoving) {
        this.updateModel({
          ...this.shape,

          x1: x,
          y1: y,
        });
      }

      if (this.secondPointMoving) {
        this.updateModel({
          ...this.shape,

          x2: x,
          y2: y,
        });
      }
    },

    onStartMoveFrom() {
      this.firstPointMoving = true;
    },

    onStartMoveTo() {
      this.secondPointMoving = true;
    },

    onMouseUp() {
      this.firstPointMoving = false;
      this.secondPointMoving = false;
    },
  },
};
</script>
