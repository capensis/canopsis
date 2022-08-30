<template lang="pug">
  g
    template(v-if="selected")
      points-line-path(
        :points="editingPoints",
        :type="shape.lineType",
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

import PointsLinePath from '../common/points-line-path.vue';

export default {
  inject: ['$mouseMove', '$mouseUp'],
  components: { PointsLinePath },
  props: {
    shape: {
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
    'shape.points': {
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
