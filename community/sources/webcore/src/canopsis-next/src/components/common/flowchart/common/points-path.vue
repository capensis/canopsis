<template lang="pug">
  path(
    :d="path",
    :fill="fill",
    pointer-events="all",
    v-on="$listeners"
  )
</template>

<script>
import { calculateCenterBetweenPoint, isCurvesControl } from '@/helpers/flowchart/points';

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
  },
  computed: {
    path() {
      const [firstPoint] = this.points;
      let path = `M ${firstPoint.x} ${firstPoint.y}`;

      for (let index = 1; index < this.points.length; index += 1) {
        const point = this.points[index];

        if (isCurvesControl(point.type)) {
          const nextPoint = this.points[index + 1];
          const centerPoint = isCurvesControl(nextPoint.type)
            ? calculateCenterBetweenPoint(point, nextPoint)
            : nextPoint;

          path += `Q ${point.x} ${point.y} ${centerPoint.x} ${centerPoint.y}`;
        } else {
          path += `L ${point.x} ${point.y}`;
        }
      }

      return path;
    },
  },
};
</script>
