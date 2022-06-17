<template lang="pug">
  g
    points-line(v-bind="shape.style", :points="shape.points", pointer-events="none")
    line-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :line="shape",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)",
      @resize="onResize"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import PointsLine from '../common/points-line.vue';

import LineShapeSelection from './line-shape-selection.vue';

export default {
  components: { PointsLine, LineShapeSelection },
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
    readonly: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    move(newOffset, oldOffset) {
      const { points } = this.shape;

      this.updateModel({
        ...this.shape,

        points: points.map(({ x, y }) => ({
          x: (x - oldOffset.x) + newOffset.x,
          y: (y - oldOffset.y) + newOffset.y,
        })),
      });
    },

    onResize(line) {
      this.updateModel({
        ...this.shape,
        ...line,
      });
    },
  },
};
</script>
