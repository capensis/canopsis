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
        :points="editedPoints",
        :stroke="color",
        fill="transparent",
        stroke-width="1",
        stroke-dasharray="4 4",
        pointer-events="none"
      )
      circle(
        v-for="(point, index) in points",
        :key="index",
        :cx="point.x",
        :cy="point.y",
        :fill="color",
        :r="cornerRadius",
        :opacity="point.ghost ? 0.4 : 1",
        cursor="crosshair",
        :pointer-events="moving ? 'none' : 'all'",
        @mousedown.stop="onStartMovePoint(index, point)"
      )
</template>

<script>
import { getPointsWithGhosts } from '../../utils/points';

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
      editedPoints: [],
      moving: undefined,
      movingPointIndex: undefined,
    };
  },
  computed: {
    points() {
      return getPointsWithGhosts(this.editedPoints);
    },
  },
  watch: {
    'line.points': {
      immediate: true,
      handler(points) {
        this.editedPoints = [...points];
      },
    },
  },
  methods: {
    movePoint({ x, y }) {
      this.editedPoints.splice(this.movingPointIndex, 1, { x, y });
    },

    addPointAfterIndex(index, { x, y }) {
      this.editedPoints.splice(index, 0, { x, y });
    },

    onStartMovePoint(index, point) {
      this.moving = true;

      if (point.ghost) {
        this.addPointAfterIndex(point.index, point);
      }

      this.movingPointIndex = point.index;

      this.$mouseMove.register(this.movePoint);
      this.$mouseUp.register(this.finishMovePoints);
    },

    finishMovePoints() {
      this.$emit('resize', { points: this.editedPoints });

      this.$mouseMove.unregister(this.movePoint);
      this.$mouseUp.unregister(this.onMouseUp);

      this.moving = false;
    },
  },
};
</script>
