<template lang="pug">
  g
    points-path(
      :points="line.points",
      cursor="move",
      stroke-width="10",
      :pointer-events="moving ? 'none' : 'stroke'",
      @mousedown.stop="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
    template(v-if="selected")
      points-path(
        :points="editingPoints",
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
      circle(
        v-for="(point, index) in ghostPoints",
        :key="`${point._id}_ghost`",
        :cx="point.x",
        :cy="point.y",
        :fill="color",
        :r="cornerRadius",
        :opacity="0.4",
        cursor="crosshair",
        :pointer-events="moving ? 'none' : 'all'",
        @mousedown.stop="onStartGhostMovePoint(index, point)"
      )
</template>

<script>
import { cloneDeep } from 'lodash';
import { getGhostPoints } from '@/helpers/flowchart/points';

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
  computed: {
    ghostPoints() {
      return getGhostPoints(this.editingPoints);
    },
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

    onStartGhostMovePoint(index, point) {
      const newPointIndex = index + 1;

      this.editingPoints.splice(newPointIndex, 0, point);

      this.onStartMovePoint(newPointIndex);
    },

    onStartMovePoint(index) {
      this.moving = true;

      this.movingPointIndex = index;

      this.$mouseMove.register(this.movePoint);
      this.$mouseUp.register(this.finishMovePoints);

      this.$emit('edit:point', this.editingPoints[this.movingPointIndex]);
    },

    finishMovePoints() {
      this.$emit('resize', { points: this.editingPoints });

      this.$mouseMove.unregister(this.movePoint);
      this.$mouseUp.unregister(this.onMouseUp);

      this.moving = false;
    },
  },
};
</script>
