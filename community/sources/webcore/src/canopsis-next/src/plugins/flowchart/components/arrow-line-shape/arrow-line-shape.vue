<template lang="pug">
  g
    marker#arrow-line-triangle(
      refX="20",
      refY="10",
      markerWidth="60",
      markerHeight="60",
      markerUnits="userSpaceOnUse",
      orient="auto"
    )
      arrow(:fill="shape.style['stroke']")
    points-line(
      v-bind="shape.style",
      :points="shape.points",
      pointer-events="none",
      marker-end="url(#arrow-line-triangle)"
    )
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

import LineShapeSelection from '../line-shape/line-shape-selection.vue';

import Arrow from '../common/arrow.vue';
import PointsLine from '../common/points-line.vue';

export default {
  components: { LineShapeSelection, PointsLine, Arrow },
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
