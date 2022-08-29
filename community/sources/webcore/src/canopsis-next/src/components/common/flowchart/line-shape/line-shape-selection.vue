<template lang="pug">
  g(@dblclick="$emit('dblclick', $event)")
    points-path(
      :points="line.points",
      :type="line.lineType",
      :pointer-events="moving ? 'none' : 'stroke'",
      cursor="move",
      stroke-width="10",
      @mousedown.stop="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
    template(v-if="selected")
      points-path(
        :points="editingPoints",
        :type="line.lineType",
        :stroke="color",
        fill="transparent",
        stroke-width="1",
        stroke-dasharray="4 4",
        pointer-events="none"
      )
      circle(
        v-for="(point, index) in editingPoints",
        :key="`${point._id}`",
        :cx="point.x",
        :cy="point.y",
        :fill="color",
        :r="cornerRadius",
        cursor="crosshair",
        :pointer-events="moving ? 'none' : 'all'",
        @mousedown.stop="onStartMovePoint(index)"
      )
</template>

<script>
import { cloneDeep } from 'lodash';

import PointsPath from '../common/points-path.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { PointsPath },
  props: {
    line: {
      type: Object,
      required: true,
    },
    selected: {
      type: Boolean,
      default: false,
    },
    color: {
      type: String,
      default: 'blue',
    },
    cornerRadius: {
      type: Number,
      default: 4,
    },
  },
  data() {
    return {
      editingPoints: [],
      moving: undefined,
      movingPointIndex: undefined,
    };
  },
  watch: {
    'line.points': {
      immediate: true,
      deep: true,
      handler(points) {
        this.editingPoints = cloneDeep(points);
      },
    },
  },
  methods: {
    movePoint({ x, y }) {
      const point = this.editingPoints[this.movingPointIndex];

      point.x = x;
      point.y = y;
    },

    onStartMovePoint(index) {
      this.moving = true;

      this.movingPointIndex = index;

      this.$mouseMove.register(this.movePoint);
      this.$mouseUp.register(this.finishMovePoints);

      this.$emit('edit:point', this.editingPoints[this.movingPointIndex]);
    },

    finishMovePoints() {
      this.$emit('update', { points: this.editingPoints });

      this.$mouseMove.unregister(this.movePoint);
      this.$mouseUp.unregister(this.onMouseUp);

      this.moving = false;
    },
  },
};
</script>
