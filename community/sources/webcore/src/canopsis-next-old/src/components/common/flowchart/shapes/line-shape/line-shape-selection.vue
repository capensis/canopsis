<template>
  <g v-if="selected">
    <points-line-path
      :points="editingPoints"
      :type="shape.lineType"
      :stroke="color"
      fill="transparent"
      stroke-width="1"
      stroke-dasharray="4 4"
      pointer-events="none"
    />
    <circle
      v-for="(point, index) in editingPoints"
      :key="`${point._id}`"
      :cx="point.x"
      :cy="point.y"
      :fill="color"
      :r="cornerRadius"
      :pointer-events="moving ? 'none' : 'all'"
      cursor="crosshair"
      @mousedown.stop="onStartMovePoint(index)"
    />
  </g>
</template>

<script>
import { cloneDeep } from 'lodash';

import PointsLinePath from '../../common/points-line-path.vue';

export default {
  inject: ['$flowchart'],
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
    movePoint({ cursor }) {
      const point = this.editingPoints[this.movingPointIndex];

      point.x = cursor.x;
      point.y = cursor.y;
    },

    onStartMovePoint(index) {
      this.moving = true;

      this.movingPointIndex = index;

      this.$flowchart.on('mousemove', this.movePoint);
      this.$flowchart.on('mouseup', this.finishMovePoints);

      this.$emit('edit:point', this.editingPoints[this.movingPointIndex]);
    },

    finishMovePoints() {
      this.$emit('update', { points: this.editingPoints });

      this.$flowchart.off('mousemove', this.movePoint);
      this.$flowchart.off('mouseup', this.finishMovePoints);

      this.moving = false;
    },
  },
};
</script>
