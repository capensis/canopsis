<template lang="pug">
  path(
    ref="path",
    :d="path",
    :fill="fill",
    pointer-events="all",
    v-on="$listeners"
  )
</template>

<script>
import { LINE_TYPES } from '@/constants';

import { calculateCenterBetweenPoint } from '@/helpers/flowchart/points';

export default {
  props: {
    points: {
      type: Array,
      required: true,
    },
    fill: {
      type: String,
      default: 'transparent',
    },
    type: {
      type: String,
      default: LINE_TYPES.sharp,
    },
  },
  computed: {
    sharpPath() {
      return this.points.map(({ x, y }, index) => `${index === 0 ? 'M' : 'L'} ${x} ${y}`).join(' ');
    },

    curvesPath() {
      const lastIndex = this.points.length - 1;

      return this.points.reduce((acc, point, index) => {
        if (!index) {
          return `M ${point.x} ${point.y}`;
        }

        if (index === lastIndex) {
          return this.points.length === 2
            ? `${acc}L ${point.x} ${point.y}`
            : acc;
        }

        const nextIndex = index + 1;
        const nextPoint = this.points[nextIndex];
        const centerPoint = nextIndex !== lastIndex
          ? calculateCenterBetweenPoint(point, nextPoint)
          : nextPoint;

        return `${acc}Q ${point.x} ${point.y} ${centerPoint.x} ${centerPoint.y}`;
      }, '');
    },

    path() {
      return this.type === LINE_TYPES.sharp ? this.sharpPath : this.curvesPath;
    },
  },
  methods: {
    getCenterPoint() {
      const length = this.$refs.path.getTotalLength();

      return this.$refs.path.getPointAtLength(length * 0.5);
    },
  },
};
</script>
