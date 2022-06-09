<template lang="pug">
  g
    points-path(
      :points="line.points",
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
      movedPointIndex: undefined,
    };
  },
  computed: {
    points() {
      return getPointsWithGhosts(this.line.points);
    },
  },
  methods: {
    movePoint({ x, y }) {
      const points = this.line.points.slice();

      points.splice(this.movedPointIndex, 1, { x, y });

      this.$emit('resize', { points });
    },

    addPointAfterIndex(index, { x, y }) {
      const points = this.line.points.slice();

      points.splice(index, 0, { x, y });

      this.$emit('resize', { points });
    },

    onStartMovePoint(index, point) {
      if (point.ghost) {
        this.addPointAfterIndex(point.index, point);
      }

      this.movedPointIndex = point.index;

      this.$mouseMove.register(this.movePoint);
      this.$mouseUp.register(this.finishMovePoints);
    },

    finishMovePoints() {
      this.$mouseMove.unregister(this.movePoint);
      this.$mouseUp.unregister(this.onMouseUp);
    },
  },
};
</script>
