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
      @mousedown="$listeners.mousedown",
      @mouseup="$listeners.mouseup",
      @update="$listeners.update",
      @edit:point="$listeners['edit:point']"
    )
</template>

<script>
import PointsPath from '../common/points-path.vue';

import LineShapeSelection from './line-shape-selection.vue';

export default {
  components: { PointsPath, LineShapeSelection },
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
};
</script>
