<template lang="pug">
  g(
    cursor="move",
    @mousedown="$emit('mousedown', $event)",
    @mouseup="$emit('mouseup', $event)",
    @dblclick.stop="",
    @click.stop=""
  )
    points-path(
      v-bind="shape.style",
      :points="shape.points"
    )
    points-path(
      :points="shape.points",
      stroke-width="10"
    )
    line-shape-selection(
      v-if="selected",
      :line="shape",
      @resize="onResize"
    )
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import PointsPath from '../common/points-path.vue';

import LineShapeSelection from './line-shape-selection.vue';

export default {
  components: { PointsPath, LineShapeSelection },
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
