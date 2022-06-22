<template lang="pug">
  g
    slot
    points-path(
      v-bind="shape.style",
      :points="shape.points",
      :marker-end="markerEnd",
      :marker-start="markerStart",
      pointer-events="none"
    )
    line-shape-selection(
      v-if="!readonly",
      :selected="selected",
      :line="shape",
      @mousedown="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)",
      @resize="onResize",
      @edit:point="$emit('edit:point', $event)"
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
    readonly: {
      type: Boolean,
      default: false,
    },
    markerStart: {
      type: String,
      required: false,
    },
    markerEnd: {
      type: String,
      required: false,
    },
  },
  methods: {
    move(newOffset, oldOffset) {
      const { points } = this.shape;

      this.updateModel({
        ...this.shape,

        points: points.map(point => ({
          ...point,
          x: (point.x - oldOffset.x) + newOffset.x,
          y: (point.y - oldOffset.y) + newOffset.y,
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
