<template lang="pug">
  g(
    cursor="move",
    @mousedown="$emit('mousedown', $event)",
    @mouseup="$emit('mouseup', $event)",
    @dblclick.stop="",
    @click.stop=""
  )
    marker#arrowLineTriangle(
      refX="20",
      refY="10",
      markerWidth="60",
      markerHeight="60",
      markerUnits="userSpaceOnUse",
      orient="auto"
    )
      arrow-shape(:fill="shape.style['stroke']")
    points-path(
      v-bind="shape.style",
      :points="shape.points",
      marker-end="url(#arrowLineTriangle)",
      pointer-events="all"
    )
    points-path(
      :points="shape.points",
      stroke-width="10",
      pointer-events="all"
    )
    line-shape-selection(
      v-if="selected",
      :line="shape",
      @resize="onResize"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import LineShapeSelection from '../line-shape/line-shape-selection.vue';

import ArrowShape from '../common/arrow.vue';
import PointsPath from '../common/points-path.vue';

export default {
  components: { LineShapeSelection, PointsPath, ArrowShape },
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
