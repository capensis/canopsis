<template lang="pug">
  g
    points-line(
      :points="line.points",
      cursor="move",
      stroke-width="10",
      pointer-events="stroke",
      @mousedown.stop="$emit('mousedown', $event)",
      @mouseup="$emit('mouseup', $event)"
    )
    template(v-if="selected")
      points-line(
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
        @mousedown.stop="onStartMovePoint(index, point)"
      )
</template>

<script>
import { getPointsWithGhosts } from '../../utils/points';

import PointsLine from '../common/points-line.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { PointsLine },
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
      movedPointIndex: undefined,
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
      this.editedPoints.splice(this.movedPointIndex, 1, { x, y });
    },

    addPointAfterIndex(index, { x, y }) {
      this.editedPoints.splice(index, 0, { x, y });
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
      this.$emit('resize', { points: this.editedPoints });

      this.$mouseMove.unregister(this.movePoint);
      this.$mouseUp.unregister(this.onMouseUp);
    },
  },
};
</script>
